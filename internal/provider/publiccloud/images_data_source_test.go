package publiccloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publiccloud"
	"github.com/stretchr/testify/assert"
)

func Test_imageDetailsList_findById(t *testing.T) {
	t.Run("image returned when found", func(t *testing.T) {
		image := publiccloud.ImageDetails{Id: "id"}
		list := imageDetailsList{image}
		got := list.findById("id")

		assert.Equal(t, image, *got)
	})

	t.Run("returns nil when nothing is found", func(t *testing.T) {
		image := publiccloud.ImageDetails{Id: "id"}
		list := imageDetailsList{image}
		got := list.findById("tralala")

		assert.Nil(t, got)
	})
}

func Test_adaptImageDetailsToImageDataSource(t *testing.T) {
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

	want := imageModelDataSource{
		ID:           basetypes.NewStringValue("imageId"),
		Name:         basetypes.NewStringValue("name"),
		Custom:       basetypes.NewBoolValue(true),
		State:        basetypes.NewStringValue("READY"),
		MarketApps:   []string{"CPANEL_30"},
		StorageTypes: []string{"CENTRAL"},
		Flavour:      basetypes.NewStringValue("flavour"),
		Region:       basetypes.NewStringValue("eu-west-3"),
	}
	got := adaptImageDetailsToImageDataSource(sdkImageDetails)

	assert.Equal(t, want, got)
}
