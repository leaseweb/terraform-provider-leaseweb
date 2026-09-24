package objectstorage

import (
	"context"
	"strings"

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
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

type userResourceModel struct {
	FullName        types.String `tfsdk:"full_name"`
	Groups          types.Set    `tfsdk:"groups"`
	ID              types.String `tfsdk:"id"`
	ObjectStorageID types.String `tfsdk:"object_storage_id"`
	UniqueName      types.String `tfsdk:"unique_name"`
}

func adaptUserToUserResource(
	sdkUser objectstorage.User,
	objectStorageID string,
	diags *diag.Diagnostics,
	ctx context.Context,
) *userResourceModel {
	groups, groupDiags := types.SetValueFrom(
		ctx,
		types.StringType,
		sdkUser.GetGroups(),
	)
	if groupDiags.HasError() {
		diags.Append(groupDiags...)
		return nil
	}

	return &userResourceModel{
		FullName:        basetypes.NewStringValue(sdkUser.GetFullName()),
		Groups:          groups,
		ID:              basetypes.NewStringValue(sdkUser.GetId()),
		ObjectStorageID: basetypes.NewStringValue(objectStorageID),
		UniqueName:      basetypes.NewStringValue(sdkUser.GetUniqueName()),
	}
}

type userResource struct {
	utils.ResourceAPI
}

func (u userResource) ImportState(
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

func (u userResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	response *resource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		Description: "Manage an object storage user",
		Attributes: map[string]schema.Attribute{
			"object_storage_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the object storage the user belongs to",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "User ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"full_name": schema.StringAttribute{
				Required:    true,
				Description: "Full name of the user",
			},
			"unique_name": schema.StringAttribute{
				Required:    true,
				Description: "Unique name identifier for the user",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"groups": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "IDs of the groups the user belongs to",
			},
		},
	}
}

func (u userResource) Create(
	ctx context.Context,
	request resource.CreateRequest,
	response *resource.CreateResponse,
) {
	var plan userResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	var groups []string
	response.Diagnostics.Append(plan.Groups.ElementsAs(ctx, &groups, false)...)
	if response.Diagnostics.HasError() {
		return
	}

	opts := objectstorage.NewCreateObjectStorageUserOpts(
		plan.FullName.ValueString(),
		plan.UniqueName.ValueString(),
		groups,
	)

	sdkUser, httpResponse, err := u.ObjectstorageAPI.CreateObjectStorageUser(
		ctx,
		plan.ObjectStorageID.ValueString(),
	).CreateObjectStorageUserOpts(*opts).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	state := adaptUserToUserResource(
		*sdkUser,
		plan.ObjectStorageID.ValueString(),
		&response.Diagnostics,
		ctx,
	)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (u userResource) Read(
	ctx context.Context,
	request resource.ReadRequest,
	response *resource.ReadResponse,
) {
	var originalState userResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &originalState)...)
	if response.Diagnostics.HasError() {
		return
	}

	objectStorageID := originalState.ObjectStorageID.ValueString()

	sdkUser, httpResponse, err := u.ObjectstorageAPI.GetObjectStorageUser(
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

	state := adaptUserToUserResource(
		*sdkUser,
		objectStorageID,
		&response.Diagnostics,
		ctx,
	)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (u userResource) Update(
	ctx context.Context,
	request resource.UpdateRequest,
	response *resource.UpdateResponse,
) {
	var plan userResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	var groups []string
	response.Diagnostics.Append(plan.Groups.ElementsAs(ctx, &groups, false)...)
	if response.Diagnostics.HasError() {
		return
	}

	opts := objectstorage.NewUpdateObjectStorageUserOpts(
		plan.FullName.ValueString(),
		groups,
	)

	sdkUser, httpResponse, err := u.ObjectstorageAPI.UpdateObjectStorageUser(
		ctx,
		plan.ObjectStorageID.ValueString(),
		plan.ID.ValueString(),
	).UpdateObjectStorageUserOpts(*opts).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	state := adaptUserToUserResource(
		*sdkUser,
		plan.ObjectStorageID.ValueString(),
		&response.Diagnostics,
		ctx,
	)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (u userResource) Delete(
	ctx context.Context,
	request resource.DeleteRequest,
	response *resource.DeleteResponse,
) {
	var state userResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	httpResponse, err := u.ObjectstorageAPI.DeleteObjectStorageUser(
		ctx,
		state.ObjectStorageID.ValueString(),
		state.ID.ValueString(),
	).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
	}
}

func NewUserResource() resource.Resource {
	return &userResource{
		ResourceAPI: utils.ResourceAPI{
			Name: "object_storage_user",
		},
	}
}
