package objectstorage

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	_ resource.ResourceWithConfigure   = &accessKeyResource{}
	_ resource.ResourceWithImportState = &accessKeyResource{}
)

type accessKeyResourceModel struct {
	AccessKey       types.String `tfsdk:"access_key"`
	AccountID       types.String `tfsdk:"account_id"`
	DisplayName     types.String `tfsdk:"display_name"`
	ExpiresAt       types.String `tfsdk:"expires_at"`
	ID              types.String `tfsdk:"id"`
	ObjectStorageID types.String `tfsdk:"object_storage_id"`
	SecretAccessKey types.String `tfsdk:"secret_access_key"`
	UserID          types.String `tfsdk:"user_id"`
}

type accessKeyResource struct {
	utils.ResourceAPI
}

// getAccessKey looks an access key up by id. The API exposes no endpoint to
// fetch a single access key, so the list is filtered client side. A nil return
// without diagnostics means the access key no longer exists.
func (a accessKeyResource) getAccessKey(
	ctx context.Context,
	objectStorageID string,
	userID string,
	id string,
	diags *diag.Diagnostics,
) *objectstorage.AccessKey {
	result, httpResponse, err := a.ObjectstorageAPI.GetObjectStorageAccessKeyList(
		ctx,
		objectStorageID,
		userID,
	).Execute()
	if err != nil {
		utils.SdkError(ctx, diags, err, httpResponse)
		return nil
	}

	for _, sdkAccessKey := range result.GetKeys() {
		if sdkAccessKey.GetId() == id {
			return &sdkAccessKey
		}
	}

	return nil
}

func (a accessKeyResource) ImportState(
	ctx context.Context,
	request resource.ImportStateRequest,
	response *resource.ImportStateResponse,
) {
	idParts := strings.Split(request.ID, ",")

	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		utils.UnexpectedImportIdentifierError(
			&response.Diagnostics,
			"object_storage_id,user_id,id",
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
		path.Root("user_id"),
		idParts[1],
	)...)
	response.Diagnostics.Append(response.State.SetAttribute(
		ctx,
		path.Root("id"),
		idParts[2],
	)...)
}

func (a accessKeyResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	response *resource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		MarkdownDescription: "Manage an object storage access key.\n\n**Note:** `access_key` & `secret_access_key` are only returned when the access key is created. They cannot be retrieved afterwards, so an imported access key has both values set to null.",
		Attributes: map[string]schema.Attribute{
			"object_storage_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the object storage the access key belongs to",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the user the access key belongs to",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Access key ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"expires_at": schema.StringAttribute{
				Optional:    true,
				Description: "The date and time the access key expires, in RFC3339 format (`yyyy-mm-ddThh:mm:ssZ`). Omit for a key that does not expire",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the access key",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the account the access key belongs to",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"access_key": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The access key. Only returned when the access key is created",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"secret_access_key": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The access key secret. Only returned when the access key is created",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (a accessKeyResource) Create(
	ctx context.Context,
	request resource.CreateRequest,
	response *resource.CreateResponse,
) {
	var plan accessKeyResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	opts := objectstorage.NewCreateObjectStorageAccessKeyOpts()
	if !plan.ExpiresAt.IsNull() && !plan.ExpiresAt.IsUnknown() {
		expiresAt, err := time.Parse(time.RFC3339, plan.ExpiresAt.ValueString())
		if err != nil {
			response.Diagnostics.AddAttributeError(
				path.Root("expires_at"),
				"Invalid date and time",
				"Attribute expires_at must be specified using the RFC3339 format (`yyyy-mm-ddThh:mm:ssZ`)",
			)
			return
		}
		opts.SetExpiresAt(expiresAt)
	}

	sdkAccessKey, httpResponse, err := a.ObjectstorageAPI.CreateObjectStorageAccessKey(
		ctx,
		plan.ObjectStorageID.ValueString(),
		plan.UserID.ValueString(),
	).CreateObjectStorageAccessKeyOpts(*opts).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	expiresAt, _ := sdkAccessKey.GetExpiresAtOk()

	state := accessKeyResourceModel{
		AccessKey:       basetypes.NewStringValue(sdkAccessKey.GetAccessKey()),
		AccountID:       basetypes.NewStringValue(sdkAccessKey.GetAccountId()),
		DisplayName:     basetypes.NewStringValue(sdkAccessKey.GetDisplayName()),
		ExpiresAt:       adaptNullableTimeToStringValue(expiresAt),
		ID:              basetypes.NewStringValue(sdkAccessKey.GetId()),
		ObjectStorageID: plan.ObjectStorageID,
		SecretAccessKey: basetypes.NewStringValue(sdkAccessKey.GetSecretAccessKey()),
		UserID:          plan.UserID,
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (a accessKeyResource) Read(
	ctx context.Context,
	request resource.ReadRequest,
	response *resource.ReadResponse,
) {
	var originalState accessKeyResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &originalState)...)
	if response.Diagnostics.HasError() {
		return
	}

	sdkAccessKey := a.getAccessKey(
		ctx,
		originalState.ObjectStorageID.ValueString(),
		originalState.UserID.ValueString(),
		originalState.ID.ValueString(),
		&response.Diagnostics,
	)
	if response.Diagnostics.HasError() {
		return
	}
	if sdkAccessKey == nil {
		response.State.RemoveResource(ctx)
		return
	}

	expiresAt, _ := sdkAccessKey.GetExpiresAtOk()

	// access_key & secret_access_key are only returned on creation, so the
	// values already in state are carried over untouched.
	state := accessKeyResourceModel{
		AccessKey:       originalState.AccessKey,
		AccountID:       basetypes.NewStringValue(sdkAccessKey.GetAccountId()),
		DisplayName:     basetypes.NewStringValue(sdkAccessKey.GetDisplayName()),
		ExpiresAt:       adaptNullableTimeToStringValue(expiresAt),
		ID:              basetypes.NewStringValue(sdkAccessKey.GetId()),
		ObjectStorageID: originalState.ObjectStorageID,
		SecretAccessKey: originalState.SecretAccessKey,
		UserID:          originalState.UserID,
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (a accessKeyResource) Update(
	_ context.Context,
	_ resource.UpdateRequest,
	_ *resource.UpdateResponse,
) {
	// Every attribute requires replacement, so Update is never called.
}

func (a accessKeyResource) Delete(
	ctx context.Context,
	request resource.DeleteRequest,
	response *resource.DeleteResponse,
) {
	var state accessKeyResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	httpResponse, err := a.ObjectstorageAPI.DeleteObjectStorageAccessKey(
		ctx,
		state.ObjectStorageID.ValueString(),
		state.UserID.ValueString(),
		state.ID.ValueString(),
	).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
	}
}

func NewAccessKeyResource() resource.Resource {
	return &accessKeyResource{
		ResourceAPI: utils.ResourceAPI{
			Name: "object_storage_access_key",
		},
	}
}
