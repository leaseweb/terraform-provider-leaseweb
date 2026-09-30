package objectstorage

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/leaseweb/leaseweb-go-sdk/objectstorage"
	"github.com/stretchr/testify/assert"
)

func Test_adaptUserToUserResource(t *testing.T) {
	t.Run("all fields are set", func(t *testing.T) {
		sdkUser := objectstorage.User{
			Id:         objectstorage.PtrString("userId"),
			FullName:   objectstorage.PtrString("Backup agent"),
			UniqueName: objectstorage.PtrString("backup-agent"),
			Groups:     []string{"groupId1", "groupId2"},
		}

		diags := diag.Diagnostics{}
		got := adaptUserToUserResource(
			sdkUser,
			"objectStorageId",
			&diags,
			context.TODO(),
		)

		assert.False(t, diags.HasError())
		assert.Equal(t, "userId", got.ID.ValueString())
		assert.Equal(t, "Backup agent", got.FullName.ValueString())
		assert.Equal(t, "backup-agent", got.UniqueName.ValueString())
		assert.Equal(t, "objectStorageId", got.ObjectStorageID.ValueString())

		var groups []string
		got.Groups.ElementsAs(context.TODO(), &groups, false)
		assert.ElementsMatch(t, []string{"groupId1", "groupId2"}, groups)
	})

	t.Run("a user without groups has an empty set", func(t *testing.T) {
		sdkUser := objectstorage.User{
			Id:         objectstorage.PtrString("userId"),
			FullName:   objectstorage.PtrString("Backup agent"),
			UniqueName: objectstorage.PtrString("backup-agent"),
		}

		diags := diag.Diagnostics{}
		got := adaptUserToUserResource(
			sdkUser,
			"objectStorageId",
			&diags,
			context.TODO(),
		)

		assert.False(t, diags.HasError())

		var groups []string
		got.Groups.ElementsAs(context.TODO(), &groups, false)
		assert.Empty(t, groups)
	})
}
