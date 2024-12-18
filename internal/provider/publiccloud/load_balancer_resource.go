package publiccloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/v3/publiccloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.ResourceWithConfigure   = &loadBalancerResource{}
	_ resource.ResourceWithImportState = &loadBalancerResource{}
)

type loadBalancerResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Region    types.String `tfsdk:"region"`
	Type      types.String `tfsdk:"type"`
	Reference types.String `tfsdk:"reference"`
	Contract  types.Object `tfsdk:"contract"`
}

func (l *loadBalancerResourceModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":        types.StringType,
		"region":    types.StringType,
		"type":      types.StringType,
		"reference": types.StringType,
		"contract": types.ObjectType{
			AttrTypes: contractResourceModel{}.attributeTypes(),
		},
	}
}

func (l *loadBalancerResourceModel) GetLaunchLoadBalancerOpts(ctx context.Context) (
	*publiccloud.LaunchLoadBalancerOpts,
	error,
) {
	contract := contractResourceModel{}
	contractDiags := l.Contract.As(ctx, &contract, basetypes.ObjectAsOptions{})
	if contractDiags != nil {
		return nil, utils.ReturnError("GetLaunchLoadBalancerOpts", contractDiags)
	}

	contractType, err := publiccloud.NewContractTypeFromValue(contract.Type.ValueString())
	if err != nil {
		return nil, err
	}

	contractTerm, err := publiccloud.NewContractTermFromValue(contract.Term.ValueInt32())
	if err != nil {
		return nil, err
	}

	billingFrequency, err := publiccloud.NewBillingFrequencyFromValue(contract.BillingFrequency.ValueInt32())
	if err != nil {
		return nil, err
	}

	regionName, err := publiccloud.NewRegionNameFromValue(l.Region.ValueString())
	if err != nil {
		return nil, err
	}

	typeName, err := publiccloud.NewTypeNameFromValue(l.Type.ValueString())
	if err != nil {
		return nil, err
	}

	opts := publiccloud.NewLaunchLoadBalancerOpts(
		*regionName,
		*typeName,
		*contractType,
		*contractTerm,
		*billingFrequency,
	)

	opts.Reference = utils.AdaptStringPointerValueToNullableString(l.Reference)

	return opts, nil
}

func (l *loadBalancerResourceModel) GetUpdateLoadBalancerOpts() (
	*publiccloud.UpdateLoadBalancerOpts,
	error,
) {
	opts := publiccloud.NewUpdateLoadBalancerOpts()
	opts.Reference = utils.AdaptStringPointerValueToNullableString(l.Reference)

	if l.Type.ValueString() != "" {
		instanceType, err := publiccloud.NewTypeNameFromValue(l.Type.ValueString())
		if err != nil {
			return nil, fmt.Errorf("GetUpdateLoadBalancerOpts: %w", err)
		}
		opts.Type = instanceType
	}

	return opts, nil
}

func adaptLoadBalancerDetailsToLoadBalancerResource(
	loadBalancerDetails publiccloud.LoadBalancerDetails,
	ctx context.Context,
) (*loadBalancerResourceModel, error) {
	loadBalancer := loadBalancerResourceModel{
		ID:        basetypes.NewStringValue(loadBalancerDetails.GetId()),
		Region:    basetypes.NewStringValue(string(loadBalancerDetails.GetRegion())),
		Type:      basetypes.NewStringValue(string(loadBalancerDetails.GetType())),
		Reference: basetypes.NewStringPointerValue(loadBalancerDetails.Reference.Get()),
	}

	contract, err := utils.AdaptSdkModelToResourceObject(
		loadBalancerDetails.Contract,
		contractResourceModel{}.attributeTypes(),
		ctx,
		adaptContractToContractResource,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptLoadBalancerDetailsToLoadBalancerResource: %w", err)
	}
	loadBalancer.Contract = contract

	return &loadBalancer, nil
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

	opts, err := plan.GetLaunchLoadBalancerOpts(ctx)
	if err != nil {
		utils.GeneralError(&response.Diagnostics, ctx, err)
		return
	}

	loadBalancer, httpResponse, err := l.PubliccloudAPI.LaunchLoadBalancer(ctx).
		LaunchLoadBalancerOpts(*opts).
		Execute()

	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	state, err := adaptLoadBalancerDetailsToLoadBalancerResource(
		*loadBalancer,
		ctx,
	)
	if err != nil {
		utils.GeneralError(&response.Diagnostics, ctx, err)
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

	newState, err := adaptLoadBalancerDetailsToLoadBalancerResource(*loadBalancerDetails, ctx)
	if err != nil {
		utils.GeneralError(&response.Diagnostics, ctx, err)
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

	opts, err := plan.GetUpdateLoadBalancerOpts()
	if err != nil {
		utils.GeneralError(&response.Diagnostics, ctx, err)
		return
	}

	loadBalancerDetails, httpResponse, err := l.PubliccloudAPI.
		UpdateLoadBalancer(ctx, plan.ID.ValueString()).
		UpdateLoadBalancerOpts(*opts).
		Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}
	state, err := adaptLoadBalancerDetailsToLoadBalancerResource(
		*loadBalancerDetails,
		ctx,
	)
	if err != nil {
		utils.GeneralError(&response.Diagnostics, ctx, err)
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

	httpResponse, err := l.PubliccloudAPI.TerminateLoadBalancer(
		ctx,
		state.ID.ValueString(),
	).Execute()
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
