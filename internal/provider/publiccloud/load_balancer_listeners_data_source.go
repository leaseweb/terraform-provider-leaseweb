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
	_ datasource.DataSourceWithConfigure = &loadBalancerListenersDataSource{}
)

type loadBalancerListenersDataSourceModel struct {
	LoadBalancerID types.String                          `tfsdk:"load_balancer_id"`
	Listeners      []loadBalancerListenerDataSourceModel `tfsdk:"listeners"`
}

type loadBalancerListenerDataSourceModel struct {
	ID types.String `tfsdk:"id"`
}

type loadBalancerListenersDataSource struct {
	utils.DataSourceAPI
}

func (l *loadBalancerListenersDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	response *datasource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"load_balancer_id": schema.StringAttribute{
				Required:    true,
				Description: "Load balancer ID",
			},
			"listeners": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The listener unique identifier",
						},
					},
				},
			},
		},
	}
}

func (l *loadBalancerListenersDataSource) Read(
	ctx context.Context,
	request datasource.ReadRequest,
	response *datasource.ReadResponse,
) {
	var config loadBalancerListenersDataSourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	var loadBalancerListeners []publiccloud.LoadBalancerListener
	var offset *int32
	loadBalancerListenerRequest := l.PubliccloudAPI.GetLoadBalancerListenerList(ctx, config.LoadBalancerID.ValueString())
	for {
		result, httpResponse, err := loadBalancerListenerRequest.Execute()
		if err != nil {
			utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
			return
		}

		loadBalancerListeners = append(loadBalancerListeners, result.GetListeners()...)

		metadata := result.GetMetadata()

		offset = utils.NewOffset(
			metadata.GetLimit(),
			metadata.GetOffset(),
			metadata.GetTotalCount(),
		)

		if offset == nil {
			break
		}

		loadBalancerListenerRequest = loadBalancerListenerRequest.Offset(*offset)
	}

	var state loadBalancerListenersDataSourceModel
	for _, loadBalancerListener := range loadBalancerListeners {
		listener := loadBalancerListenerDataSourceModel{
			ID: basetypes.NewStringValue(loadBalancerListener.GetId()),
		}
		state.Listeners = append(state.Listeners, listener)
	}
	state.LoadBalancerID = basetypes.NewStringValue(config.LoadBalancerID.ValueString())

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func NewLoadBalancerListenersDataSource() datasource.DataSource {
	return &loadBalancerListenersDataSource{
		DataSourceAPI: utils.DataSourceAPI{
			Name: "public_cloud_load_balancer_listeners",
		},
	}
}
