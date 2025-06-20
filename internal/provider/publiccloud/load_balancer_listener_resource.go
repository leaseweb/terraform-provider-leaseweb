package publiccloud

import (
	"context"
	"strings"

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
	"github.com/leaseweb/leaseweb-go-sdk/publiccloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.ResourceWithConfigure   = &loadBalancerListenerResource{}
	_ resource.ResourceWithImportState = &loadBalancerListenerResource{}
)

type loadBalancerListenerDefaultRuleResourceModel struct {
	TargetGroupID types.String `tfsdk:"target_group_id"`
}

func (l loadBalancerListenerDefaultRuleResourceModel) generateLoadBalancerListenerDefaultRule() publiccloud.LoadBalancerListenerDefaultRule {
	return *publiccloud.NewLoadBalancerListenerDefaultRule(l.TargetGroupID.ValueString())
}

func (l loadBalancerListenerDefaultRuleResourceModel) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"target_group_id": types.StringType,
	}
}

func adaptLoadBalancerListenerRuleToLoadBalancerListenerDefaultRuleResource(loadBalancerListenerRule publiccloud.LoadBalancerListenerRule) loadBalancerListenerDefaultRuleResourceModel {
	return loadBalancerListenerDefaultRuleResourceModel{
		TargetGroupID: basetypes.NewStringValue(loadBalancerListenerRule.GetTargetGroupId()),
	}
}

type loadBalancerListenerCertificateResourceModel struct {
	PrivateKey  types.String `tfsdk:"private_key"`
	Certificate types.String `tfsdk:"certificate"`
	Chain       types.String `tfsdk:"chain"`
}

func (l loadBalancerListenerCertificateResourceModel) generateSslCertificate() publiccloud.SslCertificate {
	sslCertificate := publiccloud.NewSslCertificate(
		l.PrivateKey.ValueString(),
		l.Certificate.ValueString(),
	)
	if !l.Chain.IsNull() && l.Chain.ValueString() != "" {
		sslCertificate.SetChain(l.Chain.ValueString())
	}

	return *sslCertificate
}

type loadBalancerListenerResourceModel struct {
	ListenerID     types.String `tfsdk:"listener_id"`
	LoadBalancerID types.String `tfsdk:"load_balancer_id"`
	Protocol       types.String `tfsdk:"protocol"`
	Port           types.Int32  `tfsdk:"port"`
	Certificate    types.Object `tfsdk:"certificate"`
	DefaultRule    types.Object `tfsdk:"default_rule"`
}

func adaptLoadBalancerListenerToLoadBalancerListenerResource(
	loadBalancerListener publiccloud.LoadBalancerListener,
	ctx context.Context,
	diags *diag.Diagnostics,
) *loadBalancerListenerResourceModel {
	listener := loadBalancerListenerResourceModel{
		ListenerID: basetypes.NewStringValue(loadBalancerListener.GetId()),
		Protocol:   basetypes.NewStringValue(string(loadBalancerListener.Protocol)),
		Port:       basetypes.NewInt32Value(loadBalancerListener.GetPort()),
	}

	if len(loadBalancerListener.Rules) > 0 {
		defaultRule := utils.AdaptSdkModelToResourceObject(
			loadBalancerListener.Rules[0],
			loadBalancerListenerDefaultRuleResourceModel{}.attributeTypes(),
			ctx,
			adaptLoadBalancerListenerRuleToLoadBalancerListenerDefaultRuleResource,
			diags,
		)
		if diags.HasError() {
			return nil
		}
		listener.DefaultRule = defaultRule
	}

	return &listener
}

type loadBalancerListenerResource struct {
	utils.ResourceAPI
}

