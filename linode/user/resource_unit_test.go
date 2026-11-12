//go:build unit

package user

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpandGrantsEntitiesWithRemovals(t *testing.T) {
	oldEntities := []any{
		map[string]any{"id": 101, "permissions": "read_write"},
		map[string]any{"id": 202, "permissions": "read_only"},
	}

	newEntities := []any{
		map[string]any{"id": 202, "permissions": "read_only"},
	}

	readOnly := linodego.AccessLevelReadOnly

	expected := []linodego.EntityUserGrant{
		{ID: 101, Permissions: nil},
		{ID: 202, Permissions: &readOnly},
	}

	actual := expandGrantsEntitiesWithRemovals(oldEntities, newEntities)

	assert.ElementsMatch(t, expected, actual)
}

func TestExpandGrantsEntityChanges_omittedSetPlansRevocation(t *testing.T) {
	testSchema := map[string]*schema.Schema{
		"domain_grant":  resourceSchema["domain_grant"],
		"global_grants": resourceSchema["global_grants"],
	}
	oldData := schema.TestResourceDataRaw(t, testSchema, map[string]any{
		"domain_grant": []any{
			map[string]any{"id": 101, "permissions": "read_write"},
			map[string]any{"id": 202, "permissions": "read_only"},
		},
		"global_grants": []any{
			map[string]any{"add_domains": true},
		},
	})
	oldData.SetId("test-user")
	oldState := oldData.State()

	schemaMap := schema.InternalMap(testSchema)
	diff, err := schemaMap.Diff(
		context.Background(), oldState, terraform.NewResourceConfigRaw(map[string]any{}),
		nil, nil, true,
	)
	require.NoError(t, err)
	require.NotNil(t, diff, "omitting domain_grant should plan a change to an empty grant set")

	data, err := schemaMap.Data(oldState, diff)
	require.NoError(t, err)

	assert.True(t, data.HasChange("domain_grant"))
	assert.Zero(t, data.Get("domain_grant").(*schema.Set).Len())
	assert.False(t, diff.Destroy)
	assert.False(t, diff.RequiresNew())
	assert.False(t, data.HasChange("global_grants"))
	assert.Equal(t, oldData.Get("global_grants"), data.Get("global_grants"))

	actual := expandGrantsEntityChanges(data, "domain_grant")
	expected := []linodego.EntityUserGrant{
		{ID: 101, Permissions: nil},
		{ID: 202, Permissions: nil},
	}
	assert.ElementsMatch(t, expected, actual)
}
