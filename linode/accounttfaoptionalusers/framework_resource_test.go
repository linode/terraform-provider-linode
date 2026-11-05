//go:build integration || accounttfaoptionalusers

package accounttfaoptionalusers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/accounttfaoptionalusers/tmpl"
)

const (
	testResName = "linode_account_tfa_optional_users.test"
	testEmail   = "dl-linode-dev-dx@akamai.com"
)

// The tests replace the account-wide exemption list, so skip unless it
// starts empty to avoid wiping existing exemptions.
func preCheckNoOptionalUsers(t *testing.T) {
	acceptance.PreCheck(t)

	client, err := acceptance.GetTestClient()
	if err != nil {
		t.Fatalf("failed to get client: %s", err)
	}

	users, err := client.ListTFAOptionalUsers(context.Background(), nil)
	if err != nil {
		t.Fatalf("failed to list TFA optional users: %s", err)
	}

	if len(users) > 0 {
		t.Skipf("skipping: the account's TFA optional users list must be empty, got %d user(s)", len(users))
	}
}

// The exemption list is account-wide, so this test is opt-in to avoid
// conflicting with concurrent runs.
func TestAccResourceAccountTFAOptionalUsers_basic(t *testing.T) {
	acceptance.OptInTest(t)

	username := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { preCheckNoOptionalUsers(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, username, testEmail, true),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testResName, tfjsonpath.New("id"), knownvalue.StringExact("account")),
					statecheck.ExpectKnownValue(testResName, tfjsonpath.New("usernames"),
						knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact(username)})),
				},
			},
			{
				ResourceName:      testResName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Destroy doesn't modify the list, so clear it explicitly.
				Config: tmpl.Basic(t, username, testEmail, false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testResName, tfjsonpath.New("usernames"), knownvalue.SetSizeExact(0)),
				},
			},
		},
	})
}
