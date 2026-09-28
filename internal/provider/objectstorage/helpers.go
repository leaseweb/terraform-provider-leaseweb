package objectstorage

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// adaptNullableTimeToStringValue converts a nullable time to a Terraform
// StringValue in RFC3339 format. The object storage API accepts & returns
// RFC3339, so timestamps have to round trip in that same format, otherwise a
// configured value does not match the value read back.
func adaptNullableTimeToStringValue(value *time.Time) basetypes.StringValue {
	if value == nil {
		return basetypes.NewStringNull()
	}

	return basetypes.NewStringValue(value.Format(time.RFC3339))
}
