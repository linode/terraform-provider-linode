//go:build unit

package accountsettings_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode"
	"github.com/stretchr/testify/assert"
)

func TestTerraformTFAEnforcementPolicy(t *testing.T) {
	var mu sync.Mutex
	enforced := true
	var puts int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		var result any
		switch req.URL.Path {
		case "/v4/account":
			result = map[string]any{"email": "test@example.com"}
		case "/v4/account/settings":
			if req.Method == http.MethodPut {
				var body linodego.AccountSettingsUpdateOptions
				assert.NoError(t, json.NewDecoder(req.Body).Decode(&body))
				if puts == 0 {
					assert.Nil(t, body.TFAEnforced, "omitting tfa_enforced must not reset the policy")
				}
				if body.TFAEnforced != nil {
					enforced = *body.TFAEnforced
				}
				puts++
			}
			result = linodego.AccountSettings{TFAEnforced: enforced, NetworkHelper: true}
		case "/v4/profile":
			result = linodego.Profile{Username: "test-user", TFAEnforced: enforced}
		case "/v4/account/users/test-user":
			result = linodego.User{Username: "test-user", TFAEnforced: enforced}
		case "/v4/account/users":
			result = map[string]any{
				"data": []linodego.User{{Username: "test-user", TFAEnforced: enforced}},
				"page": 1, "pages": 1, "results": 1,
			}
		default:
			t.Errorf("unexpected request %s %s", req.Method, req.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		assert.NoError(t, json.NewEncoder(w).Encode(result))
	}))
	defer server.Close()

	config := func(policy string) string {
		return fmt.Sprintf(`
provider "linode" {
  token       = "test-token"
  url         = %q
  api_version = "v4"
}
resource "linode_account_settings" "test" {
  network_helper = true
  %s
}
data "linode_account_settings" "test" {
  depends_on = [linode_account_settings.test]
}
data "linode_profile" "test" {
  depends_on = [linode_account_settings.test]
}
data "linode_user" "test" {
  username   = "test-user"
  depends_on = [linode_account_settings.test]
}
data "linode_users" "test" {
  filter {
    name   = "tfa_enforced"
    values = [tostring(linode_account_settings.test.tfa_enforced)]
  }
}
`, server.URL, policy)
	}
	checks := func(expected bool) []statecheck.StateCheck {
		result := make([]statecheck.StateCheck, 0, 5)
		for _, name := range []string{
			"linode_account_settings.test",
			"data.linode_account_settings.test",
			"data.linode_profile.test",
			"data.linode_user.test",
		} {
			result = append(result, statecheck.ExpectKnownValue(name,
				tfjsonpath.New("tfa_enforced"), knownvalue.Bool(expected)))
		}
		return append(result, statecheck.ExpectKnownValue("data.linode_users.test",
			tfjsonpath.New("users").AtSliceIndex(0).AtMapKey("tfa_enforced"), knownvalue.Bool(expected)))
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"linode": providerserver.NewProtocol6WithError(linode.CreateFrameworkProvider("test")),
		},
		Steps: []resource.TestStep{
			{Config: config(""), ConfigStateChecks: checks(true)},
			{Config: config("tfa_enforced = false"), ConfigStateChecks: checks(false)},
			{Config: config("tfa_enforced = true"), ConfigStateChecks: checks(true)},
		},
	})
	mu.Lock()
	defer mu.Unlock()
	assert.True(t, enforced, "destroying account settings must not change enforcement")
	assert.Equal(t, 3, puts)
}
