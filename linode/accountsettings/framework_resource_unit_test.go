//go:build unit

package accountsettings

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountSettingsTFAEnforcedUpdate(t *testing.T) {
	for _, tt := range []struct {
		name        string
		value       types.Bool
		wantPresent bool
		want        bool
	}{
		{"enabled", types.BoolValue(true), true, true},
		{"disabled", types.BoolValue(false), true, false},
		{"unknown", types.BoolUnknown(), false, true},
		{"null", types.BoolNull(), false, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var puts int
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if req.URL.Path == "/v4/account" {
					_, err := w.Write([]byte(`{"email":"test@example.com"}`))
					assert.NoError(t, err)
					return
				}
				assert.Equal(t, http.MethodPut, req.Method)
				assert.Equal(t, "/v4/account/settings", req.URL.Path)
				puts++
				var body map[string]json.RawMessage
				assert.NoError(t, json.NewDecoder(req.Body).Decode(&body))
				value, present := body["tfa_enforced"]
				assert.Equal(t, tt.wantPresent, present)
				if present {
					assert.Equal(t, tt.value.String(), string(value))
				}
				assert.NoError(t, json.NewEncoder(w).Encode(map[string]any{"tfa_enforced": tt.want}))
			}))
			defer server.Close()
			client, err := linodego.NewClient(server.Client())
			require.NoError(t, err)
			client.SetBaseURL(server.URL)
			client.SetAPIVersion("v4")
			r := NewResource().(*Resource)
			r.Meta = &helper.FrameworkProviderMeta{Client: &client}
			plan := AccountSettingsModel{TFAEnforced: tt.value}
			var diags diag.Diagnostics
			r.createOrUpdateAccountSettings(context.Background(), &plan, &diags)
			require.False(t, diags.HasError(), "%v", diags)
			assert.Equal(t, types.BoolValue(tt.want), plan.TFAEnforced)
			assert.Equal(t, 1, puts)
		})
	}
}

func TestAccountSettingsTFAEnforcedState(t *testing.T) {
	model := AccountSettingsModel{TFAEnforced: types.BoolValue(false)}
	model.FlattenAccountSettings("", &linodego.AccountSettings{TFAEnforced: true}, true)
	assert.Equal(t, types.BoolValue(false), model.TFAEnforced)
	model.FlattenAccountSettings("", &linodego.AccountSettings{TFAEnforced: true}, false)
	assert.Equal(t, types.BoolValue(true), model.TFAEnforced)
	model.CopyFrom(AccountSettingsModel{TFAEnforced: types.BoolValue(false)}, true)
	assert.Equal(t, types.BoolValue(true), model.TFAEnforced)
	model.CopyFrom(AccountSettingsModel{TFAEnforced: types.BoolValue(false)}, false)
	assert.Equal(t, types.BoolValue(false), model.TFAEnforced)
}
