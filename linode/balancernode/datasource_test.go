package balancernode_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

func TestAccDataSourceNodeBalancerNode_basic(t *testing.T) {
	t.Parallel()

	resName := "data.linode_nodebalancer_node.foonode"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkNodeBalancerNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.AccTestWithProvider(dataSourceConfigBasic(nodebalancerName), map[string]interface{}{
					acceptance.SkipInstanceReadyPollKey: true,
				}),
				Check: resource.ComposeTestCheckFunc(
					checkNodeBalancerNodeExists,
					resource.TestCheckResourceAttr(resName, "label", nodebalancerName),
					resource.TestCheckResourceAttrSet(resName, "status"),
					resource.TestCheckResourceAttr(resName, "mode", "accept"),
					resource.TestCheckResourceAttr(resName, "weight", "50"),
				),
			},
		},
	})
}

func dataSourceConfigBasic(nodeBalancerName string) string {
	return resourceConfigBasic(nodeBalancerName) + `
data "linode_nodebalancer_node" "foonode" {
	id = "${linode_nodebalancer_node.foonode.id}"
	config_id = "${linode_nodebalancer_config.foofig.id}"
	nodebalancer_id = "${linode_nodebalancer.foobar.id}"
}
`
}
