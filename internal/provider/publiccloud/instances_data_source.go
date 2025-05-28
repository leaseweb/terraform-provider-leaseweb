package publiccloud

import (
	"context"
	"fmt"
	"net/http"
	"sort"

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
	_ datasource.DataSourceWithConfigure = &instancesDataSource{}
)

type instanceDetailsErr struct {
	err          error
	httpResponse *http.Response
}

type contractDataSourceModel struct {
	BillingFrequency types.Int32  `tfsdk:"billing_frequency"`
	EndsAt           types.String `tfsdk:"ends_at"`
	State            types.String `tfsdk:"state"`
	Term             types.Int32  `tfsdk:"term"`
	Type             types.String `tfsdk:"type"`
}

func adaptContractToContractDataSource(contract publiccloud.InstanceContract) contractDataSourceModel {
	return contractDataSourceModel{
		BillingFrequency: basetypes.NewInt32Value(int32(contract.GetBillingFrequency())),
		EndsAt:           utils.AdaptNullableTimeToStringValue(contract.EndsAt.Get()),
		State:            basetypes.NewStringValue(string(contract.GetState())),
		Term:             basetypes.NewInt32Value(int32(contract.GetTerm())),
		Type:             basetypes.NewStringValue(string(contract.GetType())),
	}
}

type instanceDataSourceModel struct {
	Contract            contractDataSourceModel `tfsdk:"contract"`
	ID                  types.String            `tfsdk:"id"`
	Image               imageModelDataSource    `tfsdk:"image"`
	IPs                 []ipDataSourceModel     `tfsdk:"ips"`
	ISO                 *isoDataSourceModel     `tfsdk:"iso"`
	MarketAppID         types.String            `tfsdk:"market_app_id"`
	Reference           types.String            `tfsdk:"reference"`
	Region              types.String            `tfsdk:"region"`
	RootDiskSize        types.Int32             `tfsdk:"root_disk_size"`
	RootDiskStorageType types.String            `tfsdk:"root_disk_storage_type"`
	State               types.String            `tfsdk:"state"`
	Type                types.String            `tfsdk:"type"`
}

type ipDataSourceModel struct {
	IP types.String `tfsdk:"ip"`
}

type instancesDataSourceModel struct {
	Instances []instanceDataSourceModel `tfsdk:"instances"`
}

func NewInstancesDataSource() datasource.DataSource {
	return &instancesDataSource{
		DataSourceAPI: utils.DataSourceAPI{
			Name: "public_cloud_instances",
		},
	}
}

type instancesDataSource struct {
	utils.DataSourceAPI
}

