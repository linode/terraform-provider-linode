//go:build integration || firewalldevice

package firewalldevice_test

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	acceptanceTmpl "github.com/linode/terraform-provider-linode/v2/linode/acceptance/tmpl"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/firewalldevice/tmpl"
)

var testRegion string

func init() {
	region, err := acceptance.GetRandomRegionWithCaps([]string{"Cloud Firewall", "NodeBalancers"})
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccResourceFirewallDevice_basic_smoke(t *testing.T) {
	t.Parallel()

	var firewall linodego.Firewall

	firewallName := "linode_firewall.foobar"
	instanceName := "linode_instance.foobar"
	deviceName := "linode_firewall_device.foobar"

	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreventPostDestroyRefresh: true,
		PreCheck:                  func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories:  acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acceptanceTmpl.ProviderNoPoll(t) + tmpl.Basic(t, label, testRegion),
				Check: resource.ComposeAggregateTestCheckFunc(
					acceptance.CheckFirewallExists(firewallName, &firewall),
					resource.TestCheckResourceAttrSet(deviceName, "created"),
				),
			},
			// Refresh the state and verify the attachment
			{
				Config: acceptanceTmpl.ProviderNoPoll(t) + tmpl.Basic(t, label, testRegion),
				Check: resource.ComposeAggregateTestCheckFunc(
					acceptance.CheckFirewallExists(firewallName, &firewall),
					resource.TestCheckResourceAttr(firewallName, "devices.#", "1"),
					resource.TestCheckResourceAttrPair(firewallName, "linodes.0", instanceName, "id"),
				),
			},
			{
				ResourceName:      deviceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: resourceImportStateID,
			},
			{
				Config: acceptanceTmpl.ProviderNoPoll(t) + tmpl.Detached(t, label, testRegion),
				Check: resource.ComposeAggregateTestCheckFunc(
					acceptance.CheckFirewallExists(firewallName, &firewall),
				),
			},
			// Refresh the state and verify the detachment
			{
				Config: acceptanceTmpl.ProviderNoPoll(t) + tmpl.Detached(t, label, testRegion),
				Check: resource.ComposeAggregateTestCheckFunc(
					acceptance.CheckFirewallExists(firewallName, &firewall),
					resource.TestCheckResourceAttr(firewallName, "devices.#", "0"),
					resource.TestCheckResourceAttr(firewallName, "linodes.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceFirewallDevice_withNodeBalancer(t *testing.T) {
	t.Parallel()

	var firewall linodego.Firewall

	firewallName := "linode_firewall.foobar"
	nodebalancerName := "linode_nodebalancer.foobar"
	deviceName := "linode_firewall_device.foobar"

	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreventPostDestroyRefresh: true,
		PreCheck:                  func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories:  acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acceptanceTmpl.ProviderNoPoll(t) + tmpl.WithNodeBalancer(t, label, testRegion),
				Check: resource.ComposeAggregateTestCheckFunc(
					acceptance.CheckFirewallExists(firewallName, &firewall),
					resource.TestCheckResourceAttrSet(deviceName, "created"),
				),
			},
			// Refresh the state and verify the attachment
			{
				Config: acceptanceTmpl.ProviderNoPoll(t) + tmpl.WithNodeBalancer(t, label, testRegion),
				Check: resource.ComposeAggregateTestCheckFunc(
					acceptance.CheckFirewallExists(firewallName, &firewall),
					resource.TestCheckResourceAttr(firewallName, "devices.#", "1"),
					resource.TestCheckResourceAttrPair(firewallName, "nodebalancers.0", nodebalancerName, "id"),
				),
			},
			{
				ResourceName:      deviceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: resourceImportStateID,
			},
			{
				Config: acceptanceTmpl.ProviderNoPoll(t) + tmpl.Detached(t, label, testRegion),
				Check: resource.ComposeAggregateTestCheckFunc(
					acceptance.CheckFirewallExists(firewallName, &firewall),
				),
			},
			// Refresh the state and verify the detachment
			{
				Config: acceptanceTmpl.ProviderNoPoll(t) + tmpl.Detached(t, label, testRegion),
				Check: resource.ComposeAggregateTestCheckFunc(
					acceptance.CheckFirewallExists(firewallName, &firewall),
					resource.TestCheckResourceAttr(firewallName, "devices.#", "0"),
					resource.TestCheckResourceAttr(firewallName, "linodes.#", "0"),
				),
			},
		},
	})
}

func resourceImportStateID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_firewall_device" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return "", fmt.Errorf("Error parsing ID %v to int", rs.Primary.ID)
		}

		firewallID, err := strconv.Atoi(rs.Primary.Attributes["firewall_id"])
		if err != nil {
			return "", fmt.Errorf("Error parsing firewall_id %v to int", rs.Primary.Attributes["firewall_id"])
		}
		return fmt.Sprintf("%d,%d", firewallID, id), nil
	}

	return "", fmt.Errorf("Error finding firewall_device")
}
