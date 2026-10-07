package objectstorage

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/objectstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// stubBucketListAPI serves a canned sequence of bucket lists so the deletion
// poll can be driven without real API calls or real waiting. Embedding the
// interface leaves every other method unimplemented, so an unexpected call
// panics rather than silently succeeding.
type stubBucketListAPI struct {
	objectstorage.ObjectstorageAPI
	responses [][]objectstorage.Bucket
	err       error
	calls     int
}

func (s *stubBucketListAPI) GetObjectStorageBucketList(
	_ context.Context,
	_ string,
) objectstorage.ApiGetObjectStorageBucketListRequest {
	return objectstorage.ApiGetObjectStorageBucketListRequest{ApiService: s}
}

func (s *stubBucketListAPI) GetObjectStorageBucketListExecute(
	_ objectstorage.ApiGetObjectStorageBucketListRequest,
) (*objectstorage.GetObjectStorageBucketListResult, *http.Response, error) {
	s.calls++

	if s.err != nil {
		return nil, &http.Response{StatusCode: 500, Body: http.NoBody}, s.err
	}

	// The last response repeats, so a test can poll indefinitely.
	i := min(s.calls-1, len(s.responses)-1)

	result := objectstorage.NewGetObjectStorageBucketListResult()
	result.SetBuckets(s.responses[i])

	return result, &http.Response{StatusCode: 200}, nil
}

func newTestBucketResource(api objectstorage.ObjectstorageAPI) bucketResource {
	bucket := bucketResource{
		deletionPollInterval: time.Millisecond,
		deletionMaxRetries:   3,
	}
	bucket.ObjectstorageAPI = api

	return bucket
}

func beingDeletedBucket(name string, isBeingDeleted bool) []objectstorage.Bucket {
	return []objectstorage.Bucket{
		{
			Name:           objectstorage.PtrString(name),
			IsBeingDeleted: objectstorage.PtrBool(isBeingDeleted),
		},
	}
}

func Test_bucketResource_waitUntilBucketDeleted(t *testing.T) {
	t.Run("returns as soon as the bucket is gone", func(t *testing.T) {
		api := &stubBucketListAPI{responses: [][]objectstorage.Bucket{{}}}
		diags := diag.Diagnostics{}

		err := newTestBucketResource(api).
			waitUntilBucketDeleted(context.TODO(), "objectStorageId", "my-bucket", &diags)

		require.NoError(t, err)
		assert.False(t, diags.HasError())
		assert.Equal(t, 1, api.calls)
	})

	// The API flags a bucket while its removal is in flight, so the poll is
	// equally done once that flag clears.
	t.Run("returns once the bucket stops reporting itself as being deleted", func(t *testing.T) {
		api := &stubBucketListAPI{
			responses: [][]objectstorage.Bucket{beingDeletedBucket("my-bucket", false)},
		}
		diags := diag.Diagnostics{}

		err := newTestBucketResource(api).
			waitUntilBucketDeleted(context.TODO(), "objectStorageId", "my-bucket", &diags)

		require.NoError(t, err)
		assert.Equal(t, 1, api.calls)
	})

	t.Run("keeps polling while the bucket is still being deleted", func(t *testing.T) {
		api := &stubBucketListAPI{
			responses: [][]objectstorage.Bucket{
				beingDeletedBucket("my-bucket", true),
				beingDeletedBucket("my-bucket", true),
				{},
			},
		}
		diags := diag.Diagnostics{}

		err := newTestBucketResource(api).
			waitUntilBucketDeleted(context.TODO(), "objectStorageId", "my-bucket", &diags)

		require.NoError(t, err)
		assert.False(t, diags.HasError())
		assert.Equal(t, 3, api.calls)
	})

	// Another bucket being deleted at the same time must not end this poll.
	t.Run("ignores other buckets in the list", func(t *testing.T) {
		api := &stubBucketListAPI{
			responses: [][]objectstorage.Bucket{beingDeletedBucket("another-bucket", true)},
		}
		diags := diag.Diagnostics{}

		err := newTestBucketResource(api).
			waitUntilBucketDeleted(context.TODO(), "objectStorageId", "my-bucket", &diags)

		require.NoError(t, err)
		assert.Equal(t, 1, api.calls)
	})

	t.Run("gives up once the retry limit is reached", func(t *testing.T) {
		api := &stubBucketListAPI{
			responses: [][]objectstorage.Bucket{beingDeletedBucket("my-bucket", true)},
		}
		diags := diag.Diagnostics{}

		err := newTestBucketResource(api).
			waitUntilBucketDeleted(context.TODO(), "objectStorageId", "my-bucket", &diags)

		require.ErrorContains(t, err, `timed out waiting for bucket "my-bucket" to be deleted`)
		assert.Equal(t, 3, api.calls)
	})

	t.Run("reports an error from the API", func(t *testing.T) {
		api := &stubBucketListAPI{err: errors.New("boom")}
		diags := diag.Diagnostics{}

		err := newTestBucketResource(api).
			waitUntilBucketDeleted(context.TODO(), "objectStorageId", "my-bucket", &diags)

		require.ErrorContains(t, err, "could not determine deletion status")
		assert.True(t, diags.HasError())
		assert.Equal(t, 1, api.calls)
	})

	// A cancelled apply must not keep polling for the remaining retries.
	t.Run("stops when the context is cancelled", func(t *testing.T) {
		api := &stubBucketListAPI{
			responses: [][]objectstorage.Bucket{beingDeletedBucket("my-bucket", true)},
		}
		diags := diag.Diagnostics{}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := newTestBucketResource(api).
			waitUntilBucketDeleted(ctx, "objectStorageId", "my-bucket", &diags)

		require.ErrorIs(t, err, context.Canceled)
		assert.Equal(t, 1, api.calls)
	})
}
