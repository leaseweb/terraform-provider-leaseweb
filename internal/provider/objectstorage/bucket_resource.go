package objectstorage

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/objectstorage"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.ResourceWithConfigure   = &bucketResource{}
	_ resource.ResourceWithImportState = &bucketResource{}
)

type bucketResourceModel struct {
	CreationTime        types.String `tfsdk:"creation_time"`
	ID                  types.String `tfsdk:"id"`
	IsBeingDeleted      types.Bool   `tfsdk:"is_being_deleted"`
	IsVersioningEnabled types.Bool   `tfsdk:"is_versioning_enabled"`
	Name                types.String `tfsdk:"name"`
	ObjectCount         types.Int32  `tfsdk:"object_count"`
	ObjectStorageID     types.String `tfsdk:"object_storage_id"`
	Quota               types.Int32  `tfsdk:"quota"`
	Region              types.String `tfsdk:"region"`
	Used                types.Object `tfsdk:"used"`
}

type bucketUsedSpaceResourceModel struct {
	Unit  types.String  `tfsdk:"unit"`
	Value types.Float64 `tfsdk:"value"`
}

var bucketUsedSpaceAttributeTypes = map[string]attr.Type{
	"unit":  types.StringType,
	"value": types.Float64Type,
}

func adaptBucketToBucketResource(
	sdkBucket objectstorage.Bucket,
	objectStorageID string,
	diags *diag.Diagnostics,
	ctx context.Context,
) *bucketResourceModel {
	sdkUsed, _ := sdkBucket.GetUsedOk()
	used := utils.AdaptNullableSdkModelToResourceObject(
		sdkUsed,
		bucketUsedSpaceAttributeTypes,
		ctx,
		func(usedSpace objectstorage.BucketUsedSpace) bucketUsedSpaceResourceModel {
			return bucketUsedSpaceResourceModel{
				Unit:  basetypes.NewStringValue(usedSpace.GetUnit()),
				Value: basetypes.NewFloat64Value(float64(usedSpace.GetValue())),
			}
		},
		diags,
	)
	if diags.HasError() {
		return nil
	}

	// The API accepts a quota as an integer amount of GB but returns it as a
	// {value, unit} object. The spec pins unit to GB, so the value maps back
	// onto the same integer the practitioner configured.
	quota := basetypes.NewInt32Null()
	if sdkQuota, ok := sdkBucket.GetQuotaOk(); ok && sdkQuota != nil {
		quota = basetypes.NewInt32Value(int32(sdkQuota.GetValue()))
	}

	objectCount := basetypes.NewInt32Null()
	if sdkObjectCount, ok := sdkBucket.GetObjectCountOk(); ok && sdkObjectCount != nil {
		objectCount = basetypes.NewInt32Value(*sdkObjectCount)
	}

	creationTime, _ := sdkBucket.GetCreationTimeOk()

	return &bucketResourceModel{
		CreationTime:        adaptNullableTimeToStringValue(creationTime),
		ID:                  basetypes.NewStringValue(sdkBucket.GetName()),
		IsBeingDeleted:      basetypes.NewBoolValue(sdkBucket.GetIsBeingDeleted()),
		IsVersioningEnabled: basetypes.NewBoolValue(sdkBucket.GetIsVersioningEnabled()),
		Name:                basetypes.NewStringValue(sdkBucket.GetName()),
		ObjectCount:         objectCount,
		ObjectStorageID:     basetypes.NewStringValue(objectStorageID),
		Quota:               quota,
		Region:              basetypes.NewStringValue(sdkBucket.GetRegion()),
		Used:                used,
	}
}

type bucketResource struct {
	utils.ResourceAPI
}

// getBucket looks a bucket up by name. The API exposes no endpoint to fetch a
// single bucket, so the list is filtered client side. A nil return without
// diagnostics means the bucket no longer exists.
func (b bucketResource) getBucket(
	ctx context.Context,
	objectStorageID string,
	name string,
	diags *diag.Diagnostics,
) *objectstorage.Bucket {
	result, httpResponse, err := b.ObjectstorageAPI.GetObjectStorageBucketList(
		ctx,
		objectStorageID,
	).Execute()
	if err != nil {
		utils.SdkError(ctx, diags, err, httpResponse)
		return nil
	}

	for _, sdkBucket := range result.GetBuckets() {
		if sdkBucket.GetName() == name {
			return &sdkBucket
		}
	}

	return nil
}

func (b bucketResource) ImportState(
	ctx context.Context,
	request resource.ImportStateRequest,
	response *resource.ImportStateResponse,
) {
	idParts := strings.Split(request.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		utils.UnexpectedImportIdentifierError(
			&response.Diagnostics,
			"object_storage_id,name",
			request.ID,
		)
		return
	}

	response.Diagnostics.Append(response.State.SetAttribute(
		ctx,
		path.Root("object_storage_id"),
		idParts[0],
	)...)
	response.Diagnostics.Append(response.State.SetAttribute(
		ctx,
		path.Root("name"),
		idParts[1],
	)...)
	response.Diagnostics.Append(response.State.SetAttribute(
		ctx,
		path.Root("id"),
		idParts[1],
	)...)
}

