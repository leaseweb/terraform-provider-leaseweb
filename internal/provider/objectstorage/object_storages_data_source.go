package objectstorage

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/objectstorage"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ datasource.DataSourceWithConfigure = &objectStoragesDataSource{}
)

type objectStoragesDataSourceModel struct {
	ObjectStorages []objectStorageDataSourceModel `tfsdk:"object_storages"`
}

type objectStorageDataSourceModel struct {
	ContractID  types.String `tfsdk:"contract_id"`
	CustomerID  types.String `tfsdk:"customer_id"`
	Description types.String `tfsdk:"description"`
	ID          types.String `tfsdk:"id"`
	Quantity    types.String `tfsdk:"quantity"`
	RegionURL   types.String `tfsdk:"region_url"`
	SalesOrgID  types.String `tfsdk:"sales_org_id"`
}

func adaptObjectStorageListItemToObjectStorageDataSource(
	sdkObjectStorage objectstorage.ObjectStorageListItem,
) objectStorageDataSourceModel {
	return objectStorageDataSourceModel{
		ContractID:  basetypes.NewStringValue(sdkObjectStorage.GetContractId()),
		CustomerID:  basetypes.NewStringValue(sdkObjectStorage.GetCustomerId()),
		Description: basetypes.NewStringValue(sdkObjectStorage.GetDescription()),
		ID:          basetypes.NewStringValue(sdkObjectStorage.GetId()),
		Quantity:    basetypes.NewStringValue(sdkObjectStorage.GetQuantity()),
		RegionURL:   basetypes.NewStringValue(sdkObjectStorage.GetRegionUrl()),
		SalesOrgID:  basetypes.NewStringValue(sdkObjectStorage.GetSalesOrgId()),
	}
}

type objectStoragesDataSource struct {
	utils.DataSourceAPI
}

func (o objectStoragesDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	response *datasource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		Description: "List all object storages",
		Attributes: map[string]schema.Attribute{
			"object_storages": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"contract_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the contract the object storage belongs to",
						},
						"customer_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the customer the object storage belongs to",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "Description of the object storage",
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Object storage ID",
						},
						"quantity": schema.StringAttribute{
							Computed:    true,
							Description: "Quantity of the object storage",
						},
						"region_url": schema.StringAttribute{
							Computed:    true,
							Description: "The region specific endpoint URL for this object storage",
						},
						"sales_org_id": schema.StringAttribute{
							Computed:    true,
							Description: "Sales organization ID of the object storage",
						},
					},
				},
			},
		},
	}
}

func (o objectStoragesDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	response *datasource.ReadResponse,
) {
	var objectStorages []objectstorage.ObjectStorageListItem
	var offset *int32
	var state objectStoragesDataSourceModel

	objectStorageRequest := o.ObjectstorageAPI.GetObjectStorageList(ctx)

	for {
		result, httpResponse, err := objectStorageRequest.Execute()
		if err != nil {
			utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
			return
		}

		objectStorages = append(objectStorages, result.GetObjectStorages()...)

		metadata := result.GetMetadata()

		offset = utils.NewOffset(
			metadata.GetLimit(),
			metadata.GetOffset(),
			metadata.GetTotalCount(),
		)

		if offset == nil {
			break
		}

		objectStorageRequest = objectStorageRequest.Offset(*offset)
	}

	for _, sdkObjectStorage := range objectStorages {
		state.ObjectStorages = append(
			state.ObjectStorages,
			adaptObjectStorageListItemToObjectStorageDataSource(sdkObjectStorage),
		)
	}

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func NewObjectStoragesDataSource() datasource.DataSource {
	return &objectStoragesDataSource{
		DataSourceAPI: utils.DataSourceAPI{
			Name: "object_storages",
		},
	}
}
