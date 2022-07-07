package ipv6range_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/ipv6range/tmpl"
)

func TestAccDataSourceIPv6Range_basic(t *testing.T) {
	t.Parallel()

	instanceLabel := acctest.RandomWithPrefix("tf_test")
	resourceName := "data.linode_ipv6_range.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, instanceLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "range"),

					resource.TestCheckResourceAttr(resourceName, "linodes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "is_bgp", "false"),
					resource.TestCheckResourceAttr(resourceName, "prefix", "64"),
					resource.TestCheckResourceAttr(resourceName, "region", "us-southeast"),
				),
			},
		},
	})
}
