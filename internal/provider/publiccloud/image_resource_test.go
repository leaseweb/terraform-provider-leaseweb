package publiccloud

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publiccloud"
	"github.com/stretchr/testify/assert"
)

func Test_adaptImageDetailsToImageResource(t *testing.T) {
	state := publiccloud.IMAGESTATENAME_READY
	region := publiccloud.IMAGEREGION_EU_WEST_3

	sdkImageDetails := publiccloud.ImageDetails{
		Id:           "imageId",
		Name:         "name",
		Custom:       true,
		State:        *publiccloud.NewNullableImageStateName(&state),
		MarketApps:   []publiccloud.MarketAppId{publiccloud.MARKETAPPID_CPANEL_30},
		StorageTypes: []publiccloud.StorageType{publiccloud.STORAGETYPE_CENTRAL},
		Flavour:      "flavour",
		Region:       *publiccloud.NewNullableImageRegion(&region),
	}

	marketApps, _ := basetypes.NewListValueFrom(
		context.TODO(),
		types.StringType,
		[]string{"CPANEL_30"},
	)
	storageTypes, _ := basetypes.NewListValueFrom(
		context.TODO(),
		types.StringType,
		[]string{"CENTRAL"},
	)

	diags := diag.Diagnostics{}

	got := adaptImageDetailsToImageResource(
		context.TODO(),
		sdkImageDetails,
		&diags,
	)

	want := imageResourceModel{
		ID:           basetypes.NewStringValue("imageId"),
		Name:         basetypes.NewStringValue("name"),
		Custom:       basetypes.NewBoolValue(true),
		State:        basetypes.NewStringValue("READY"),
		MarketApps:   marketApps,
		StorageTypes: storageTypes,
		Flavour:      basetypes.NewStringValue("flavour"),
		Region:       basetypes.NewStringValue("eu-west-3"),
	}

	assert.False(t, diags.HasError())
	assert.Equal(t, want, *got)
}
