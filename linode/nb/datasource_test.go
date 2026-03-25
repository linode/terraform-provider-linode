//go:build integration || nb

package nb_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/nb/tmpl"
)

func TestAccDataSourceNodeBalancer_basic(t *testing.T) {
	t.Parallel()

	resName := "data.linode_nodebalancer.foobar"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")
	nbType := "common"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkNodeBalancerDestroy,

		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, nodebalancerName, testRegion, nbType),
				Check: resource.ComposeTestCheckFunc(
					checkNodeBalancerExists,
					resource.TestCheckResourceAttr(resName, "label", nodebalancerName),
					resource.TestCheckResourceAttr(resName, "type", "common"),
					resource.TestCheckResourceAttr(resName, "client_conn_throttle", "20"),
					resource.TestCheckResourceAttr(resName, "client_udp_sess_throttle", "10"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					resource.TestCheckResourceAttrSet(resName, "hostname"),
					resource.TestCheckResourceAttrSet(resName, "ipv4"),
					resource.TestCheckResourceAttrSet(resName, "ipv6"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttrSet(resName, "updated"),
					resource.TestCheckResourceAttr(resName, "transfer.#", "1"),
					resource.TestCheckResourceAttrSet(resName, "transfer.0.in"),
					resource.TestCheckResourceAttrSet(resName, "transfer.0.out"),
					resource.TestCheckResourceAttrSet(resName, "transfer.0.total"),
					resource.TestCheckResourceAttr(resName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resName, "tags.0", "tf_test"),
					resource.TestCheckResourceAttr(resName, "frontend_address_type", "public"),
					resource.TestCheckNoResourceAttr(resName, "frontend_vpc_subnet_id"),
				),
			},
		},
	})
}

func TestAccDataSourceNodeBalancer_firewalls(t *testing.T) {
	t.Parallel()

	resName := "data.linode_nodebalancer.foobar"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkNodeBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataFirewalls(t, nodebalancerName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkNodeBalancerExists,
					resource.TestCheckResourceAttr(resName, "label", nodebalancerName),
					resource.TestCheckResourceAttr(resName, "type", "common"),
					resource.TestCheckResourceAttr(resName, "client_conn_throttle", "20"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					resource.TestCheckResourceAttrSet(resName, "hostname"),
					resource.TestCheckResourceAttrSet(resName, "ipv4"),
					resource.TestCheckResourceAttrSet(resName, "ipv6"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttrSet(resName, "updated"),
					resource.TestCheckResourceAttr(resName, "transfer.#", "1"),
					resource.TestCheckResourceAttrSet(resName, "transfer.0.in"),
					resource.TestCheckResourceAttrSet(resName, "transfer.0.out"),
					resource.TestCheckResourceAttrSet(resName, "transfer.0.total"),
					resource.TestCheckResourceAttr(resName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resName, "tags.0", "tf_test"),
					acceptance.CheckResourceAttrGreaterThan(resName, "firewalls.#", 0),
					resource.TestCheckResourceAttr(resName, "firewalls.0.label", fmt.Sprintf("%v-fw", nodebalancerName)),
					resource.TestCheckResourceAttr(resName, "firewalls.0.inbound_policy", "DROP"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.inbound.#", "1"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.inbound.0.action", "ACCEPT"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.inbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.inbound.0.ports", "80"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.inbound.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.inbound.0.ipv4.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.inbound.0.ipv6.#", "1"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.inbound.0.ipv6.0", "::/0"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.outbound_policy", "DROP"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.outbound.#", "1"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.outbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.outbound.0.ports", "80"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.outbound.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.outbound.0.ipv4.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.outbound.0.ipv6.#", "1"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.outbound.0.ipv6.0", "2001:db8::/32"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.tags.#", "1"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.tags.0", "test"),
					resource.TestCheckResourceAttr(resName, "frontend_address_type", "public"),
					resource.TestCheckNoResourceAttr(resName, "frontend_vpc_subnet_id"),
				),
			},
		},
	})
}

func TestAccDataSourceNodeBalancer_vpc(t *testing.T) {
	t.Parallel()

	dsName := "data.linode_nodebalancer.test"
	nodebalancerName := acctest.RandomWithPrefix("tf-test")

	// Use random premium region, as not all regions support VPCs.
	targetRegion, err := acceptance.GetRandomPremiumRegion()
	if err != nil {
		log.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkNodeBalancerDestroy,

		Steps: []resource.TestStep{
			{
				Config: tmpl.DataVPC(t, nodebalancerName, targetRegion),
				Check:  checkNodeBalancerExists,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						dsName,
						tfjsonpath.New("type"),
						knownvalue.StringExact("common"),
					),
					statecheck.ExpectKnownValue(
						dsName,
						tfjsonpath.New("vpcs").AtSliceIndex(0).AtMapKey("subnet_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dsName,
						tfjsonpath.New("vpcs").AtSliceIndex(0).AtMapKey("ipv4_range"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dsName,
						tfjsonpath.New("frontend_address_type"),
						knownvalue.StringExact("public"),
					),
					statecheck.ExpectKnownValue(
						dsName,
						tfjsonpath.New("frontend_vpc_subnet_id"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

func TestAccDataSourceNodeBalancer_frontendVPC(t *testing.T) {
	t.Parallel()

	dsName := "data.linode_nodebalancer.test"
	nodebalancerName := acctest.RandomWithPrefix("tf-test")

	// Use random premium region, as not all regions support VPCs.
	targetRegion, err := acceptance.GetRandomPremiumRegion()
	if err != nil {
		log.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkNodeBalancerDestroy,

		Steps: []resource.TestStep{
			{
				Config: tmpl.DataFrontendVPC(t, nodebalancerName, targetRegion),
				Check:  checkNodeBalancerExists,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						dsName,
						tfjsonpath.New("type"),
						knownvalue.StringExact("premium"),
					),
					statecheck.ExpectKnownValue(
						dsName,
						tfjsonpath.New("vpcs").AtSliceIndex(0).AtMapKey("subnet_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dsName,
						tfjsonpath.New("vpcs").AtSliceIndex(0).AtMapKey("ipv4_range"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dsName,
						tfjsonpath.New("frontend_address_type"),
						knownvalue.StringExact("vpc"),
					),
					statecheck.ExpectKnownValue(
						dsName,
						tfjsonpath.New("frontend_vpc_subnet_id"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}
