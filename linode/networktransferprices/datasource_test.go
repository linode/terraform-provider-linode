//go:build integration || networktransferprices

package networktransferprices_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/networktransferprices/tmpl"
)

func TestAccDataSourceNetworkTransferPrices_basic(t *testing.T) {
	t.Parallel()

	dataSourceName := "data.linode_network_transfer_prices.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "types.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "types.0.id", "distributed_network_transfer"),
					resource.TestCheckResourceAttr(dataSourceName, "types.0.label", "Distributed Network Transfer"),
					resource.TestCheckResourceAttrSet(dataSourceName, "types.0.transfer"),
					resource.TestCheckResourceAttrSet(dataSourceName, "types.0.price.0.hourly"),
					resource.TestCheckResourceAttrSet(dataSourceName, "types.0.price.0.monthly"),
					resource.TestCheckResourceAttrSet(dataSourceName, "types.1.region_prices.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "types.1.region_prices.0.hourly"),
					resource.TestCheckResourceAttrSet(dataSourceName, "types.1.region_prices.0.monthly"),
				),
			},
		},
	})
}
