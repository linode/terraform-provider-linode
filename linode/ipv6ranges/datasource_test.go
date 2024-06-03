//go:build integration || ipv6ranges

package ipv6ranges_test

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/ipv6ranges/tmpl"
)

var testRegion string

func init() {
	region, err := acceptance.GetRandomRegionWithCaps([]string{"Linodes"})
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccDataSourceIPv6Ranges_basic(t *testing.T) {
	t.Parallel()

	instanceLabel := acctest.RandomWithPrefix("tf_test")
	dataSourceName := "data.linode_ipv6_ranges.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, instanceLabel, testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "ranges.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ranges.0.range"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ranges.0.route_target"),
					resource.TestCheckResourceAttr(dataSourceName, "ranges.0.region", testRegion),
					resource.TestCheckResourceAttr(dataSourceName, "ranges.0.prefix", "64"),
				),
			},
		},
	})
}
