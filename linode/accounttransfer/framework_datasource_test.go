//go:build integration || accounttransfer || act_tests

package accounttransfer_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/accounttransfer/tmpl"
)

func TestAccDataSourceAccountTransfer_basic(t *testing.T) {
	acceptance.OptInTest(t)

	t.Parallel()

	resourceName := "data.linode_account_transfer.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.StringExact("account_transfer")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("billable"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("quota"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("used"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("region_transfers"), knownvalue.NotNull()),
				},
			},
		},
	})
}
