//go:build integration || nb

package nb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/nb/tmpl"
)

func TestAccDataSourceNodeBalancer_basic(t *testing.T) {
	t.Parallel()

	resName := "data.linode_nodebalancer.foobar"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkNodeBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, nodebalancerName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkNodeBalancerExists,
					resource.TestCheckResourceAttr(resName, "label", nodebalancerName),
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
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkNodeBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataFirewalls(t, nodebalancerName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkNodeBalancerExists,
					resource.TestCheckResourceAttr(resName, "label", nodebalancerName),
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
				),
			},
		},
	})
}
