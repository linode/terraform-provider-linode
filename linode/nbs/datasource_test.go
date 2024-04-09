//go:build integration || nbs

package nbs_test

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/nbs/tmpl"
)

func TestAccDataSourceNodeBalancers_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_nodebalancers.nbs"

	nbLabel := acctest.RandomWithPrefix("tf_test")
	nbRegion, err := acceptance.GetRandomRegionWithCaps([]string{"NodeBalancers"})
	if err != nil {
		log.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, nbLabel, nbRegion),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet(resourceName, "nodebalancers.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "nodebalancers.0.client_conn_throttle"),
					resource.TestCheckResourceAttrSet(resourceName, "nodebalancers.0.created"),
					resource.TestCheckResourceAttrSet(resourceName, "nodebalancers.0.hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "nodebalancers.0.ipv4"),
					resource.TestCheckResourceAttrSet(resourceName, "nodebalancers.0.ipv6"),
					resource.TestCheckResourceAttr(resourceName, "nodebalancers.0.label", nbLabel+"-0"),
					resource.TestCheckResourceAttr(resourceName, "nodebalancers.0.region", nbRegion),
					resource.TestCheckResourceAttr(resourceName, "nodebalancers.0.transfer.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "nodebalancers.0.transfer.0.in"),
					resource.TestCheckResourceAttrSet(resourceName, "nodebalancers.0.transfer.0.out"),
					resource.TestCheckResourceAttrSet(resourceName, "nodebalancers.0.transfer.0.total"),
					resource.TestCheckResourceAttr(resourceName, "nodebalancers.0.tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "nodebalancers.0.tags.0", "tf_test_1"),
					resource.TestCheckResourceAttrSet(resourceName, "nodebalancers.0.updated"),
				),
			},
			{
				Config: tmpl.DataFilter(t, nbLabel, nbRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "nodebalancers.#", "1"),
				),
			},
			{
				Config: tmpl.DataFilterEmpty(t, nbLabel, nbRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "nodebalancers.#", "0"),
				),
			},
			{
				Config: tmpl.DataOrder(t, nbLabel, nbRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "nodebalancers.0.label", nbLabel+"-0"),
				),
			},
		},
	})
}
