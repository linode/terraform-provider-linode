//go:build integration || nbnode

package nbnode_test

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/nbnode/tmpl"
)

func TestAccDataSourceNodeBalancerNode_basic(t *testing.T) {
	t.Parallel()

	resName := "data.linode_nodebalancer_node.foonode"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")
	rootPass := acctest.RandString(64)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
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

func TestAccDataSourceNodeBalancerNode_vpc(t *testing.T) {
	t.Parallel()

	dsName := "data.linode_nodebalancer_node.test"
	label := acctest.RandomWithPrefix("tf-test")
	rootPass := acctest.RandString(64)

	targetRegion, err := acceptance.GetRandomRegionWithCaps([]string{"NodeBalancers", "VPCs"}, "core")
	if err != nil {
		log.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkNodeBalancerNodeDestroy,

		Steps: []resource.TestStep{
			{
				Config: tmpl.DataVPC(t, label, targetRegion, rootPass),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						dsName,
						tfjsonpath.New("nodebalancer_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dsName,
						tfjsonpath.New("config_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dsName,
						tfjsonpath.New("subnet_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dsName,
						tfjsonpath.New("address"),
						knownvalue.StringExact("10.0.0.5:80"),
					),
					statecheck.ExpectKnownValue(
						dsName,
						tfjsonpath.New("label"),
						knownvalue.StringExact(label),
					),
					statecheck.ExpectKnownValue(
						dsName,
						tfjsonpath.New("weight"),
						knownvalue.Int64Exact(50),
					),
				},
			},
		},
	})
}
