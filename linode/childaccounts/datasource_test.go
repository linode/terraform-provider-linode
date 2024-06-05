//go:build integration && parent_child

package childaccounts_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/childaccounts/tmpl"
)

func TestAccDataSourceChildAccounts_basic(t *testing.T) {
	t.Parallel()

	allName := "data.linode_child_accounts.all"
	filterName := "data.linode_child_accounts.filter"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(allName, "child_accounts.0.euuid"),
					resource.TestCheckResourceAttrSet(allName, "child_accounts.0.email"),
					acceptance.CheckResourceAttrGreaterThan(allName, "child_accounts.#", 0),

					resource.TestCheckResourceAttrPair(
						filterName, "child_accounts.0.euuid",
						allName, "child_accounts.0.euuid",
					),
					resource.TestCheckResourceAttrPair(
						filterName, "child_accounts.0.email",
						allName, "child_accounts.0.email",
					),
					resource.TestCheckResourceAttr(filterName, "child_accounts.#", "1"),
				),
			},
		},
	})
}