func (b bucketResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	response *resource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		Description: "Manage an object storage bucket",
		Attributes: map[string]schema.Attribute{
			"object_storage_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the object storage the bucket belongs to",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Bucket name, used as the resource identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Bucket name. Must be 3-63 characters, alphanumeric and hyphens only",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 63),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9]$`),
						"must start and end with an alphanumeric character and may only contain alphanumeric characters and hyphens",
					),
				},
			},
			"is_versioning_enabled": schema.BoolAttribute{
				Required:    true,
				Description: "Whether versioning is enabled for the bucket",
			},
			"quota": schema.Int32Attribute{
				Optional:    true,
				Description: "Quota for the size of the bucket in GB. Omit for no quota",
				Validators: []validator.Int32{
					int32validator.Between(1, 1000000),
				},
			},
			"region": schema.StringAttribute{
				Computed:    true,
				Description: "Region the bucket resides in",
			},
			"object_count": schema.Int32Attribute{
				Computed:    true,
				Description: "Number of objects in the bucket",
			},
			"used": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Space used by the bucket",
				Attributes: map[string]schema.Attribute{
					"unit": schema.StringAttribute{
						Computed:    true,
						Description: "Unit of the used space, e.g. GB",
					},
					"value": schema.Float64Attribute{
						Computed:    true,
						Description: "Value of the used space",
					},
				},
			},
			"creation_time": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time the bucket was created",
			},
			"is_being_deleted": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the bucket is currently being deleted",
			},
		},
	}
}

func (b bucketResource) Create(
	ctx context.Context,
	request resource.CreateRequest,
	response *resource.CreateResponse,
) {
	var plan bucketResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	opts := objectstorage.NewCreateObjectStorageBucketOpts(
		plan.Name.ValueString(),
		plan.IsVersioningEnabled.ValueBool(),
	)
	if !plan.Quota.IsNull() && !plan.Quota.IsUnknown() {
		opts.SetQuota(plan.Quota.ValueInt32())
	}

	sdkBucket, httpResponse, err := b.ObjectstorageAPI.CreateObjectStorageBucket(
		ctx,
		plan.ObjectStorageID.ValueString(),
	).CreateObjectStorageBucketOpts(*opts).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	state := adaptBucketToBucketResource(
		*sdkBucket,
		plan.ObjectStorageID.ValueString(),
		&response.Diagnostics,
		ctx,
	)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (b bucketResource) Read(
	ctx context.Context,
	request resource.ReadRequest,
	response *resource.ReadResponse,
) {
	var originalState bucketResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &originalState)...)
	if response.Diagnostics.HasError() {
		return
	}

	objectStorageID := originalState.ObjectStorageID.ValueString()

	sdkBucket := b.getBucket(
		ctx,
		objectStorageID,
		originalState.Name.ValueString(),
		&response.Diagnostics,
	)
	if response.Diagnostics.HasError() {
		return
	}
	if sdkBucket == nil {
		response.State.RemoveResource(ctx)
		return
	}

	state := adaptBucketToBucketResource(
		*sdkBucket,
		objectStorageID,
		&response.Diagnostics,
		ctx,
	)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (b bucketResource) Update(
	ctx context.Context,
	request resource.UpdateRequest,
	response *resource.UpdateResponse,
) {
	var plan bucketResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	opts := objectstorage.NewUpdateObjectStorageBucketOpts(
		plan.IsVersioningEnabled.ValueBool(),
	)
	if plan.Quota.IsNull() {
		opts.SetQuotaNil()
	} else if !plan.Quota.IsUnknown() {
		opts.SetQuota(plan.Quota.ValueInt32())
	}

	sdkBucket, httpResponse, err := b.ObjectstorageAPI.UpdateObjectStorageBucket(
		ctx,
		plan.ObjectStorageID.ValueString(),
		plan.Name.ValueString(),
	).UpdateObjectStorageBucketOpts(*opts).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	state := adaptBucketToBucketResource(
		*sdkBucket,
		plan.ObjectStorageID.ValueString(),
		&response.Diagnostics,
		ctx,
	)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (b bucketResource) Delete(
	ctx context.Context,
	request resource.DeleteRequest,
	response *resource.DeleteResponse,
) {
	var state bucketResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	objectStorageID := state.ObjectStorageID.ValueString()
	name := state.Name.ValueString()

	httpResponse, err := b.ObjectstorageAPI.DeleteObjectStorageBucket(
		ctx,
		objectStorageID,
		name,
	).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	// Deletion is asynchronous, so wait until the bucket has actually gone
	// before reporting success. Otherwise a recreate of a bucket with the same
	// name races the deletion.
	err = b.waitUntilBucketDeleted(
		ctx,
		objectStorageID,
		name,
		&response.Diagnostics,
	)
	if err != nil && !response.Diagnostics.HasError() {
		utils.GeneralError(&response.Diagnostics, ctx, err)
	}
}

// waitUntilBucketDeleted handles polling with retry and timeout.
func (b bucketResource) waitUntilBucketDeleted(
	ctx context.Context,
	objectStorageID string,
	name string,
	diags *diag.Diagnostics,
) error {
	// Create a constant backoff with a 10-second retry interval
	bo := backoff.NewConstantBackOff(10 * time.Second)

	// Set the retry limit to 30 retries (5 minutes total)
	retryCount := 0
	maxRetries := 30

	// Start polling and retrying
	for {
		if retryCount >= maxRetries {
			return errors.New("timed out waiting for bucket to be deleted after 5 minutes")
		}

		sdkBucket := b.getBucket(ctx, objectStorageID, name, diags)
		if diags.HasError() {
			return fmt.Errorf(
				"could not determine deletion status of bucket %q",
				name,
			)
		}

		// The bucket is gone once it drops out of the bucket list. Deletion is
		// equally done once it stops reporting itself as being deleted, which
		// is the flag the API sets while the removal is still in flight.
		if sdkBucket == nil || !sdkBucket.GetIsBeingDeleted() {
			return nil
		}

		// Sleep for the backoff interval before retrying, while staying
		// responsive to a cancelled apply
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(bo.NextBackOff()):
		}
		retryCount++
	}
}

func NewBucketResource() resource.Resource {
	return &bucketResource{
		ResourceAPI: utils.ResourceAPI{
			Name: "object_storage_bucket",
		},
	}
}
