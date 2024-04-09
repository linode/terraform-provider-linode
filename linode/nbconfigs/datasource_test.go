//go:build integration || nbconfigs

package nbconfigs_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/nbconfigs/tmpl"
)

func TestAccDataSourceNodeBalancerConfigs_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_nodebalancer_configs.foo"

	nbLabel := acctest.RandomWithPrefix("tf_test")
	nbRegion, err := acceptance.GetRandomRegionWithCaps([]string{"NodeBalancers"})
	if err != nil {
		log.Fatal(err)
	}

	port := 80

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, nbLabel, nbRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "nodebalancer_configs.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "nodebalancer_configs.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "nodebalancer_configs.0.nodebalancer_id"),
					resource.TestCheckResourceAttrSet(resourceName, "nodebalancer_configs.0.protocol"),
					resource.TestCheckResourceAttrSet(resourceName, "nodebalancer_configs.0.proxy_protocol"),
					resource.TestCheckResourceAttrSet(resourceName, "nodebalancer_configs.0.port"),
					resource.TestCheckResourceAttrSet(resourceName, "nodebalancer_configs.0.check_interval"),
					resource.TestCheckResourceAttrSet(resourceName, "nodebalancer_configs.0.check_passive"),
					resource.TestCheckResourceAttrSet(resourceName, "nodebalancer_configs.0.cipher_suite"),
					resource.TestCheckNoResourceAttr(resourceName, "nodebalancer_configs.0.ssl_common"),
					resource.TestCheckNoResourceAttr(resourceName, "nodebalancer_configs.0.ssl_ciphersuite"),
					resource.TestCheckResourceAttr(resourceName, "nodebalancer_configs.0.node_status.0.up", "0"),
					resource.TestCheckResourceAttr(resourceName, "nodebalancer_configs.0.node_status.0.down", "0"),
					resource.TestCheckNoResourceAttr(resourceName, "nodebalancer_configs.0.ssl_cert"),
					resource.TestCheckNoResourceAttr(resourceName, "nodebalancer_configs.0.ssl_key"),
				),
			},
			{
				Config: tmpl.DataFilter(t, nbLabel, nbRegion, port),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "nodebalancer_configs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "nodebalancer_configs.0.port", fmt.Sprint(port)),
				),
			},
		},
	})
}
