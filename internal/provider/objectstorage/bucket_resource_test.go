package objectstorage

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/objectstorage"
	"github.com/stretchr/testify/assert"
)

func Test_adaptBucketToBucketResource(t *testing.T) {
	t.Run("all fields are set", func(t *testing.T) {
		creationTime, _ := time.Parse(time.RFC3339, "2024-01-15T08:30:00Z")

		sdkBucket := objectstorage.Bucket{
			Name:                objectstorage.PtrString("my-bucket"),
			Region:              objectstorage.PtrString("nl-01"),
			IsVersioningEnabled: objectstorage.PtrBool(true),
			Quota: *objectstorage.NewNullableBucketQuota(&objectstorage.BucketQuota{
				Value: objectstorage.PtrFloat32(1000),
				Unit:  objectstorage.PtrString("GB"),
			}),
			ObjectCount: *objectstorage.NewNullableInt32(objectstorage.PtrInt32(42)),
			Used: &objectstorage.BucketUsedSpace{
				Value: objectstorage.PtrFloat32(500),
				Unit:  objectstorage.PtrString("GB"),
			},
			CreationTime:   *objectstorage.NewNullableTime(&creationTime),
			IsBeingDeleted: objectstorage.PtrBool(false),
		}

		diags := diag.Diagnostics{}
		got := adaptBucketToBucketResource(
			sdkBucket,
			"objectStorageId",
			&diags,
			context.TODO(),
		)

		assert.False(t, diags.HasError())
		assert.Equal(t, "my-bucket", got.Name.ValueString())
		assert.Equal(t, "my-bucket", got.ID.ValueString())
		assert.Equal(t, "objectStorageId", got.ObjectStorageID.ValueString())
		assert.Equal(t, "nl-01", got.Region.ValueString())
		assert.True(t, got.IsVersioningEnabled.ValueBool())
		assert.False(t, got.IsBeingDeleted.ValueBool())
		assert.Equal(t, int32(42), got.ObjectCount.ValueInt32())

		// The API returns a quota object but accepts an integer amount of GB,
		// so it has to map back onto the integer the practitioner configured.
		assert.Equal(t, int32(1000), got.Quota.ValueInt32())

		// Timestamps have to round trip as RFC3339, the format the API accepts.
		assert.Equal(t, "2024-01-15T08:30:00Z", got.CreationTime.ValueString())

		used := bucketUsedSpaceResourceModel{}
		got.Used.As(context.TODO(), &used, basetypes.ObjectAsOptions{})
		assert.Equal(t, "GB", used.Unit.ValueString())
		assert.InDelta(t, float64(500), used.Value.ValueFloat64(), 0)
	})

	t.Run("nullable fields are null", func(t *testing.T) {
		sdkBucket := objectstorage.Bucket{
			Name:                objectstorage.PtrString("my-bucket"),
			Region:              objectstorage.PtrString("nl-01"),
			IsVersioningEnabled: objectstorage.PtrBool(false),
			IsBeingDeleted:      objectstorage.PtrBool(true),
		}

		diags := diag.Diagnostics{}
		got := adaptBucketToBucketResource(
			sdkBucket,
			"objectStorageId",
			&diags,
			context.TODO(),
		)

		assert.False(t, diags.HasError())
		assert.True(t, got.Quota.IsNull())
		assert.True(t, got.ObjectCount.IsNull())
		assert.True(t, got.CreationTime.IsNull())
		assert.True(t, got.Used.IsNull())
		assert.True(t, got.IsBeingDeleted.ValueBool())
	})
}

func Test_adaptNullableTimeToStringValue(t *testing.T) {
	t.Run("a time is formatted as RFC3339", func(t *testing.T) {
		value, _ := time.Parse(time.RFC3339, "2026-07-10T13:14:42Z")

		got := adaptNullableTimeToStringValue(&value)

		assert.False(t, got.IsNull())
		assert.Equal(t, "2026-07-10T13:14:42Z", got.ValueString())
	})

	t.Run("nil becomes a null string", func(t *testing.T) {
		got := adaptNullableTimeToStringValue(nil)

		assert.True(t, got.IsNull())
	})
}
