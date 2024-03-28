//go:build integration && parent_child

package childaccount_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/childaccount/tmpl"
)

func TestAccDataSourceChildAccount_basic_smoke(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_child_account.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				Check: resource.ComposeTestCheckFunc(
					// Most of the account fields can be empty so we can't consistently
					// check them here
					resource.TestCheckResourceAttrSet(resourceName, "euuid"),
					resource.TestCheckResourceAttrSet(resourceName, "email"),
					resource.TestCheckResourceAttrSet(resourceName, "balance"),
				),
			},
		},
	})
}
