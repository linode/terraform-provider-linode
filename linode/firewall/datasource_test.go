//go:build integration || firewall

package firewall_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/firewall/tmpl"
)

const testFirewallDataName = "data.linode_firewall.test"

func TestAccDataSourceFirewall_basic(t *testing.T) {
	t.Parallel()

	firewallName := acctest.RandomWithPrefix("tf_test")
	devicePrefix := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, firewallName, devicePrefix, testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testFirewallDataName, "label", firewallName),
					resource.TestCheckResourceAttr(testFirewallDataName, "disabled", "false"),
					resource.TestCheckResourceAttr(testFirewallDataName, "inbound_policy", "DROP"),
					resource.TestCheckResourceAttr(testFirewallDataName, "inbound.#", "1"),
					resource.TestCheckResourceAttr(testFirewallDataName, "inbound.0.action", "ACCEPT"),
					resource.TestCheckResourceAttr(testFirewallDataName, "inbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallDataName, "inbound.0.ports", "80"),
					resource.TestCheckResourceAttr(testFirewallDataName, "inbound.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(testFirewallDataName, "inbound.0.ipv4.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallDataName, "inbound.0.ipv6.#", "1"),
					resource.TestCheckResourceAttr(testFirewallDataName, "inbound.0.ipv6.0", "::/0"),
					resource.TestCheckResourceAttr(testFirewallDataName, "outbound_policy", "DROP"),
					resource.TestCheckResourceAttr(testFirewallDataName, "outbound.#", "1"),
					resource.TestCheckResourceAttr(testFirewallDataName, "outbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallDataName, "outbound.0.ports", "80"),
					resource.TestCheckResourceAttr(testFirewallDataName, "outbound.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(testFirewallDataName, "outbound.0.ipv4.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallDataName, "outbound.0.ipv6.#", "1"),
					resource.TestCheckResourceAttr(testFirewallDataName, "outbound.0.ipv6.0", "2001:db8::/32"),
					resource.TestCheckResourceAttr(testFirewallDataName, "devices.#", "2"),
					resource.TestCheckResourceAttrSet(testFirewallDataName, "devices.0.type"),
					resource.TestCheckResourceAttr(testFirewallDataName, "nodebalancers.#", "1"),
					resource.TestCheckResourceAttr(testFirewallDataName, "linodes.#", "1"),
					resource.TestCheckResourceAttr(testFirewallDataName, "tags.#", "1"),
					resource.TestCheckResourceAttr(testFirewallDataName, "tags.0", "test"),
					resource.TestCheckResourceAttrSet(testFirewallDataName, "devices.0.url"),
					resource.TestCheckResourceAttrSet(testFirewallDataName, "devices.0.id"),
					resource.TestCheckResourceAttrSet(testFirewallDataName, "devices.0.entity_id"),
					resource.TestCheckResourceAttrSet(testFirewallDataName, "devices.0.label"),
				),
			},
		},
	})
}
