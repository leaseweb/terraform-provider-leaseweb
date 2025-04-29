package publiccloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publiccloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ datasource.DataSourceWithConfigure = &loadBalancersDataSource{}
)

type loadBalancerDataSourceModel struct {
	ID        types.String            `tfsdk:"id"`
	IPs       []ipDataSourceModel     `tfsdk:"ips"`
	Reference types.String            `tfsdk:"reference"`
	Contract  contractDataSourceModel `tfsdk:"contract"`
	State     types.String            `tfsdk:"state"`
	Region    types.String            `tfsdk:"region"`
	Type      types.String            `tfsdk:"type"`
}

type loadBalancersDataSourceModel struct {
	LoadBalancers []loadBalancerDataSourceModel `tfsdk:"load_balancers"`
}

type loadBalancersDataSource struct {
	utils.DataSourceAPI
}

func (l *loadBalancersDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	response *datasource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		Description: utils.BetaDescription,
		Attributes: map[string]schema.Attribute{
			"load_balancers": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The load balancer unique identifier",
						},
						"ips": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"ip": schema.StringAttribute{Computed: true},
								},
							},
						},
						"reference": schema.StringAttribute{
							Computed: true,
						},
						"contract": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"billing_frequency": schema.Int32Attribute{
									Computed:    true,
									Description: "The billing frequency (in months)",
								},
								"term": schema.Int32Attribute{
									Computed:    true,
									Description: "Contract term (in months)",
								},
								"type": schema.StringAttribute{
									Computed: true,
								},
								"ends_at": schema.StringAttribute{Computed: true},
								"state": schema.StringAttribute{
									Computed: true,
								},
							},
						},
						"state": schema.StringAttribute{
							Computed: true,
						},
						"region": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (l *loadBalancersDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	response *datasource.ReadResponse,
) {
	var loadBalancers []publiccloud.LoadBalancer
	var offset *int32

	loadBalancerRequest := l.PubliccloudAPI.GetLoadBalancerList(ctx)
	for {
		result, httpResponse, err := loadBalancerRequest.Execute()
		if err != nil {
			utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
			return
		}

		loadBalancers = append(loadBalancers, result.GetLoadBalancers()...)

		metadata := result.GetMetadata()

		offset = utils.NewOffset(
			metadata.GetLimit(),
			metadata.GetOffset(),
			metadata.GetTotalCount(),
		)

		if offset == nil {
			break
		}

		loadBalancerRequest = loadBalancerRequest.Offset(*offset)
	}

	var state loadBalancersDataSourceModel
	for _, sdkLoadBalancer := range loadBalancers {
		var ips []ipDataSourceModel
		for _, ip := range sdkLoadBalancer.Ips {
			ips = append(ips, ipDataSourceModel{IP: basetypes.NewStringValue(ip.GetIp())})
		}

		loadBalancer := loadBalancerDataSourceModel{
			ID:        basetypes.NewStringValue(sdkLoadBalancer.GetId()),
			IPs:       ips,
			Reference: basetypes.NewStringPointerValue(sdkLoadBalancer.Reference.Get()),
			Contract:  adaptContractToContractDataSource(sdkLoadBalancer.GetContract()),
			State:     basetypes.NewStringValue(string(sdkLoadBalancer.GetState())),
			Region:    basetypes.NewStringValue(string(sdkLoadBalancer.GetRegion())),
			Type:      basetypes.NewStringValue(string(sdkLoadBalancer.GetType())),
		}
		state.LoadBalancers = append(state.LoadBalancers, loadBalancer)
	}

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func NewLoadBalancersDataSource() datasource.DataSource {
	return &loadBalancersDataSource{
		DataSourceAPI: utils.DataSourceAPI{
			Name: "public_cloud_load_balancers",
		},
	}
}
