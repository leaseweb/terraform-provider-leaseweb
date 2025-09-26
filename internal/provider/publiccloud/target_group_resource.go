package publiccloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publiccloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.ResourceWithConfigure   = &targetGroupResource{}
	_ resource.ResourceWithImportState = &targetGroupResource{}
)

type targetGroupResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Protocol    types.String `tfsdk:"protocol"`
	Port        types.Int32  `tfsdk:"port"`
	Region      types.String `tfsdk:"region"`
	HealthCheck types.Object `tfsdk:"health_check"`
}

func adaptTargetGroupToTargetGroupResource(
	sdkTargetGroup publiccloud.TargetGroup,
	ctx context.Context,
	diags *diag.Diagnostics,
) *targetGroupResourceModel {
	targetGroup := targetGroupResourceModel{
		ID:       basetypes.NewStringValue(sdkTargetGroup.GetId()),
		Name:     basetypes.NewStringValue(sdkTargetGroup.GetName()),
		Protocol: basetypes.NewStringValue(string(sdkTargetGroup.GetProtocol())),
		Port:     basetypes.NewInt32Value(sdkTargetGroup.GetPort()),
		Region:   basetypes.NewStringValue(string(sdkTargetGroup.GetRegion())),
	}

	sdkHealthCheck, _ := sdkTargetGroup.GetHealthCheckOk()

	healthCheck := utils.AdaptNullableSdkModelToResourceObject(
		sdkHealthCheck,
		map[string]attr.Type{
			"protocol": types.StringType,
			"method":   types.StringType,
			"uri":      types.StringType,
			"host":     types.StringType,
			"port":     types.Int32Type,
		},
		ctx,
		func(sdkHealthCheck publiccloud.HealthCheck) healthCheckResourceModel {
			healthCheck := healthCheckResourceModel{
				Protocol: basetypes.NewStringValue(string(sdkHealthCheck.GetProtocol())),
				URI:      basetypes.NewStringValue(sdkHealthCheck.GetUri()),
				Port:     basetypes.NewInt32Value(sdkHealthCheck.GetPort()),
			}

			method, _ := sdkHealthCheck.GetMethodOk()
			healthCheck.Method = basetypes.NewStringPointerValue((*string)(method))

			host, _ := sdkHealthCheck.GetHostOk()
			healthCheck.Host = basetypes.NewStringPointerValue(host)

			return healthCheck
		},
		diags,
	)
	if diags.HasError() {
		return nil
	}
	targetGroup.HealthCheck = healthCheck

	return &targetGroup
}

type healthCheckResourceModel struct {
	Protocol types.String `tfsdk:"protocol"`
	Method   types.String `tfsdk:"method"`
	URI      types.String `tfsdk:"uri"`
	Host     types.String `tfsdk:"host"`
	Port     types.Int32  `tfsdk:"port"`
}

func (h healthCheckResourceModel) generateOpts() publiccloud.HealthCheckOpts {
	opts := publiccloud.NewHealthCheckOpts(
		publiccloud.Protocol(h.Protocol.ValueString()),
		h.URI.ValueString(),
		h.Port.ValueInt32(),
	)

	if !h.Method.IsNull() {
		opts.SetMethod(publiccloud.HttpMethodOpt(h.Method.ValueString()))
	}
	opts.Host = utils.AdaptStringPointerValueToNullableString(h.Host)

	return *opts
}

type targetGroupResource struct {
	utils.ResourceAPI
}

func (t *targetGroupResource) ImportState(
	ctx context.Context,
	request resource.ImportStateRequest,
	response *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(
		ctx,
		path.Root("id"),
		request,
		response,
	)
}

