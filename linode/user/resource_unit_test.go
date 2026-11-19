//go:build unit

package user

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserResourceTFAEnforced(t *testing.T) {
	for _, enforced := range []bool{false, true} {
		t.Run(strconv.FormatBool(enforced), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/v4/account/users/test-user", req.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				assert.NoError(t, json.NewEncoder(w).Encode(linodego.User{
					Username: "test-user", Email: "test@example.com",
					TFAEnforced: enforced, TFAEnabled: !enforced,
				}))
			}))
			defer server.Close()
			client, err := linodego.NewClient(server.Client())
			require.NoError(t, err)
			client.SetBaseURL(server.URL)
			client.SetAPIVersion("v4")
			data := schema.TestResourceDataRaw(t, resourceSchema, nil)
			data.SetId("test-user")
			diags := readResource(context.Background(), data, &helper.ProviderMeta{Client: client})
			require.False(t, diags.HasError(), "%v", diags)
			assert.Equal(t, enforced, data.Get("tfa_enforced"))
			assert.Equal(t, !enforced, data.Get("tfa_enabled"))
			assert.True(t, resourceSchema["tfa_enforced"].Computed)
			assert.False(t, resourceSchema["tfa_enforced"].Optional)
		})
	}
}
