//go:build integration || nbtypes

package nbtypes_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/nbtypes/tmpl"
)

func TestAccDataSourceNodeBalancerTypes_basic(t *testing.T) {
	t.Parallel()

	dataSourceName := "data.linode_nb_types.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "types.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "types.0.id", "nodebalancer"),
					resource.TestCheckResourceAttr(dataSourceName, "types.0.label", "NodeBalancer"),
					resource.TestCheckResourceAttrSet(dataSourceName, "types.0.price.0.hourly"),
					resource.TestCheckResourceAttrSet(dataSourceName, "types.0.price.0.monthly"),
					resource.TestCheckResourceAttrSet(dataSourceName, "types.0.region_prices.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "types.0.region_prices.0.hourly"),
					resource.TestCheckResourceAttrSet(dataSourceName, "types.0.region_prices.0.monthly"),
				),
			},
		},
	})
}
