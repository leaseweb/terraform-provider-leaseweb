package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain/entity"
)

func Test_newDDos(t *testing.T) {
	ddos := entity.NewDdos("detectionProfile", "protectionType")
	got := newDdos(ddos)

	assert.Equal(
		t,
		"detectionProfile",
		got.DetectionProfile.ValueString(),
		"detectionProfile should be set",
	)
	assert.Equal(
		t,
		"protectionType",
		got.ProtectionType.ValueString(),
		"protectionType should be set",
	)
}
