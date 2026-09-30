package objectstorage

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/leaseweb/leaseweb-go-sdk/objectstorage"
	"github.com/stretchr/testify/assert"
)

func Test_adaptGroupToGroupResource(t *testing.T) {
	t.Run("all fields are set", func(t *testing.T) {
		s3Policies := `{"Statement":[{"Effect":"Deny","Action":"s3:*"}]}`

		sdkGroup := objectstorage.Group{
			Id:          objectstorage.PtrString("groupId"),
			DisplayName: objectstorage.PtrString("Read only"),
			UniqueName:  objectstorage.PtrString("read-only"),
			S3Policies:  *objectstorage.NewNullableString(&s3Policies),
		}

		got := adaptGroupToGroupResource(sdkGroup, "objectStorageId")

		assert.Equal(t, "groupId", got.ID.ValueString())
		assert.Equal(t, "Read only", got.DisplayName.ValueString())
		assert.Equal(t, "read-only", got.UniqueName.ValueString())
		assert.Equal(t, "objectStorageId", got.ObjectStorageID.ValueString())
		assert.Equal(t, s3Policies, got.S3Policies.ValueString())
	})

	t.Run("a null policy stays null", func(t *testing.T) {
		sdkGroup := objectstorage.Group{
			Id:          objectstorage.PtrString("groupId"),
			DisplayName: objectstorage.PtrString("Read only"),
			UniqueName:  objectstorage.PtrString("read-only"),
		}

		got := adaptGroupToGroupResource(sdkGroup, "objectStorageId")

		assert.True(t, got.S3Policies.IsNull())
	})

	// s3_policies is compared semantically so that the formatting & property
	// order the API returns does not differ from what jsonencode produces.
	t.Run("policies differing only in formatting are equal", func(t *testing.T) {
		fromConfig := jsontypes.NewNormalizedValue(
			`{"Statement":[{"Action":"s3:*","Effect":"Deny"}]}`,
		)
		fromAPI := jsontypes.NewNormalizedValue(
			`{"Statement": [{"Effect": "Deny","Action": "s3:*"}]}`,
		)

		equal, diags := fromConfig.StringSemanticEquals(context.TODO(), fromAPI)

		assert.False(t, diags.HasError())
		assert.True(t, equal)
	})
}
