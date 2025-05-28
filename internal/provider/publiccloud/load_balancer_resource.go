package publiccloud

import (
	"context"
	"fmt"

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
	_ resource.ResourceWithConfigure   = &loadBalancerResource{}
	_ resource.ResourceWithImportState = &loadBalancerResource{}
)

type loadBalancerIPResourceModel struct {
	ReverseLookup  types.String `tfsdk:"reverse_lookup"`
	LoadBalancerID types.String `tfsdk:"load_balancer_id"`
	IP             types.String `tfsdk:"ip"`
}

type loadBalancerResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Region    types.String `tfsdk:"region"`
	Type      types.String `tfsdk:"type"`
	Reference types.String `tfsdk:"reference"`
	Contract  types.Object `tfsdk:"contract"`
	IPs       types.List   `tfsdk:"ips"`
}

func adaptLoadBalancerDetailsToLoadBalancerResource(
	loadBalancerDetails publiccloud.LoadBalancerDetails,
	ctx context.Context,
	diags *diag.Diagnostics,
) *loadBalancerResourceModel {
	loadBalancer := loadBalancerResourceModel{
		ID:        basetypes.NewStringValue(loadBalancerDetails.GetId()),
		Region:    basetypes.NewStringValue(string(loadBalancerDetails.GetRegion())),
		Type:      basetypes.NewStringValue(string(loadBalancerDetails.GetType())),
		Reference: basetypes.NewStringPointerValue(loadBalancerDetails.Reference.Get()),
	}

	contract := utils.AdaptSdkModelToResourceObject(
		loadBalancerDetails.Contract,
		contractResourceModel{}.attributeTypes(),
		ctx,
		adaptContractToContractResource,
		diags,
	)
	if diags.HasError() {
		return nil
	}
	loadBalancer.Contract = contract

	ips := utils.AdaptSdkModelsToListValue(
		loadBalancerDetails.Ips,
		map[string]attr.Type{
			"reverse_lookup":   types.StringType,
			"load_balancer_id": types.StringType,
			"ip":               types.StringType,
		},
		ctx,
		adaptIpDetailsToLoadBalancerIPResource,
		diags,
	)
	if diags.HasError() {
		return nil
	}
	loadBalancer.IPs = ips

	return &loadBalancer
}

func adaptIpDetailsToLoadBalancerIPResource(ipDetails publiccloud.IpDetails) loadBalancerIPResourceModel {
	reverseLookup, _ := ipDetails.GetReverseLookupOk()
	return loadBalancerIPResourceModel{
		ReverseLookup: basetypes.NewStringPointerValue(reverseLookup),
		IP:            basetypes.NewStringValue(ipDetails.GetIp()),
	}
}

type loadBalancerResource struct {
	utils.ResourceAPI
}

func (l *loadBalancerResource) ImportState(
	ctx context.Context,
	request resource.ImportStateRequest,
	response *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (l *loadBalancerResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	response *resource.SchemaResponse,
) {
	warningError := "**WARNING!** Changing this value once running will cause this loadbalancer to be destroyed and a new one to be created."

	contractTerms := utils.NewIntMarkdownList(publiccloud.AllowedContractTermEnumValues)
	// 0 has to be prepended manually as it's a valid option.
	billingFrequencies := utils.NewIntMarkdownList(
		append(
			[]publiccloud.BillingFrequency{0},
			publiccloud.AllowedBillingFrequencyEnumValues...,
		),
	)

	response.Schema = schema.Schema{
		Description: utils.BetaDescription,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The load balancer unique identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"reference": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "An identifying name you can refer to the load balancer",
			},
			"ips": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip": schema.StringAttribute{Computed: true},
						"load_balancer_id": schema.StringAttribute{
							Computed: true,
						},
						"reverse_lookup": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"contract": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"billing_frequency": schema.Int32Attribute{
						Required:    true,
						Description: "The billing frequency (in months). Valid options are " + billingFrequencies.Markdown(),
						Validators: []validator.Int32{
							int32validator.OneOf(billingFrequencies.ToInt32()...),
						},
					},
					"term": schema.Int32Attribute{
						Required:    true,
						Description: "Contract term (in months). Used only when type is *MONTHLY*. Valid options are " + contractTerms.Markdown(),
						Validators: []validator.Int32{
							int32validator.OneOf(contractTerms.ToInt32()...),
						},
					},
					"type": schema.StringAttribute{
						Required:    true,
						Description: "Select *HOURLY* for billing based on hourly usage, else *MONTHLY* for billing per month usage",
						Validators: []validator.String{
							stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publiccloud.AllowedContractTypeEnumValues)...),
						},
					},
					"ends_at": schema.StringAttribute{Computed: true},
					"state": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"region": schema.StringAttribute{
				Required: true,
				Description: fmt.Sprintf(
					"%s Valid options are %s",
					warningError,
					utils.StringTypeArrayToMarkdown(publiccloud.AllowedRegionNameEnumValues),
				),
				Validators: []validator.String{
					stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publiccloud.AllowedRegionNameEnumValues)...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required: true,
				Description: fmt.Sprintf(
					"%s Valid options are %s",
					warningError,
					utils.StringTypeArrayToMarkdown(publiccloud.AllowedTypeNameEnumValues),
				),
				Validators: []validator.String{
					stringvalidator.AlsoRequires(
						path.Expressions{path.MatchRoot("region")}...,
					),
					stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publiccloud.AllowedTypeNameEnumValues)...),
				},
			},
		},
	}
}

