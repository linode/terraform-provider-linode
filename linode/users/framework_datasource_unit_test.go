//go:build unit

package users

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/frameworkfilter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersTFAEnforcedFilter(t *testing.T) {
	for _, enforced := range []bool{false, true} {
		t.Run(strconv.FormatBool(enforced), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/v4/account/users", req.URL.Path)
				assert.JSONEq(t, `{"+and":[{"+or":[{"tfa_enforced":`+strconv.FormatBool(enforced)+`}]}],"+order":"desc"}`,
					req.Header.Get("X-Filter"))
				w.Header().Set("Content-Type", "application/json")
				assert.NoError(t, json.NewEncoder(w).Encode(map[string]any{
					"data": []linodego.User{{
						Username: "test-user", TFAEnforced: enforced, TFAEnabled: !enforced,
					}},
					"page": 1, "pages": 1, "results": 1,
				}))
			}))
			defer server.Close()
			client, err := linodego.NewClient(server.Client())
			require.NoError(t, err)
			client.SetBaseURL(server.URL)
			client.SetAPIVersion("v4")
			ctx := context.Background()
			results, diagnostic := filterConfig.GetAndFilter(ctx, &client,
				frameworkfilter.FiltersModelType{{
					Name:    types.StringValue("tfa_enforced"),
					Values:  []types.String{types.StringValue(strconv.FormatBool(enforced))},
					MatchBy: types.StringValue("exact"),
				}}, listUsers, types.StringNull(), types.StringNull())
			require.Nil(t, diagnostic)
			var model UserFilterModel
			diags := model.parseUsers(ctx, &client, helper.AnySliceToTyped[linodego.User](results))
			require.False(t, diags.HasError(), "%v", diags)
			require.Len(t, model.Users, 1)
			assert.Equal(t, types.BoolValue(enforced), model.Users[0].TFAEnforced)
			assert.Equal(t, types.BoolValue(!enforced), model.Users[0].TFAEnabled)
		})
	}
}
