package objectstorage

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/objectstorage"
	"github.com/stretchr/testify/assert"
)

func Test_adaptAccessKeyToAccessKeyResource(t *testing.T) {
	originalState := accessKeyResourceModel{
		AccessKey:       basetypes.NewStringValue("BHGWSRW4OQ7QKIV5YTP0"),
		SecretAccessKey: basetypes.NewStringValue("WsndsW+WjRz14LP6p33L"),
		ObjectStorageID: basetypes.NewStringValue("objectStorageId"),
		UserID:          basetypes.NewStringValue("userId"),
	}

	t.Run("all fields are set", func(t *testing.T) {
		expiresAt, _ := time.Parse(time.RFC3339, "2026-07-10T13:14:42Z")

		sdkAccessKey := objectstorage.AccessKey{
			Id:          objectstorage.PtrString("accessKeyId"),
			DisplayName: objectstorage.PtrString("Root"),
			AccountId:   objectstorage.PtrString("accountId"),
			ExpiresAt:   *objectstorage.NewNullableTime(&expiresAt),
		}

		got := adaptAccessKeyToAccessKeyResource(sdkAccessKey, originalState)

		assert.Equal(t, "accessKeyId", got.ID.ValueString())
		assert.Equal(t, "Root", got.DisplayName.ValueString())
		assert.Equal(t, "accountId", got.AccountID.ValueString())
		assert.Equal(t, "objectStorageId", got.ObjectStorageID.ValueString())
		assert.Equal(t, "userId", got.UserID.ValueString())

		// expires_at is configurable, so it has to round trip as RFC3339.
		assert.Equal(t, "2026-07-10T13:14:42Z", got.ExpiresAt.ValueString())
	})

	// The list endpoint never returns the secrets, so a read must not blank
	// the values captured when the key was created.
	t.Run("secrets are carried over from state", func(t *testing.T) {
		sdkAccessKey := objectstorage.AccessKey{
			Id:        objectstorage.PtrString("accessKeyId"),
			AccountId: objectstorage.PtrString("accountId"),
		}

		got := adaptAccessKeyToAccessKeyResource(sdkAccessKey, originalState)

		assert.Equal(t, "BHGWSRW4OQ7QKIV5YTP0", got.AccessKey.ValueString())
		assert.Equal(t, "WsndsW+WjRz14LP6p33L", got.SecretAccessKey.ValueString())
	})

	t.Run("a key without an expiry stays null", func(t *testing.T) {
		sdkAccessKey := objectstorage.AccessKey{
			Id:        objectstorage.PtrString("accessKeyId"),
			AccountId: objectstorage.PtrString("accountId"),
		}

		got := adaptAccessKeyToAccessKeyResource(sdkAccessKey, originalState)

		assert.True(t, got.ExpiresAt.IsNull())
	})
}