func (d *instancesDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var instances []publiccloud.Instance
	var offset *int32

	// Get instances
	request := d.PubliccloudAPI.GetInstanceList(ctx)
	for {
		result, httpResponse, err := request.Execute()
		if err != nil {
			utils.SdkError(ctx, &resp.Diagnostics, err, httpResponse)
			return
		}

		instances = append(instances, result.Instances...)

		metadata := result.GetMetadata()
		offset = utils.NewOffset(
			metadata.GetLimit(),
			metadata.GetOffset(),
			metadata.GetTotalCount(),
		)
		if offset == nil {
			break
		}
		request = request.Offset(*offset)
	}

	//Get images once
	images := getAllImages(ctx, d.PubliccloudAPI, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get instanceDetails for each instance
	var instanceDetailsList []publiccloud.InstanceDetails
	resultChan := make(chan publiccloud.InstanceDetails)
	errorChan := make(chan instanceDetailsErr)
	for _, instance := range instances {
		go func(id string) {
			instanceDetails, httpResponse, err := d.PubliccloudAPI.GetInstance(
				ctx,
				id,
			).Execute()
			if err != nil {
				errorChan <- instanceDetailsErr{
					err:          err,
					httpResponse: httpResponse,
				}
				return
			}
			resultChan <- *instanceDetails
		}(instance.Id)
	}
	for i := 0; i < len(instances); i++ {
		select {
		case err := <-errorChan:
			utils.SdkError(ctx, &resp.Diagnostics, err.err, err.httpResponse)
			return
		case res := <-resultChan:
			instanceDetailsList = append(instanceDetailsList, res)
		}
	}

	var state instancesDataSourceModel

	sort.Slice(instanceDetailsList, func(i, j int) bool {
		return instanceDetailsList[i].Id < instanceDetailsList[j].Id
	})
	for _, instanceDetails := range instanceDetailsList {
		var ips []ipDataSourceModel
		for _, ip := range instanceDetails.Ips {
			ips = append(
				ips,
				ipDataSourceModel{
					IP: basetypes.NewStringValue(ip.GetIp()),
				},
			)
		}

		instance := instanceDataSourceModel{
			Contract:            adaptContractToContractDataSource(instanceDetails.GetContract()),
			ID:                  basetypes.NewStringValue(instanceDetails.GetId()),
			IPs:                 ips,
			MarketAppID:         basetypes.NewStringPointerValue(instanceDetails.MarketAppId.Get()),
			RootDiskSize:        basetypes.NewInt32Value(instanceDetails.GetRootDiskSize()),
			RootDiskStorageType: basetypes.NewStringValue(string(instanceDetails.GetRootDiskStorageType())),
			Region:              basetypes.NewStringValue(string(instanceDetails.GetRegion())),
			Reference:           basetypes.NewStringPointerValue(instanceDetails.Reference.Get()),
			State:               basetypes.NewStringValue(string(instanceDetails.GetState())),
			Type:                basetypes.NewStringValue(string(instanceDetails.GetType())),
		}
		imageDetails := images.findById(instanceDetails.Image.Id)
		if imageDetails == nil {
			utils.GeneralError(
				&resp.Diagnostics,
				ctx,
				fmt.Errorf("imageDetails %s not found", instance.Image.ID),
			)
			return
		}
		instance.Image = adaptImageDetailsToImageDataSource(*imageDetails)

		sdkIso, _ := instanceDetails.GetIsoOk()
		if sdkIso != nil {
			iso := adaptIsoToISODataSource(*sdkIso)
			instance.ISO = &iso
		}

		state.Instances = append(state.Instances, instance)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (d *instancesDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	// 0 has to be prepended manually as it's a valid option.
	billingFrequencies := utils.NewIntMarkdownList(
		append(
			[]publiccloud.BillingFrequency{0},
			publiccloud.AllowedBillingFrequencyEnumValues...,
		),
	)
	contractTerms := utils.NewIntMarkdownList(publiccloud.AllowedContractTermEnumValues)

	resp.Schema = schema.Schema{
		Description: utils.BetaDescription,
		Attributes: map[string]schema.Attribute{
			"instances": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The instance unique identifier",
						},
						"iso": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:    true,
									Description: "The ISO ID.",
								},
								"name": schema.StringAttribute{
									Computed: true,
								},
							},
						},
						"contract": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"billing_frequency": schema.Int32Attribute{
									Computed:    true,
									Description: "The billing frequency (in months). Valid options are " + billingFrequencies.Markdown(),
									Validators: []validator.Int32{
										int32validator.OneOf(billingFrequencies.ToInt32()...),
									},
								},
								"term": schema.Int32Attribute{
									Computed:    true,
									Description: "Contract term (in months). Used only when type is *MONTHLY*. Valid options are " + contractTerms.Markdown(),
									Validators: []validator.Int32{
										int32validator.OneOf(contractTerms.ToInt32()...),
									},
								},
								"type": schema.StringAttribute{
									Computed:    true,
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
						"market_app_id": schema.StringAttribute{
							Computed:    true,
							Description: "Market App ID",
						},
						"reference": schema.StringAttribute{
							Computed:    true,
							Description: "The identifying name set to the instance",
						},
						"region": schema.StringAttribute{
							Computed: true,
						},
						"root_disk_size": schema.Int32Attribute{
							Computed:    true,
							Description: "The root disk's size in GB. Must be at least 5 GB for Linux and FreeBSD instances and 50 GB for Windows instances",
						},
						"root_disk_storage_type": schema.StringAttribute{
							Computed:    true,
							Description: "The root disk's storage type",
						},
						"image": schema.SingleNestedAttribute{
							Computed:   true,
							Attributes: imageSchemaAttributes(),
						},
						"ips": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"ip": schema.StringAttribute{Computed: true},
								},
							},
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "The instance's current state",
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
