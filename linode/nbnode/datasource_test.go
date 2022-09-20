package nbnode_test

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/nbnode/tmpl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

func TestAccDataSourceNodeBalancerNode_basic(t *testing.T) {
	t.Parallel()

	resName := "data.linode_nodebalancer_node.foonode"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    testProviders,
		CheckDestroy: checkNodeBalancerNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, nodebalancerName, testRegion),
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
