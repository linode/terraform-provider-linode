//go:build integration || firewalls

package firewalls_test

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/firewalls/tmpl"
)

const testFirewallDataName = "data.linode_firewalls.test"

var testRegion string

func init() {
	region, err := acceptance.GetRandomRegionWithCaps([]string{"Linodes"})
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccDataSourceFirewalls_basic(t *testing.T) {
	t.Parallel()

	firewallName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataAll(t, firewallName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckResourceAttrGreaterThan(testFirewallDataName, "firewalls.#", 0),
					resource.TestCheckResourceAttrSet(testFirewallDataName, "firewalls.0.label"),
					resource.TestCheckResourceAttrSet(testFirewallDataName, "firewalls.0.tags.#"),
					resource.TestCheckResourceAttrSet(testFirewallDataName, "firewalls.0.created"),
					resource.TestCheckResourceAttrSet(testFirewallDataName, "firewalls.0.updated"),
					resource.TestCheckResourceAttrSet(testFirewallDataName, "firewalls.0.inbound_policy"),
					resource.TestCheckResourceAttrSet(testFirewallDataName, "firewalls.0.outbound_policy"),
				),
			},
			{
				Config: tmpl.DataFilter(t, firewallName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testFirewallDataName, "firewalls.0.label", firewallName),
					resource.TestCheckResourceAttr(testFirewallDataName, "firewalls.0.inbound_policy", "DROP"),
					resource.TestCheckResourceAttr(testFirewallDataName, "firewalls.0.outbound_policy", "ACCEPT"),
					resource.TestCheckResourceAttr(testFirewallDataName, "firewalls.0.disabled", "false"),
					resource.TestCheckResourceAttr(testFirewallDataName, "firewalls.0.status", "enabled"),
					resource.TestCheckResourceAttr(testFirewallDataName, "firewalls.0.tags.#", "2"),
					resource.TestCheckResourceAttr(testFirewallDataName, "firewalls.0.linodes.#", "1"),
					resource.TestCheckResourceAttrSet(testFirewallDataName, "firewalls.0.created"),
					resource.TestCheckResourceAttrSet(testFirewallDataName, "firewalls.0.updated"),

					resource.TestCheckResourceAttr(testFirewallDataName, "firewalls.0.inbound.0.label", "allow-http"),
					resource.TestCheckResourceAttr(testFirewallDataName, "firewalls.0.inbound.0.action", "ACCEPT"),
					resource.TestCheckResourceAttr(testFirewallDataName, "firewalls.0.inbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallDataName, "firewalls.0.inbound.0.ports", "80"),
					resource.TestCheckResourceAttr(testFirewallDataName, "firewalls.0.inbound.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(testFirewallDataName, "firewalls.0.inbound.0.ipv4.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallDataName, "firewalls.0.inbound.0.ipv6.#", "1"),
					resource.TestCheckResourceAttr(testFirewallDataName, "firewalls.0.inbound.0.ipv6.0", "::/0"),

					resource.TestCheckResourceAttr(testFirewallDataName, "firewalls.0.outbound.0.label", "reject-http"),
					resource.TestCheckResourceAttr(testFirewallDataName, "firewalls.0.outbound.0.action", "DROP"),
					resource.TestCheckResourceAttr(testFirewallDataName, "firewalls.0.outbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallDataName, "firewalls.0.outbound.0.ports", "80"),
					resource.TestCheckResourceAttr(testFirewallDataName, "firewalls.0.outbound.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(testFirewallDataName, "firewalls.0.outbound.0.ipv4.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallDataName, "firewalls.0.outbound.0.ipv6.#", "1"),
					resource.TestCheckResourceAttr(testFirewallDataName, "firewalls.0.outbound.0.ipv6.0", "::/0"),

					resource.TestCheckResourceAttr(testFirewallDataName, "firewalls.0.devices.#", "2"),
					resource.TestCheckResourceAttrSet(testFirewallDataName, "firewalls.0.devices.0.label"),
					resource.TestCheckResourceAttrSet(testFirewallDataName, "firewalls.0.devices.0.type"),
				),
			},
		},
	})
}
