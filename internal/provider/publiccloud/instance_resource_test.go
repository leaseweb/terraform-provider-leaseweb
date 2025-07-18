package publiccloud

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publiccloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_adaptContractToContractResource(t *testing.T) {
	endsAt, _ := time.Parse("2006-01-02 15:04:05", "2023-12-14 17:09:47")
	sdkContract := publiccloud.InstanceContractDetails{
		BillingFrequency: publiccloud.BILLINGFREQUENCY__1,
		Term:             publiccloud.CONTRACTTERM__3,
		Type:             publiccloud.CONTRACTTYPE_HOURLY,
		EndsAt:           *publiccloud.NewNullableTime(&endsAt),
		State:            publiccloud.CONTRACTSTATE_ACTIVE,
	}

	want := contractResourceModel{
		BillingFrequency: basetypes.NewInt32Value(1),
		Term:             basetypes.NewInt32Value(3),
		Type:             basetypes.NewStringValue("HOURLY"),
		EndsAt:           basetypes.NewStringValue("2023-12-14 17:09:47 +0000 UTC"),
		State:            basetypes.NewStringValue("ACTIVE"),
	}
	got := adaptContractToContractResource(sdkContract)

	assert.Equal(t, want, got)
}

func Test_adaptInstanceDetailsToInstanceResource(t *testing.T) {
	marketAppId := "marketAppId"
	reference := "reference"
	isoSdk := publiccloud.Iso{
		Id: "isoId",
	}

	instance := publiccloud.InstanceDetails{
		Id:                  "id",
		Type:                publiccloud.TYPENAME_C3_2XLARGE,
		Region:              "region",
		Reference:           *publiccloud.NewNullableString(&reference),
		MarketAppId:         *publiccloud.NewNullableString(&marketAppId),
		State:               publiccloud.STATE_CREATING,
		RootDiskSize:        50,
		RootDiskStorageType: publiccloud.STORAGETYPE_CENTRAL,
		Contract: publiccloud.InstanceContractDetails{
			Type: publiccloud.CONTRACTTYPE_MONTHLY,
		},
		Image: publiccloud.Image{
			Id: "UBUNTU_20_04_64BIT",
		},
		Ips: []publiccloud.IpDetails{
			{
				Ip: "127.0.0.1",
			},
		},
		Iso:               *publiccloud.NewNullableIso(&isoSdk),
		HasPrivateNetwork: true,
	}

	diags := diag.Diagnostics{}

	got := adaptInstanceDetailsToInstanceResource(
		instance,
		context.TODO(),
		&diags,
	)

	require.False(t, diags.HasError())

	assert.Equal(t, "id", got.ID.ValueString())
	assert.Equal(t, "region", got.Region.ValueString())
	assert.Equal(t, "CREATING", got.State.ValueString())
	assert.Equal(t, int32(50), got.RootDiskSize.ValueInt32())
	assert.Equal(t, "CENTRAL", got.RootDiskStorageType.ValueString())
	assert.Equal(t, "marketAppId", got.MarketAppID.ValueString())
	assert.Equal(t, "reference", got.Reference.ValueString())
	assert.Equal(t, "lsw.c3.2xlarge", got.Type.ValueString())

	image := imageResourceModel{}
	got.Image.As(context.TODO(), &image, basetypes.ObjectAsOptions{})
	assert.Equal(t, "UBUNTU_20_04_64BIT", image.ID.ValueString())

	contract := contractResourceModel{}
	got.Contract.As(context.TODO(), &contract, basetypes.ObjectAsOptions{})
	assert.Equal(t, "MONTHLY", contract.Type.ValueString())

	var ips []ipResourceModel
	got.IPs.ElementsAs(context.TODO(), &ips, false)
	assert.Len(t, ips, 1)
	assert.Equal(t, "127.0.0.1", ips[0].IP.ValueString())

	iso := isoResourceModel{}
	got.ISO.As(context.TODO(), &iso, basetypes.ObjectAsOptions{})
	assert.Equal(t, "isoId", iso.ID.ValueString())

	assert.True(t, got.HasPrivateNetwork.ValueBool())
}