func (l *loadBalancerListenerResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	response *resource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"load_balancer_id": schema.StringAttribute{
				Required:    true,
				Description: "Load balancer ID",
			},
			"listener_id": schema.StringAttribute{
				Computed:    true,
				Description: "Listener ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "Valid options are " + utils.StringTypeArrayToMarkdown(publiccloud.AllowedProtocolEnumValues),
				Validators: []validator.String{
					stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publiccloud.AllowedProtocolEnumValues)...),
				},
			},
			"port": schema.Int32Attribute{
				Required:    true,
				Description: "Port that the listener listens to",
				Validators: []validator.Int32{
					int32validator.Between(1, 65535),
				},
			},
			"certificate": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Required only if protocol is HTTPS",
				Attributes: map[string]schema.Attribute{
					"private_key": schema.StringAttribute{
						Optional:    true,
						Description: "Client Private Key. Required only if protocol is `HTTPS`",
						Sensitive:   true,
					},
					"certificate": schema.StringAttribute{
						Optional:    true,
						Description: "Client Certificate. Required only if protocol is `HTTPS`",
						Sensitive:   true,
					},
					"chain": schema.StringAttribute{
						Optional:    true,
						Description: "CA certificate. Not required, but can be added if protocol is `HTTPS`",
						Sensitive:   true,
					},
				},
			},
			"default_rule": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"target_group_id": schema.StringAttribute{
						Optional:    true,
						Description: "Client Private Key. Required only if protocol is `HTTPS`",
					},
				},
			},
		},
	}
}