func (t *targetGroupResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	response *resource.SchemaResponse,
) {
	warningError := "**WARNING!** Changing this value once running will cause this target group to be destroyed and a new one to be created."

	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The Name of the target group",
			},
			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "Valid options are " + utils.StringTypeArrayToMarkdown(publiccloud.AllowedProtocolEnumValues) + "\n" + warningError,
				Validators: []validator.String{
					stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publiccloud.AllowedProtocolEnumValues)...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"port": schema.Int32Attribute{
				Required:    true,
				Description: "The port of the target group",
				Validators: []validator.Int32{
					int32validator.Between(1, 65535),
				},
			},
			"region": schema.StringAttribute{
				Required:    true,
				Description: "Valid options are " + utils.StringTypeArrayToMarkdown(publiccloud.AllowedRegionNameEnumValues) + "\n" + warningError,
				Validators: []validator.String{
					stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publiccloud.AllowedRegionNameEnumValues)...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"health_check": schema.SingleNestedAttribute{
				Description: "**WARNING!** Removing health_check once running will cause this target group to be destroyed and a new one to be created.",
				Optional:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplaceIf(
						func(
							ctx context.Context,
							request planmodifier.ObjectRequest,
							response *objectplanmodifier.RequiresReplaceIfFuncResponse,
						) {
							if request.PlanValue.IsNull() {
								response.RequiresReplace = true
							}
						},
						"",
						"",
					),
				},
				Attributes: map[string]schema.Attribute{
					"protocol": schema.StringAttribute{
						Required:    true,
						Description: "Valid options are " + utils.StringTypeArrayToMarkdown(publiccloud.AllowedProtocolEnumValues),
						Validators: []validator.String{
							stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publiccloud.AllowedProtocolEnumValues)...),
						},
					},
					"method": schema.StringAttribute{
						Description: "Required if `protocol` is `HTTP` or `HTTPS`. Valid options are " + utils.StringTypeArrayToMarkdown(publiccloud.AllowedHttpMethodEnumValues),
						Validators: []validator.String{
							stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publiccloud.AllowedHttpMethodEnumValues)...),
						},
						Optional: true,
					},
					"uri": schema.StringAttribute{
						Required:    true,
						Description: "URI to check in the target instances",
					},
					"host": schema.StringAttribute{
						Description: "Host for the health check if any",
						Optional:    true,
					},
					"port": schema.Int32Attribute{
						Required:    true,
						Description: "Port number",
						Validators: []validator.Int32{
							int32validator.Between(1, 65535),
						},
					},
				},
			},
		},
	}
}

func (t *targetGroupResource) Create(
	ctx context.Context,
	request resource.CreateRequest,
	response *resource.CreateResponse,
) {
	var plan targetGroupResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	opts := publiccloud.NewCreateTargetGroupOpts(
		plan.Name.ValueString(),
		publiccloud.Protocol(plan.Protocol.ValueString()),
		plan.Port.ValueInt32(),
		publiccloud.RegionName(plan.Region.ValueString()),
	)
	if !plan.HealthCheck.IsNull() {
		healthCheck := healthCheckResourceModel{}
		healthCheckDiags := plan.HealthCheck.As(
			ctx,
			&healthCheck,
			basetypes.ObjectAsOptions{},
		)
		if healthCheckDiags != nil {
			response.Diagnostics.Append(healthCheckDiags...)
			return
		}

		opts.SetHealthCheck(healthCheck.generateOpts())
	}

	sdkTargetGroup, httpResponse, err := t.PubliccloudAPI.CreateTargetGroup(ctx).
		CreateTargetGroupOpts(*opts).
		Execute()

	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	targetGroup := adaptTargetGroupToTargetGroupResource(
		*sdkTargetGroup,
		ctx,
		&response.Diagnostics,
	)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, targetGroup)...)
}

func (t *targetGroupResource) Read(
	ctx context.Context,
	request resource.ReadRequest,
	response *resource.ReadResponse,
) {
	var state targetGroupResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	sdkTargetGroup, httpResponse, err := t.PubliccloudAPI.
		GetTargetGroup(ctx, state.ID.ValueString()).
		Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	targetGroup := adaptTargetGroupToTargetGroupResource(
		*sdkTargetGroup,
		ctx,
		&response.Diagnostics,
	)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, targetGroup)...)
}

func (t *targetGroupResource) Update(
	ctx context.Context,
	request resource.UpdateRequest,
	response *resource.UpdateResponse,
) {
	var plan targetGroupResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	opts := publiccloud.NewUpdateTargetGroupOpts()
	opts.SetName(plan.Name.ValueString())
	opts.SetPort(plan.Port.ValueInt32())
	if !plan.HealthCheck.IsNull() {
		healthCheck := healthCheckResourceModel{}
		healthCheckDiags := plan.HealthCheck.As(
			ctx,
			&healthCheck,
			basetypes.ObjectAsOptions{},
		)
		if healthCheckDiags != nil {
			response.Diagnostics.Append(healthCheckDiags...)
			return
		}

		opts.SetHealthCheck(healthCheck.generateOpts())
	}

	sdkTargetGroup, httpResponse, err := t.PubliccloudAPI.
		UpdateTargetGroup(ctx, plan.ID.ValueString()).
		UpdateTargetGroupOpts(*opts).
		Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	targetGroup := adaptTargetGroupToTargetGroupResource(
		*sdkTargetGroup,
		ctx,
		&response.Diagnostics,
	)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, targetGroup)...)
}

func (t *targetGroupResource) Delete(
	ctx context.Context,
	request resource.DeleteRequest,
	response *resource.DeleteResponse,
) {
	var state targetGroupResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	httpResponse, err := t.PubliccloudAPI.DeleteTargetGroup(
		ctx,
		state.ID.ValueString(),
	).Execute()

	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
	}
}

func NewTargetGroupResource() resource.Resource {
	return &targetGroupResource{
		ResourceAPI: utils.ResourceAPI{
			Name: "public_cloud_target_group",
		},
	}
}
