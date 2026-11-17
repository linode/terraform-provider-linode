//go:build unit

package user

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	ctyjson "github.com/hashicorp/go-cty/cty/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
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

func userTestRawConfig(t *testing.T, config map[string]any) cty.Value {
	t.Helper()
	raw, err := json.Marshal(config)
	require.NoError(t, err)
	value, err := ctyjson.Unmarshal(raw, Resource().CoreConfigSchema().ImpliedType())
	require.NoError(t, err)
	return value
}

func TestUnrestrictUserWithOmittedGrant_Update(t *testing.T) {
	oldData := schema.TestResourceDataRaw(t, resourceSchema, map[string]any{
		"username": "test-user", "email": "test@example.com", "restricted": true,
		"domain_grant":  []any{map[string]any{"id": 101, "permissions": "read_write"}},
		"global_grants": []any{map[string]any{"add_domains": true}},
	})
	oldData.SetId("test-user")
	state := oldData.State()
	config := map[string]any{"username": "test-user", "email": "test@example.com", "restricted": false}
	state.RawConfig = userTestRawConfig(t, config)
	diff, err := Resource().Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), nil)
	require.NoError(t, err)
	require.NotNil(t, diff)
	data, err := schema.InternalMap(resourceSchema).Data(state, diff)
	require.NoError(t, err)

	var requests []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.Method+" "+r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path != "/v4/account/users/test-user" {
			http.Error(w, "unexpected request", http.StatusBadRequest)
			return
		}
		if r.Method == http.MethodPut {
			var body struct {
				Restricted bool `json:"restricted"`
			}
			assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			assert.False(t, body.Restricted)
		} else if r.Method != http.MethodGet {
			http.Error(w, "unexpected method", http.StatusBadRequest)
			return
		}
		_, _ = w.Write([]byte(`{"username":"test-user","email":"test@example.com","restricted":false}`))
	}))
	defer server.Close()
	client, err := linodego.NewClient(server.Client())
	require.NoError(t, err)
	client.SetBaseURL(server.URL)
	client.SetRetryCount(0)

	diags := updateResource(context.Background(), data, &helper.ProviderMeta{Client: client})
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, []string{"PUT /v4/account/users/test-user", "GET /v4/account/users/test-user"}, requests)
	assert.False(t, data.Get("restricted").(bool))
	assert.Zero(t, data.Get("domain_grant").(*schema.Set).Len())
	assert.Empty(t, data.Get("global_grants").([]any))

	newState := data.State()
	newState.RawConfig = userTestRawConfig(t, config)
	nextDiff, err := Resource().Diff(context.Background(), newState, terraform.NewResourceConfigRaw(config), nil)
	require.NoError(t, err)
	assert.Nil(t, nextDiff)
}

func TestValidateGrantsForUnrestrictedUsers(t *testing.T) {
	for _, tc := range []struct {
		name        string
		restricted  bool
		grant       string
		value       []any
		expectError bool
	}{
		{
			name:        "grant on unrestricted user is rejected",
			restricted:  false,
			grant:       "domain_grant",
			value:       []any{map[string]any{"id": 101, "permissions": "read_only"}},
			expectError: true,
		},
		{
			name:       "grant on restricted user is allowed",
			restricted: true,
			grant:      "domain_grant",
			value:      []any{map[string]any{"id": 101, "permissions": "read_only"}},
		},
		{
			name:       "unrestricted user without grants is allowed",
			restricted: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := map[string]any{"username": "test-user", "email": "test@example.com", "restricted": tc.restricted}
			if tc.grant != "" {
				config[tc.grant] = tc.value
			}

			state := &terraform.InstanceState{RawConfig: userTestRawConfig(t, config)}
			_, err := Resource().Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), nil)
			if tc.expectError {
				require.ErrorContains(t, err, "restricted")
			} else {
				require.NoError(t, err)
			}
		})
	}
}
