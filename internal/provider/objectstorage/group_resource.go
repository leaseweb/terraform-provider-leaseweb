package objectstorage

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/objectstorage"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.ResourceWithConfigure   = &groupResource{}
	_ resource.ResourceWithImportState = &groupResource{}
)

type groupResourceModel struct {
	DisplayName     types.String         `tfsdk:"display_name"`
	ID              types.String         `tfsdk:"id"`
	ObjectStorageID types.String         `tfsdk:"object_storage_id"`
	S3Policies      jsontypes.Normalized `tfsdk:"s3_policies"`
	UniqueName      types.String         `tfsdk:"unique_name"`
}

func adaptGroupToGroupResource(
	sdkGroup objectstorage.Group,
	objectStorageID string,
) *groupResourceModel {
	s3Policies, _ := sdkGroup.GetS3PoliciesOk()

	return &groupResourceModel{
		DisplayName:     basetypes.NewStringValue(sdkGroup.GetDisplayName()),
		ID:              basetypes.NewStringValue(sdkGroup.GetId()),
		ObjectStorageID: basetypes.NewStringValue(objectStorageID),
		S3Policies:      jsontypes.NewNormalizedPointerValue(s3Policies),
		UniqueName:      basetypes.NewStringValue(sdkGroup.GetUniqueName()),
	}
}

type groupResource struct {
	utils.ResourceAPI
}

func (g groupResource) ImportState(
	ctx context.Context,
	request resource.ImportStateRequest,
	response *resource.ImportStateResponse,
) {
	idParts := strings.Split(request.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		utils.UnexpectedImportIdentifierError(
			&response.Diagnostics,
			"object_storage_id,id",
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
		path.Root("id"),
		idParts[1],
	)...)
}

func (g groupResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	response *resource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		Description: "Manage an object storage group",
		Attributes: map[string]schema.Attribute{
			"object_storage_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the object storage the group belongs to",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Group ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the group",
			},
			"unique_name": schema.StringAttribute{
				Required:    true,
				Description: "Unique name identifier for the group",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"s3_policies": schema.StringAttribute{
				Optional:    true,
				CustomType:  jsontypes.NormalizedType{},
				Description: "A JSON encoded IAM style policy document defining the S3 permissions for users in this group. Semantically equal JSON is treated as unchanged, so formatting & property order returned by the API do not cause a diff",
			},
		},
	}
}

func (g groupResource) Create(
	ctx context.Context,
	request resource.CreateRequest,
	response *resource.CreateResponse,
) {
	var plan groupResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	opts := objectstorage.NewCreateObjectStorageGroupOpts(
		plan.DisplayName.ValueString(),
		plan.UniqueName.ValueString(),
	)
	if !plan.S3Policies.IsNull() && !plan.S3Policies.IsUnknown() {
		opts.SetS3Policies(plan.S3Policies.ValueString())
	}

	sdkGroup, httpResponse, err := g.ObjectstorageAPI.CreateObjectStorageGroup(
		ctx,
		plan.ObjectStorageID.ValueString(),
	).CreateObjectStorageGroupOpts(*opts).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	state := adaptGroupToGroupResource(
		*sdkGroup,
		plan.ObjectStorageID.ValueString(),
	)

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (g groupResource) Read(
	ctx context.Context,
	request resource.ReadRequest,
	response *resource.ReadResponse,
) {
	var originalState groupResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &originalState)...)
	if response.Diagnostics.HasError() {
		return
	}

	objectStorageID := originalState.ObjectStorageID.ValueString()

	sdkGroup, httpResponse, err := g.ObjectstorageAPI.GetObjectStorageGroup(
		ctx,
		objectStorageID,
		originalState.ID.ValueString(),
	).Execute()
	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == 404 {
			response.State.RemoveResource(ctx)
			return
		}

		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	state := adaptGroupToGroupResource(*sdkGroup, objectStorageID)

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (g groupResource) Update(
	ctx context.Context,
	request resource.UpdateRequest,
	response *resource.UpdateResponse,
) {
	var plan groupResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	opts := objectstorage.NewUpdateObjectStorageGroupOpts(
		plan.DisplayName.ValueString(),
	)
	if plan.S3Policies.IsNull() {
		opts.SetS3PoliciesNil()
	} else if !plan.S3Policies.IsUnknown() {
		opts.SetS3Policies(plan.S3Policies.ValueString())
	}

	sdkGroup, httpResponse, err := g.ObjectstorageAPI.UpdateObjectStorageGroup(
		ctx,
		plan.ObjectStorageID.ValueString(),
		plan.ID.ValueString(),
	).UpdateObjectStorageGroupOpts(*opts).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	state := adaptGroupToGroupResource(
		*sdkGroup,
		plan.ObjectStorageID.ValueString(),
	)

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (g groupResource) Delete(
	ctx context.Context,
	request resource.DeleteRequest,
	response *resource.DeleteResponse,
) {
	var state groupResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	httpResponse, err := g.ObjectstorageAPI.DeleteObjectStorageGroup(
		ctx,
		state.ObjectStorageID.ValueString(),
		state.ID.ValueString(),
	).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
	}
}

func NewGroupResource() resource.Resource {
	return &groupResource{
		ResourceAPI: utils.ResourceAPI{
			Name: "object_storage_group",
		},
	}
}
