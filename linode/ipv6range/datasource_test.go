//go:build integration || ipv6range

package ipv6range_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/ipv6range/tmpl"
)

func TestAccDataSourceIPv6Range_basic(t *testing.T) {
	t.Parallel()

	instanceLabel := acctest.RandomWithPrefix("tf_test")
	resourceName := "data.linode_ipv6_range.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, instanceLabel, testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "range"),

					resource.TestCheckResourceAttr(resourceName, "linodes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "is_bgp", "false"),
					resource.TestCheckResourceAttr(resourceName, "prefix", "64"),
					resource.TestCheckResourceAttr(resourceName, "region", testRegion),
				),
			},
		},
	})
}
