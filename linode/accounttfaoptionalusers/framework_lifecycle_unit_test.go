//go:build unit

package accounttfaoptionalusers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strconv"
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

func TestTerraformTFAOptionalUsersLifecycle(t *testing.T) {
	for _, clearBeforeDestroy := range []bool{false, true} {
		t.Run("clear_before_destroy_"+strconv.FormatBool(clearBeforeDestroy), func(t *testing.T) {
			var mu sync.Mutex
			users := []string{"unmanaged-user"}
			var puts int
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				mu.Lock()
				defer mu.Unlock()
				assert.Equal(t, "/v4/account/tfa-optional-users", req.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				switch req.Method {
				case http.MethodPut:
					var body linodego.TFAOptionalUsersUpdateOptions
					if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
						t.Errorf("decode PUT: %v", err)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					assert.NotNil(t, body.Usernames)
					users = slices.Clone(body.Usernames)
					slices.Reverse(users)
					puts++
					assert.NoError(t, json.NewEncoder(w).Encode(map[string]any{
						"tfa_optional_users_count": len(users), "url": "/account/tfa-optional-users",
					}))
				case http.MethodGet:
					page, _ := strconv.Atoi(req.URL.Query().Get("page"))
					page = max(1, page)
					data := []linodego.TFAOptionalUser{}
					if page <= len(users) {
						data = append(data, linodego.TFAOptionalUser{Username: users[page-1]})
					}
					assert.NoError(t, json.NewEncoder(w).Encode(map[string]any{
						"data": data, "page": page, "pages": max(1, len(users)), "results": len(users),
					}))
				default:
					t.Errorf("unexpected method %s", req.Method)
					w.WriteHeader(http.StatusMethodNotAllowed)
				}
			}))
			defer server.Close()

			config := func(usernames string) string {
				return fmt.Sprintf(`
provider "linode" {
  token       = "test-token"
  url         = %q
  api_version = "v4"
}
resource "linode_account_tfa_optional_users" "test" {
  usernames = [%s]
}
data "linode_account_tfa_optional_users" "test" {
  depends_on = [linode_account_tfa_optional_users.test]
}
`, server.URL, usernames)
			}
			checks := func(usernames ...string) []statecheck.StateCheck {
				values := make([]knownvalue.Check, len(usernames))
				for i, username := range usernames {
					values[i] = knownvalue.StringExact(username)
				}
				return []statecheck.StateCheck{
					statecheck.ExpectKnownValue("linode_account_tfa_optional_users.test",
						tfjsonpath.New("id"), knownvalue.StringExact("account")),
					statecheck.ExpectKnownValue("linode_account_tfa_optional_users.test",
						tfjsonpath.New("usernames"), knownvalue.SetExact(values)),
					statecheck.ExpectKnownValue("data.linode_account_tfa_optional_users.test",
						tfjsonpath.New("usernames"), knownvalue.SetExact(values)),
				}
			}
			factories := map[string]func() (tfprotov6.ProviderServer, error){
				"linode": providerserver.NewProtocol6WithError(linode.CreateFrameworkProvider("test")),
			}
			steps := []resource.TestStep{
				{
					Config:            config(`"alice", "bob"`),
					ConfigStateChecks: checks("alice", "bob"),
				},
				{
					ResourceName:      "linode_account_tfa_optional_users.test",
					ImportState:       true,
					ImportStateId:     "account",
					ImportStateVerify: true,
				},
				{
					Config:   config(`"bob", "alice"`),
					PlanOnly: true,
				},
				{
					PreConfig: func() {
						mu.Lock()
						defer mu.Unlock()
						users = []string{"external-user", "alice"}
					},
					Config:            config(`"alice", "bob"`),
					ConfigStateChecks: checks("alice", "bob"),
				},
				{
					Config:            config(""),
					ConfigStateChecks: checks(),
				},
				{
					Config:            config(`"alice"`),
					ConfigStateChecks: checks("alice"),
				},
			}
			if clearBeforeDestroy {
				steps = append(steps, resource.TestStep{
					Config:            config(""),
					ConfigStateChecks: checks(),
				})
			}
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: factories,
				Steps:                    steps,
			})

			mu.Lock()
			defer mu.Unlock()
			if clearBeforeDestroy {
				assert.Empty(t, users)
				assert.Equal(t, 5, puts)
			} else {
				assert.Equal(t, []string{"alice"}, users)
				assert.Equal(t, 4, puts)
			}
		})
	}
}