func (l *loadBalancerListenerResource) Create(
	ctx context.Context,
	request resource.CreateRequest,
	response *resource.CreateResponse,
) {
	var plan loadBalancerListenerResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	defaultRule := loadBalancerListenerDefaultRuleResourceModel{}
	defaultRuleDiags := plan.DefaultRule.As(ctx, &defaultRule, basetypes.ObjectAsOptions{})
	if defaultRuleDiags != nil {
		response.Diagnostics.Append(defaultRuleDiags...)
		return
	}

	opts := publiccloud.NewLoadBalancerListenerCreateOpts(
		publiccloud.Protocol(plan.Protocol.ValueString()),
		plan.Port.ValueInt32(),
		defaultRule.generateLoadBalancerListenerDefaultRule(),
	)

	if !plan.Certificate.IsNull() {
		certificate := loadBalancerListenerCertificateResourceModel{}
		certificateDiags := plan.Certificate.As(ctx, &certificate, basetypes.ObjectAsOptions{})
		if certificateDiags != nil {
			response.Diagnostics.Append(certificateDiags...)
			return
		}

		opts.SetCertificate(certificate.generateSslCertificate())
	}

	loadBalancerListener, httpResponse, err := l.PubliccloudAPI.CreateLoadBalancerListener(
		ctx,
		plan.LoadBalancerID.ValueString(),
	).LoadBalancerListenerCreateOpts(*opts).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	state := adaptLoadBalancerListenerToLoadBalancerListenerResource(
		*loadBalancerListener,
		ctx,
		&response.Diagnostics,
	)
	if response.Diagnostics.HasError() {
		return
	}

	state.LoadBalancerID = plan.LoadBalancerID
	state.Certificate = plan.Certificate

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (l *loadBalancerListenerResource) Read(
	ctx context.Context,
	request resource.ReadRequest,
	response *resource.ReadResponse,
) {
	var state loadBalancerListenerResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	loadBalancerListenerDetails, httpResponse, err := l.PubliccloudAPI.GetLoadBalancerListener(
		ctx,
		state.LoadBalancerID.ValueString(),
		state.ListenerID.ValueString(),
	).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	newState := loadBalancerListenerResourceModel{
		ListenerID: basetypes.NewStringValue(loadBalancerListenerDetails.GetId()),
		Protocol:   basetypes.NewStringValue(string(loadBalancerListenerDetails.GetProtocol())),
		Port:       basetypes.NewInt32Value(loadBalancerListenerDetails.GetPort()),
	}
	if len(loadBalancerListenerDetails.SslCertificates) > 0 {
		certificate := utils.AdaptSdkModelToResourceObject(
			loadBalancerListenerDetails.SslCertificates[0],
			map[string]attr.Type{
				"private_key": types.StringType,
				"certificate": types.StringType,
				"chain":       types.StringType,
			},
			ctx,
			func(sslCertificate publiccloud.SslCertificate) loadBalancerListenerCertificateResourceModel {
				listener := loadBalancerListenerCertificateResourceModel{
					PrivateKey:  basetypes.NewStringValue(sslCertificate.GetPrivateKey()),
					Certificate: basetypes.NewStringValue(sslCertificate.GetCertificate()),
				}

				chain, _ := sslCertificate.GetChainOk()
				if chain != nil && *chain != "" {
					listener.Chain = basetypes.NewStringPointerValue(chain)
				}

				return listener
			},
			&response.Diagnostics,
		)
		if response.Diagnostics.HasError() {
			return
		}
		newState.Certificate = certificate
	}

	if len(loadBalancerListenerDetails.Rules) > 0 {
		defaultRule := utils.AdaptSdkModelToResourceObject(
			loadBalancerListenerDetails.Rules[0],
			loadBalancerListenerDefaultRuleResourceModel{}.attributeTypes(),
			ctx,
			adaptLoadBalancerListenerRuleToLoadBalancerListenerDefaultRuleResource,
			&response.Diagnostics,
		)
		if response.Diagnostics.HasError() {
			return
		}
		newState.DefaultRule = defaultRule
	}

	newState.LoadBalancerID = state.LoadBalancerID
	newState.ListenerID = state.ListenerID
	if newState.Certificate.IsNull() {
		newState.Certificate = state.Certificate
	}

	response.Diagnostics.Append(response.State.Set(ctx, newState)...)
}

func (l *loadBalancerListenerResource) Update(
	ctx context.Context,
	request resource.UpdateRequest,
	response *resource.UpdateResponse,
) {
	var plan loadBalancerListenerResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	opts := publiccloud.NewLoadBalancerListenerOpts()
	opts.SetProtocol(publiccloud.Protocol(plan.Protocol.ValueString()))
	opts.SetPort(plan.Port.ValueInt32())

	if !plan.Certificate.IsNull() {
		certificate := loadBalancerListenerCertificateResourceModel{}
		certificateDiags := plan.Certificate.As(
			ctx,
			&certificate,
			basetypes.ObjectAsOptions{},
		)
		if certificateDiags != nil {
			response.Diagnostics.Append(certificateDiags...)
			return
		}

		opts.SetCertificate(certificate.generateSslCertificate())
	}

	if !plan.DefaultRule.IsNull() {
		defaultRule := loadBalancerListenerDefaultRuleResourceModel{}
		defaultRuleDiags := plan.DefaultRule.As(
			ctx,
			&defaultRule,
			basetypes.ObjectAsOptions{},
		)
		if defaultRuleDiags != nil {
			response.Diagnostics.Append(defaultRuleDiags...)
			return
		}

		opts.SetDefaultRule(defaultRule.generateLoadBalancerListenerDefaultRule())
	}

	loadBalancerListener, httpResponse, err := l.PubliccloudAPI.
		UpdateLoadBalancerListener(
			ctx,
			plan.LoadBalancerID.ValueString(),
			plan.ListenerID.ValueString(),
		).
		LoadBalancerListenerOpts(*opts).
		Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	state := adaptLoadBalancerListenerToLoadBalancerListenerResource(
		*loadBalancerListener,
		ctx,
		&response.Diagnostics,
	)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Certificate.IsNull() {
		state.Certificate = plan.Certificate
	}
	state.LoadBalancerID = plan.LoadBalancerID

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (l *loadBalancerListenerResource) Delete(
	ctx context.Context,
	request resource.DeleteRequest,
	response *resource.DeleteResponse,
) {
	var state loadBalancerListenerResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	httpResponse, err := l.PubliccloudAPI.DeleteLoadBalancerListener(
		ctx,
		state.LoadBalancerID.ValueString(),
		state.ListenerID.ValueString(),
	).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
	}
}

func (l *loadBalancerListenerResource) ImportState(
	ctx context.Context,
	request resource.ImportStateRequest,
	response *resource.ImportStateResponse,
) {
	idParts := strings.Split(request.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		utils.UnexpectedImportIdentifierError(
			&response.Diagnostics,
			"load_balancer_id,listener_id",
			request.ID,
		)
		return
	}

	response.Diagnostics.Append(response.State.SetAttribute(
		ctx,
		path.Root("load_balancer_id"),
		idParts[0],
	)...)
	response.Diagnostics.Append(response.State.SetAttribute(
		ctx,
		path.Root("listener_id"),
		idParts[1],
	)...)
}

func NewLoadBalancerListenerResource() resource.Resource {
	return &loadBalancerListenerResource{
		ResourceAPI: utils.ResourceAPI{
			Name: "public_cloud_load_balancer_listener",
		},
	}
}
