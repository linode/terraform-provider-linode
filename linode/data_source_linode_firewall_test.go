package linode

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

const testFirewallDataName = "data.linode_firewall.test"

func TestAccDataSourceLinodeFirewall_basic(t *testing.T) {
	t.Parallel()

	firewallName := acctest.RandomWithPrefix("tf_test")
	devicePrefix := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: accTestWithProvider(testDataSourceLinodeFirewallBasic(firewallName, devicePrefix), map[string]interface{}{
					providerKeySkipInstanceReadyPoll: true,
				}),
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
					resource.TestCheckResourceAttr(testFirewallDataName, "devices.#", "1"),
					resource.TestCheckResourceAttr(testFirewallDataName, "devices.0.type", "linode"),
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

func testDataSourceLinodeFirewallBasic(firewallName, devicePrefix string) string {
	return testAccCheckLinodeFirewallBasic(firewallName, devicePrefix) + `
data "linode_firewall" "test" {
	id = linode_firewall.test.id
}
`
}
