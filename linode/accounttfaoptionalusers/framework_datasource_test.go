//go:build integration || accounttfaoptionalusers

package accounttfaoptionalusers_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/accounttfaoptionalusers/tmpl"
)

const testDataSourceName = "data.linode_account_tfa_optional_users.test"

// The exemption list is account-wide, so this test is opt-in to avoid
// conflicting with concurrent runs.
func TestAccDataSourceAccountTFAOptionalUsers_basic(t *testing.T) {
	acceptance.OptInTest(t)

	username := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { preCheckNoOptionalUsers(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, username, testEmail, true),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testDataSourceName, tfjsonpath.New("id"), knownvalue.StringExact("account")),
					statecheck.ExpectKnownValue(testDataSourceName, tfjsonpath.New("usernames"),
						knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact(username)})),
				},
			},
			{
				// Destroy doesn't modify the list, so clear it explicitly.
				Config: tmpl.DataBasic(t, username, testEmail, false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testDataSourceName, tfjsonpath.New("usernames"), knownvalue.SetSizeExact(0)),
				},
			},
		},
	})
}