func (l *loadBalancerResource) Create(
	ctx context.Context,
	request resource.CreateRequest,
	response *resource.CreateResponse,
) {
	var plan loadBalancerResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	contract := contractResourceModel{}
	contractDiags := plan.Contract.As(ctx, &contract, basetypes.ObjectAsOptions{})
	if contractDiags != nil {
		response.Diagnostics.Append(contractDiags...)
		return
	}
	opts := publiccloud.NewLaunchLoadBalancerOpts(
		publiccloud.RegionName(plan.Region.ValueString()),
		publiccloud.TypeName(plan.Type.ValueString()),
		publiccloud.ContractType(contract.Type.ValueString()),
	)
	opts.Reference = utils.AdaptStringPointerValueToNullableString(plan.Reference)
	opts.SetContractTerm(publiccloud.ContractTerm(contract.Term.ValueInt32()))
	opts.SetBillingFrequency(publiccloud.BillingFrequency(contract.BillingFrequency.ValueInt32()))

	loadBalancer, httpResponse, err := l.PubliccloudAPI.LaunchLoadBalancer(ctx).
		LaunchLoadBalancerOpts(*opts).
		Execute()

	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	state := adaptLoadBalancerDetailsToLoadBalancerResource(
		*loadBalancer,
		ctx,
		&response.Diagnostics,
	)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (l *loadBalancerResource) Read(
	ctx context.Context,
	request resource.ReadRequest,
	response *resource.ReadResponse,
) {
	var state loadBalancerResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	loadBalancerDetails, httpResponse, err := l.PubliccloudAPI.
		GetLoadBalancer(ctx, state.ID.ValueString()).
		Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	newState := adaptLoadBalancerDetailsToLoadBalancerResource(
		*loadBalancerDetails,
		ctx,
		&response.Diagnostics,
	)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, newState)...)
}

func (l *loadBalancerResource) Update(
	ctx context.Context,
	request resource.UpdateRequest,
	response *resource.UpdateResponse,
) {
	var plan loadBalancerResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	opts := publiccloud.NewUpdateLoadBalancerOpts()
	opts.Reference = utils.AdaptStringPointerValueToNullableString(plan.Reference)
	if plan.Type.ValueString() != "" {
		opts.SetType(publiccloud.TypeName(plan.Type.ValueString()))
	}

	loadBalancerDetails, httpResponse, err := l.PubliccloudAPI.
		UpdateLoadBalancer(ctx, plan.ID.ValueString()).
		UpdateLoadBalancerOpts(*opts).
		Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}
	state := adaptLoadBalancerDetailsToLoadBalancerResource(
		*loadBalancerDetails,
		ctx,
		&response.Diagnostics,
	)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (l *loadBalancerResource) Delete(
	ctx context.Context,
	request resource.DeleteRequest,
	response *resource.DeleteResponse,
) {
	var state loadBalancerResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	opts := publiccloud.NewTerminateLoadBalancerOpts("CANCEL_OTHER")
	opts.SetReason("Terraform")

	httpResponse, err := l.PubliccloudAPI.TerminateLoadBalancer(
		ctx,
		state.ID.ValueString(),
	).TerminateLoadBalancerOpts(*opts).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
	}
}

func NewLoadBalancerResource() resource.Resource {
	return &loadBalancerResource{
		ResourceAPI: utils.ResourceAPI{
			Name: "public_cloud_load_balancer",
		},
	}
}
