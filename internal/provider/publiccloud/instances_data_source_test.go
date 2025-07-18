package publiccloud

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publiccloud"
	"github.com/stretchr/testify/assert"
)

func Test_adaptContractToContractDataSource(t *testing.T) {
	endsAt, _ := time.Parse(
		"2006-01-02 15:04:05",
		"2023-12-14 17:09:47",
	)
	sdkContract := publiccloud.InstanceContractDetails{
		BillingFrequency: publiccloud.BILLINGFREQUENCY__1,
		EndsAt:           *publiccloud.NewNullableTime(&endsAt),
		Term:             publiccloud.CONTRACTTERM__3,
		Type:             publiccloud.CONTRACTTYPE_HOURLY,
		State:            publiccloud.CONTRACTSTATE_ACTIVE,
	}

	want := contractDataSourceModel{
		BillingFrequency: basetypes.NewInt32Value(1),
		EndsAt:           basetypes.NewStringValue("2023-12-14 17:09:47 +0000 UTC"),
		Term:             basetypes.NewInt32Value(3),
		Type:             basetypes.NewStringValue("HOURLY"),
		State:            basetypes.NewStringValue("ACTIVE"),
	}
	got := adaptContractToContractDataSource(sdkContract)

	assert.Equal(t, want, got)
}
