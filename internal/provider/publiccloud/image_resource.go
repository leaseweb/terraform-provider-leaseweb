package publiccloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publiccloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.ResourceWithConfigure = &imageResource{}
)

type imageResourceModel struct {
	ID           types.String `tfsdk:"id"`
	InstanceID   types.String `tfsdk:"instance_id"`
	Name         types.String `tfsdk:"name"`
	Custom       types.Bool   `tfsdk:"custom"`
	State        types.String `tfsdk:"state"`
	MarketApps   types.List   `tfsdk:"market_apps"`
	StorageTypes types.List   `tfsdk:"storage_types"`
	Flavour      types.String `tfsdk:"flavour"`
	Region       types.String `tfsdk:"region"`
}

func adaptImageDetailsToImageResource(
	ctx context.Context,
	imageDetails publiccloud.ImageDetails,
	diags *diag.Diagnostics,
) *imageResourceModel {
	marketApps, marketAppsDiags := basetypes.NewListValueFrom(
		ctx,
		basetypes.StringType{},
		imageDetails.MarketApps,
	)
	if marketAppsDiags.HasError() {
		diags.Append(marketAppsDiags...)
		return nil
	}

	storageTypes, storageTypesDiags := basetypes.NewListValueFrom(
		ctx,
		basetypes.StringType{},
		imageDetails.StorageTypes,
	)
	if storageTypesDiags.HasError() {
		diags.Append(storageTypesDiags...)
		return nil
	}

	image := imageResourceModel{
		ID:           basetypes.NewStringValue(imageDetails.GetId()),
		Name:         basetypes.NewStringValue(imageDetails.GetName()),
		Custom:       basetypes.NewBoolValue(imageDetails.GetCustom()),
		State:        basetypes.NewStringValue(string(imageDetails.GetState())),
		MarketApps:   marketApps,
		StorageTypes: storageTypes,
		Flavour:      basetypes.NewStringValue(string(imageDetails.GetFlavour())),
		Region:       basetypes.NewStringValue(string(imageDetails.GetRegion())),
	}

	return &image
}

type imageResource struct {
	utils.ResourceAPI
}

func (i *imageResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	response *resource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Can be either an Operating System or a UUID in case of a Custom Image",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required: true,
				Description: `
The id of the instance which the custom image is based on. The following rules apply:
  - instance exists for instanceId
  - instance has state *STOPPED*
  - instance has a maximum rootDiskSize of 100 GB
  - instance OS must not be *windows*`,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Custom image name",
			},
			"custom": schema.BoolAttribute{
				Computed:    true,
				Description: "Standard or Custom image",
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"market_apps": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"storage_types": schema.ListAttribute{
				Computed:    true,
				Description: "The supported storage types for the instance type",
				ElementType: types.StringType,
			},
			"flavour": schema.StringAttribute{
				Computed: true,
			},
			"region": schema.StringAttribute{
				Computed: true,
			},
		},
	}

	utils.AddUnsupportedActionsNotation(
		response,
		[]utils.Action{utils.DeleteAction},
	)
}

func (i *imageResource) Create(
	ctx context.Context,
	request resource.CreateRequest,
	response *resource.CreateResponse,
) {
	var plan imageResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	imageDetails, httpResponse, err := i.PubliccloudAPI.CreateImage(ctx).
		CreateImageOpts(
			*publiccloud.NewCreateImageOpts(
				plan.Name.ValueString(),
				plan.InstanceID.ValueString(),
			),
		).
		Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	state := adaptImageDetailsToImageResource(
		ctx,
		*imageDetails,
		&response.Diagnostics,
	)
	if response.Diagnostics.HasError() {
		return
	}
	// instanceId has to be set manually as it isn't returned from the API
	state.InstanceID = plan.InstanceID

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (i *imageResource) Read(
	ctx context.Context,
	request resource.ReadRequest,
	response *resource.ReadResponse,
) {
	var currentState imageResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &currentState)...)
	if response.Diagnostics.HasError() {
		return
	}

	images := getAllImages(ctx, i.PubliccloudAPI, &response.Diagnostics)
	if response.Diagnostics.HasError() {
		return
	}
	imageDetails := images.findById(currentState.ID.ValueString())
	if imageDetails == nil {
		utils.GeneralError(
			&response.Diagnostics,
			ctx,
			fmt.Errorf(
				"imageDetails  %s not found",
				currentState.ID.ValueString(),
			),
		)
		return
	}

	state := adaptImageDetailsToImageResource(
		ctx,
		*imageDetails,
		&response.Diagnostics,
	)
	if response.Diagnostics.HasError() {
		return
	}
	// instanceId has to be set manually as it isn't returned from the API
	state.InstanceID = currentState.InstanceID

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (i *imageResource) Update(
	ctx context.Context,
	request resource.UpdateRequest,
	response *resource.UpdateResponse,
) {
	var plan imageResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	imageDetails, httpResponse, err := i.PubliccloudAPI.UpdateImage(
		ctx,
		plan.ID.ValueString(),
	).UpdateImageOpts(*publiccloud.NewUpdateImageOpts(plan.Name.ValueString())).
		Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	state := adaptImageDetailsToImageResource(
		ctx,
		*imageDetails,
		&response.Diagnostics,
	)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

// Delete does nothing as there is no endpoint to delete an Image.
func (i *imageResource) Delete(
	_ context.Context,
	_ resource.DeleteRequest,
	_ *resource.DeleteResponse,
) {
}

func NewImageResource() resource.Resource {
	return &imageResource{
		ResourceAPI: utils.ResourceAPI{
			Name: "public_cloud_image",
		},
	}
}
