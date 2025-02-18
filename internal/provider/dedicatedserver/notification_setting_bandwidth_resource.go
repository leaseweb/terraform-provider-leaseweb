package dedicatedserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedserver/v2"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.Resource              = &notificationSettingBandwidthResource{}
	_ resource.ResourceWithConfigure = &notificationSettingBandwidthResource{}
)

type notificationSettingBandwidthResource struct {
	utils.ResourceAPI
}

type notificationSettingBandwidthResourceModel struct {
	ID                types.String `tfsdk:"id"`
	DedicatedServerID types.String `tfsdk:"dedicated_server_id"`
	Frequency         types.String `tfsdk:"frequency"`
	Threshold         types.String `tfsdk:"threshold"`
	Unit              types.String `tfsdk:"unit"`
}

func NewNotificationSettingBandwidthResource() resource.Resource {
	return &notificationSettingBandwidthResource{
		ResourceAPI: utils.ResourceAPI{
			Name: "dedicated_server_notification_setting_bandwidth",
		},
	}
}

func (n *notificationSettingBandwidthResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The notification setting bandwidth unique identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dedicated_server_id": schema.StringAttribute{
				Required:    true,
				Description: "The server unique identifier",
			},
			"frequency": schema.StringAttribute{
				Required:    true,
				Description: "The notification frequency. Valid options can be *DAILY* or *WEEKLY* or *MONTHLY*.",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"DAILY", "WEEKLY", "MONTHLY"}...),
				},
			},
			"threshold": schema.StringAttribute{
				Required:    true,
				Description: "Threshold Value. Value can be a number greater than 0.",
				Validators: []validator.String{
					greaterThanZero(),
				},
			},
			"unit": schema.StringAttribute{
				Required:    true,
				Description: "The notification unit. Valid options can be *Mbps* or *Gbps*.",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"Mbps", "Gbps"}...),
				},
			},
		},
	}
}

func (n *notificationSettingBandwidthResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan notificationSettingBandwidthResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := dedicatedserver.NewBandwidthNotificationSettingOpts(
		plan.Frequency.ValueString(),
		plan.Threshold.ValueString(),
		plan.Unit.ValueString(),
	)
	request := n.DedicatedserverAPI.CreateBandwidthNotificationSetting(
		ctx,
		plan.DedicatedServerID.ValueString(),
	).BandwidthNotificationSettingOpts(*opts)
	result, response, err := request.Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, response)
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			notificationSettingBandwidthResourceModel{
				ID:                types.StringValue(result.GetId()),
				Frequency:         types.StringValue(result.GetFrequency()),
				Threshold:         types.StringValue(result.GetThreshold()),
				Unit:              types.StringValue(result.GetUnit()),
				DedicatedServerID: plan.DedicatedServerID,
			},
		)...,
	)
}

func (n *notificationSettingBandwidthResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state notificationSettingBandwidthResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := n.DedicatedserverAPI.GetBandwidthNotificationSetting(
		ctx,
		state.DedicatedServerID.ValueString(),
		state.ID.ValueString(),
	)
	result, response, err := request.Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, response)
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			notificationSettingBandwidthResourceModel{
				ID:                types.StringValue(result.GetId()),
				Frequency:         types.StringValue(result.GetFrequency()),
				Threshold:         types.StringValue(result.GetThreshold()),
				Unit:              types.StringValue(result.GetUnit()),
				DedicatedServerID: state.DedicatedServerID,
			},
		)...,
	)
}

func (n *notificationSettingBandwidthResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan notificationSettingBandwidthResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := dedicatedserver.NewBandwidthNotificationSettingOpts(
		plan.Frequency.ValueString(),
		plan.Threshold.ValueString(),
		plan.Unit.ValueString(),
	)
	request := n.DedicatedserverAPI.UpdateBandwidthNotificationSetting(
		ctx,
		plan.DedicatedServerID.ValueString(),
		plan.ID.ValueString(),
	).BandwidthNotificationSettingOpts(*opts)
	result, response, err := request.Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, response)
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			notificationSettingBandwidthResourceModel{
				ID:                plan.ID,
				DedicatedServerID: plan.DedicatedServerID,
				Frequency:         types.StringValue(result.GetFrequency()),
				Threshold:         types.StringValue(result.GetThreshold()),
				Unit:              types.StringValue(result.GetUnit()),
			},
		)...,
	)
}

func (n *notificationSettingBandwidthResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state notificationSettingBandwidthResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := n.DedicatedserverAPI.DeleteBandwidthNotificationSetting(
		ctx,
		state.DedicatedServerID.ValueString(),
		state.ID.ValueString(),
	)
	response, err := request.Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, response)
	}
}
