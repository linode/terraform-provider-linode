package linode

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
)

const testFirewallResName = "linode_firewall.test"

func init() {
	resource.AddTestSweepers("linode_firewall", &resource.Sweeper{
		Name: "linode_firewall",
		F:    testSweepLinodeFirewall,
	})
}

func testSweepLinodeFirewall(prefix string) error {
	client, err := getClientForSweepers()
	if err != nil {
		return fmt.Errorf("failed to get client: %s", err)
	}

	firewalls, err := client.ListFirewalls(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to get firewalls: %s", err)
	}
	for _, firewall := range firewalls {
		if !shouldSweepAcceptanceTestResource(prefix, firewall.Label) {
			continue
		}
		if err := client.DeleteFirewall(context.Background(), firewall.ID); err != nil {
			return fmt.Errorf("failed to destroy firewall %d during sweep: %s", firewall.ID, err)
		}
	}

	return nil
}

func testAccCheckLinodeFirewallDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(linodego.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_firewall" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to parse Firewall ID: %s", err)
		}

		if id == 0 {
			return fmt.Errorf("should not have Firewall ID of 0")
		}

		if _, err = client.GetFirewall(context.Background(), id); err == nil {
			return fmt.Errorf("should not find Firewall %d existing after delete", id)
		} else if apiErr, ok := err.(*linodego.Error); !ok {
			return fmt.Errorf("expected API Error but got %#v", err)
		} else if apiErr.Code != 404 {
			return fmt.Errorf("expected an error 404 but got %#v", apiErr)
		}
	}

	return nil
}

func TestAccLinodeFirewall_basic(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("tf_test")
	devicePrefix := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeFirewallBasic(name, devicePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testFirewallResName, "label", name),
					resource.TestCheckResourceAttr(testFirewallResName, "disabled", "false"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ports.1889509032", "80"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.addresses.1080289494", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ports.1889509032", "80"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.addresses.1080289494", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "devices.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "devices.0.type", "linode"),
					resource.TestCheckResourceAttr(testFirewallResName, "linodes.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "tags.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "tags.3632233996", "test"),
					resource.TestCheckResourceAttrSet(testFirewallResName, "devices.0.url"),
					resource.TestCheckResourceAttrSet(testFirewallResName, "devices.0.id"),
					resource.TestCheckResourceAttrSet(testFirewallResName, "devices.0.entity_id"),
					resource.TestCheckResourceAttrSet(testFirewallResName, "devices.0.label"),
				),
			},
			{
				ResourceName:      testFirewallResName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLinodeFirewall_updates(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("tf_test")
	newName := acctest.RandomWithPrefix("tf_test")
	devicePrefix := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeFirewallBasic(name, devicePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testFirewallResName, "label", name),
					resource.TestCheckResourceAttr(testFirewallResName, "disabled", "false"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ports.1889509032", "80"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.addresses.1080289494", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ports.1889509032", "80"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.addresses.1080289494", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "devices.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "devices.0.type", "linode"),
					resource.TestCheckResourceAttr(testFirewallResName, "linodes.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "tags.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "tags.3632233996", "test"),
				),
			},
			{
				Config: testAccCheckLinodeFirewallUpdates(newName, devicePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testFirewallResName, "label", newName),
					resource.TestCheckResourceAttr(testFirewallResName, "disabled", "true"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.#", "3"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ports.1889509032", "80"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.addresses.1080289494", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.1.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.1.ports.3638101695", "443"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.1.addresses.1080289494", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.2.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.2.ports.1685985038", "22"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.2.addresses.1080289494", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.#", "0"),
					resource.TestCheckResourceAttr(testFirewallResName, "linodes.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "tags.#", "2"),
					resource.TestCheckResourceAttr(testFirewallResName, "tags.3632233996", "test"),
					resource.TestCheckResourceAttr(testFirewallResName, "tags.331058520", "test2"),
				),
			},
		},
	})
}

func testAccCheckLinodeFirewallInstance(prefix, identifier string) string {
	return fmt.Sprintf(`
resource "linode_instance" "%[1]s" {
	label = "%.15[2]s-%[1]s"
	group = "tf_test"
	type = "g6-nanode-1"
	region = "us-east"
	disk {
		label = "disk"
		image = "linode/alpine3.11"
		root_pass = "b4d_p4s5"
		authorized_keys = ["%[3]s"]
		size = 3000
	}
}`, identifier, prefix, publicKeyMaterial)
}

func testAccCheckLinodeFirewallBasic(name, devicePrefix string) string {
	return testAccCheckLinodeFirewallInstance(devicePrefix, "one") + fmt.Sprintf(`
resource "linode_firewall" "test" {
	label = "%s"
	tags  = ["test"]

	inbound {
		protocol  = "TCP"
		ports     = ["80"]
		addresses = ["0.0.0.0/0"]
	}

	outbound {
		protocol  = "TCP"
		ports     = ["80"]
		addresses = ["0.0.0.0/0"]
	}

	linodes = [linode_instance.one.id]
}`, name)
}

func testAccCheckLinodeFirewallUpdates(name, devicePrefix string) string {
	return testAccCheckLinodeFirewallInstance(devicePrefix, "one") +
		testAccCheckLinodeFirewallInstance(devicePrefix, "two") +
		fmt.Sprintf(`
resource "linode_firewall" "test" {
	label    = "%s"
	tags     = ["test", "test2"]
    disabled = true

	inbound {
		protocol  = "TCP"
		ports     = ["80"]
		addresses = ["0.0.0.0/0"]
	}

	inbound {
		protocol  = "TCP"
		ports     = ["443"]
		addresses = ["0.0.0.0/0"]
	}

	inbound {
		protocol  = "TCP"
		ports     = ["22"]
		addresses = ["0.0.0.0/0"]
	}

	linodes = [
		linode_instance.two.id,
	]
}`, name)
}
