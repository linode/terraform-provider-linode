package firewalldevice_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/firewalldevice/tmpl"
)

func TestAccResourceFirewallDevice_basic(t *testing.T) {
	t.Parallel()

	var firewall linodego.Firewall

	firewallName := "linode_firewall.foobar"
	instanceName := "linode_instance.foobar"
	deviceName := "linode_firewall_device.foobar"

	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreventPostDestroyRefresh: true,
		PreCheck:                  func() { acceptance.PreCheck(t) },
		Providers:                 acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label),
				Check: resource.ComposeAggregateTestCheckFunc(
					acceptance.CheckFirewallExists(firewallName, &firewall),
					resource.TestCheckResourceAttrSet(deviceName, "created"),
				),
			},
			// Refresh the state and verify the attachment
			{
				Config: tmpl.Basic(t, label),
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
				Config: tmpl.Detached(t, label),
				Check: resource.ComposeAggregateTestCheckFunc(
					acceptance.CheckFirewallExists(firewallName, &firewall),
				),
			},
			// Refresh the state and verify the detachment
			{
				Config: tmpl.Detached(t, label),
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
