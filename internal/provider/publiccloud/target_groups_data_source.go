package publiccloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publiccloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ datasource.DataSourceWithConfigure = &targetGroupsDataSource{}
)

type targetGroupsDataSourceModel struct {
	ID           types.String                 `tfsdk:"id"`
	Name         types.String                 `tfsdk:"name"`
	Protocol     types.String                 `tfsdk:"protocol"`
	Port         types.Int32                  `tfsdk:"port"`
	Region       types.String                 `tfsdk:"region"`
	TargetGroups []targetGroupDataSourceModel `tfsdk:"target_groups"`
}

type targetGroupDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Protocol types.String `tfsdk:"protocol"`
	Port     types.Int32  `tfsdk:"port"`
	Region   types.String `tfsdk:"region"`
}

type targetGroupsDataSource struct {
	utils.DataSourceAPI
}

func (t *targetGroupsDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	response *datasource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Description: "Target group ID",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the target group",
			},
			"protocol": schema.StringAttribute{
				Optional:    true,
				Description: "Valid options are " + utils.StringTypeArrayToMarkdown(publiccloud.AllowedProtocolEnumValues),
				Validators: []validator.String{
					stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publiccloud.AllowedProtocolEnumValues)...),
				},
			},
			"port": schema.Int32Attribute{
				Optional:    true,
				Description: "The port of the target group",
				Validators: []validator.Int32{
					int32validator.Between(1, 65535),
				},
			},
			"region": schema.StringAttribute{
				Optional:    true,
				Description: "Region name. Valid options are " + utils.StringTypeArrayToMarkdown(publiccloud.AllowedRegionNameEnumValues),
				Validators: []validator.String{
					stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publiccloud.AllowedRegionNameEnumValues)...),
				},
			},
			"target_groups": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Target group ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the target group",
						},
						"protocol": schema.StringAttribute{
							Computed: true,
						},
						"port": schema.Int32Attribute{
							Computed:    true,
							Description: "The port of the target group",
						},
						"region": schema.StringAttribute{
							Computed:    true,
							Description: "Region name",
						},
					},
				},
			},
		},
	}
}

func (t *targetGroupsDataSource) Read(
	ctx context.Context,
	request datasource.ReadRequest,
	response *datasource.ReadResponse,
) {
	var config targetGroupsDataSourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	targetGroupsRequest := t.PubliccloudAPI.GetTargetGroupList(ctx)
	if !config.ID.IsNull() {
		targetGroupsRequest = targetGroupsRequest.Id(config.ID.ValueString())
	}
	if !config.Name.IsNull() {
		targetGroupsRequest = targetGroupsRequest.Name(config.Name.ValueString())
	}
	if !config.Protocol.IsNull() {
		targetGroupsRequest = targetGroupsRequest.Protocol(publiccloud.Protocol(config.Protocol.ValueString()))
	}
	if !config.Port.IsNull() {
		targetGroupsRequest = targetGroupsRequest.Port(config.Port.ValueInt32())
	}
	if !config.Region.IsNull() {
		targetGroupsRequest = targetGroupsRequest.Region(publiccloud.RegionName(config.Region.ValueString()))
	}
	var targetGroups []publiccloud.TargetGroup
	var offset *int32
	for {
		result, httpResponse, err := targetGroupsRequest.Execute()
		if err != nil {
			utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
			return
		}

		targetGroups = append(targetGroups, result.GetTargetGroups()...)

		metadata := result.GetMetadata()
		offset = utils.NewOffset(
			metadata.GetLimit(),
			metadata.GetOffset(),
			metadata.GetTotalCount(),
		)
		if offset == nil {
			break
		}

		targetGroupsRequest = targetGroupsRequest.Offset(*offset)
	}

	state := targetGroupsDataSourceModel{}
	for _, targetGroup := range targetGroups {
		state.TargetGroups = append(
			state.TargetGroups,
			targetGroupDataSourceModel{
				ID:       basetypes.NewStringValue(targetGroup.GetId()),
				Name:     basetypes.NewStringValue(targetGroup.GetName()),
				Protocol: basetypes.NewStringValue(string(targetGroup.GetProtocol())),
				Port:     basetypes.NewInt32Value(targetGroup.GetPort()),
				Region:   basetypes.NewStringValue(string(targetGroup.GetRegion())),
			},
		)
	}
	state.ID = config.ID
	state.Name = config.Name
	state.Protocol = config.Protocol
	state.Port = config.Port
	state.Region = config.Region

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func NewTargetGroupsDataSource() datasource.DataSource {
	return &targetGroupsDataSource{
		DataSourceAPI: utils.DataSourceAPI{
			Name: "public_cloud_target_groups",
		},
	}
}
