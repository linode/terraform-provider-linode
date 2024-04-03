//go:build integration || nbnode

package nbnode_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/nbnode/tmpl"
)

func TestAccDataSourceNodeBalancerNode_basic(t *testing.T) {
	t.Parallel()

	resName := "data.linode_nodebalancer_node.foonode"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkNodeBalancerNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, nodebalancerName, testRegion, rootPass),
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
