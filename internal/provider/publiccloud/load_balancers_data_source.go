package publiccloud

import (
	"context"
	"net/http"
	"sort"

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

type loadBalancerDetailsErr struct {
	err          error
	httpResponse *http.Response
}

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

	// Get loadBalancerDetails for each loadbalancer
	var loadBalancerDetailsList []publiccloud.LoadBalancerDetails
	resultChan := make(chan publiccloud.LoadBalancerDetails)
	errorChan := make(chan loadBalancerDetailsErr)
	for _, loadBalancer := range loadBalancers {
		go func(id string) {
			loadBalancerDetails, httpResponse, err := l.PubliccloudAPI.GetLoadBalancer(
				ctx,
				id,
			).Execute()
			if err != nil {
				errorChan <- loadBalancerDetailsErr{
					err:          err,
					httpResponse: httpResponse,
				}
				return
			}
			resultChan <- *loadBalancerDetails
		}(loadBalancer.Id)
	}
	for i := 0; i < len(loadBalancers); i++ {
		select {
		case err := <-errorChan:
			utils.SdkError(ctx, &response.Diagnostics, err.err, err.httpResponse)
			return
		case res := <-resultChan:
			loadBalancerDetailsList = append(loadBalancerDetailsList, res)
		}
	}

	var state loadBalancersDataSourceModel

	sort.Slice(loadBalancerDetailsList, func(i, j int) bool {
		return loadBalancerDetailsList[i].Id < loadBalancerDetailsList[j].Id
	})

	for _, sdkLoadBalancer := range loadBalancerDetailsList {
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
