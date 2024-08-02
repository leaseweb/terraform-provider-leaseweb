package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewImage(t *testing.T) {
	state := "state"
	stateReason := "stateReason"
	region := "region"
	createdAt := time.Now()
	updatedAt := time.Now()
	custom := false

	got := NewImage(
		"UBUNTU_24_04_64BIT",
		"name",
		"version",
		"family",
		"flavour",
		"architecture",
		&state,
		&stateReason,
		&region,
		&createdAt,
		&updatedAt,
		&custom,
		[]string{"marketApp"},
		[]string{"storageType"},
	)
	want := Image{
		Id:           "UBUNTU_24_04_64BIT",
		Name:         "name",
		Version:      "version",
		Family:       "family",
		Flavour:      "flavour",
		Architecture: "architecture",
		State:        &state,
		StateReason:  &stateReason,
		Region:       &region,
		CreatedAt:    &createdAt,
		UpdatedAt:    &updatedAt,
		Custom:       &custom,
		MarketApps:   []string{"marketApp"},
		StorageTypes: []string{"storageType"},
	}

	assert.Equal(t, want, got)
}
