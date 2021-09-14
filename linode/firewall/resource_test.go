package firewall_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/firewall/tmpl"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

const testFirewallResName = "linode_firewall.test"

func init() {
	resource.AddTestSweepers("linode_firewall", &resource.Sweeper{
		Name: "linode_firewall",
		F:    sweep,
	})
}

func sweep(prefix string) error {
	client, err := acceptance.GetClientForSweepers()
	if err != nil {
		return fmt.Errorf("failed to get client: %s", err)
	}

	firewalls, err := client.ListFirewalls(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to get firewalls: %s", err)
	}
	for _, firewall := range firewalls {
		if !acceptance.ShouldSweep(prefix, firewall.Label) {
			continue
		}
		if err := client.DeleteFirewall(context.Background(), firewall.ID); err != nil {
			return fmt.Errorf("failed to destroy firewall %d during sweep: %s", firewall.ID, err)
		}
	}

	return nil
}

func TestAccLinodeFirewall_basic(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("tf_test")
	devicePrefix := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CheckLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.AccTestWithProvider(tmpl.Basic(t, name, devicePrefix),
					map[string]interface{}{
						acceptance.SkipInstanceReadyPollKey: true,
					}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testFirewallResName, "label", name),
					resource.TestCheckResourceAttr(testFirewallResName, "disabled", "false"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound_policy", "DROP"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.action", "ACCEPT"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ports", "80"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv4.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv6.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv6.0", "::/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound_policy", "DROP"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ports", "80"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ipv4.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ipv6.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ipv6.0", "2001:db8::/32"),
					resource.TestCheckResourceAttr(testFirewallResName, "devices.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "devices.0.type", "linode"),
					resource.TestCheckResourceAttr(testFirewallResName, "linodes.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "tags.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "tags.0", "test"),
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

func TestAccLinodeFirewall_minimum(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CheckLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.AccTestWithProvider(tmpl.Minimum(t, name),
					map[string]interface{}{
						acceptance.SkipInstanceReadyPollKey: true,
					}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testFirewallResName, "label", name),
					resource.TestCheckResourceAttr(testFirewallResName, "disabled", "false"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ports", ""),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv4.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv6.#", "0"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.#", "0"),
					resource.TestCheckResourceAttr(testFirewallResName, "devices.#", "0"),
					resource.TestCheckResourceAttr(testFirewallResName, "linodes.#", "0"),
					resource.TestCheckResourceAttr(testFirewallResName, "tags.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "tags.0", "test"),
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

func TestAccLinodeFirewall_multipleRules(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("tf_test")
	devicePrefix := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CheckLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.AccTestWithProvider(tmpl.MultipleRules(t, name, devicePrefix),
					map[string]interface{}{
						acceptance.SkipInstanceReadyPollKey: true,
					}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testFirewallResName, "label", name),
					resource.TestCheckResourceAttr(testFirewallResName, "disabled", "false"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound_policy", "DROP"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.#", "2"),

					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.action", "ACCEPT"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ports", "80"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv4.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv6.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv6.0", "::/0"),

					resource.TestCheckResourceAttr(testFirewallResName, "inbound.1.action", "ACCEPT"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.1.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.1.ports", "443"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.1.ipv4.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.1.ipv4.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.1.ipv6.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.1.ipv6.0", "::/0"),

					resource.TestCheckResourceAttr(testFirewallResName, "outbound_policy", "DROP"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.#", "2"),

					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ports", "80"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ipv4.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ipv6.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ipv6.0", "2001:db8::/32"),

					resource.TestCheckResourceAttr(testFirewallResName, "outbound.1.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.1.ports", "443"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.1.ipv4.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.1.ipv4.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.1.ipv6.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.1.ipv6.0", "2001:db8::/32"),

					resource.TestCheckResourceAttr(testFirewallResName, "devices.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "devices.0.type", "linode"),
					resource.TestCheckResourceAttr(testFirewallResName, "linodes.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "tags.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "tags.0", "test"),
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

func TestAccLinodeFirewall_no_device(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CheckLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.NoDevice(t, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testFirewallResName, "label", name),
					resource.TestCheckResourceAttr(testFirewallResName, "disabled", "false"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ports", "80"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv6.0", "::/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ports", "80"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ipv6.0", "::/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "devices.#", "0"),
					resource.TestCheckResourceAttr(testFirewallResName, "linodes.#", "0"),
					resource.TestCheckResourceAttr(testFirewallResName, "tags.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "tags.0", "test"),
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
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CheckLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.AccTestWithProvider(tmpl.Basic(t, name, devicePrefix),
					map[string]interface{}{
						acceptance.SkipInstanceReadyPollKey: true,
					}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testFirewallResName, "label", name),
					resource.TestCheckResourceAttr(testFirewallResName, "disabled", "false"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound_policy", "DROP"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.action", "ACCEPT"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ports", "80"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv4.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv6.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv6.0", "::/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound_policy", "DROP"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.action", "ACCEPT"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ports", "80"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ipv4.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ipv6.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ipv6.0", "2001:db8::/32"),
					resource.TestCheckResourceAttr(testFirewallResName, "devices.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "devices.0.type", "linode"),
					resource.TestCheckResourceAttr(testFirewallResName, "linodes.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "tags.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "tags.0", "test"),
				),
			},
			{
				Config: acceptance.AccTestWithProvider(tmpl.Updates(t, newName, devicePrefix),
					map[string]interface{}{
						acceptance.SkipInstanceReadyPollKey: true,
					}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testFirewallResName, "label", newName),
					resource.TestCheckResourceAttr(testFirewallResName, "disabled", "true"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound_policy", "ACCEPT"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.#", "3"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.action", "DROP"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ports", "80"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv4.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv6.#", "2"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv6.0", "::/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv6.1", "ff00::/8"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.1.action", "DROP"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.1.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.1.ports", "443"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.1.ipv4.#", "2"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.1.ipv4.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.1.ipv4.1", "127.0.0.1/32"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.1.ipv6.#", "0"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.2.action", "DROP"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.2.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.2.ports", "22"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.2.ipv4.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.2.ipv4.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.2.ipv6.#", "0"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound_policy", "ACCEPT"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.#", "0"),
					resource.TestCheckResourceAttr(testFirewallResName, "linodes.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "tags.#", "2"),
					resource.TestCheckResourceAttr(testFirewallResName, "tags.0", "test"),
					resource.TestCheckResourceAttr(testFirewallResName, "tags.1", "test2"),
				),
			},
		},
	})
}

func TestAccLinodeFirewall_externalDelete(t *testing.T) {
	t.Parallel()

	var firewall linodego.Firewall
	name := acctest.RandomWithPrefix("tf_test")
	devicePrefix := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CheckLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.AccTestWithProvider(tmpl.Basic(t, name, devicePrefix),
					map[string]interface{}{
						acceptance.SkipInstanceReadyPollKey: true,
					}),
				Check: resource.ComposeTestCheckFunc(
					checkFirewallExists(testFirewallResName, &firewall),
					resource.TestCheckResourceAttr(testFirewallResName, "label", name),
					resource.TestCheckResourceAttr(testFirewallResName, "disabled", "false"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound_policy", "DROP"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.action", "ACCEPT"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ports", "80"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv4.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv6.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv6.0", "::/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound_policy", "DROP"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.action", "ACCEPT"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ports", "80"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ipv4.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ipv6.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ipv6.0", "2001:db8::/32"),
					resource.TestCheckResourceAttr(testFirewallResName, "devices.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "devices.0.type", "linode"),
					resource.TestCheckResourceAttr(testFirewallResName, "linodes.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "tags.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "tags.0", "test"),
				),
			},
			{
				PreConfig: func() {
					// Delete the Firewall external from Terraform
					client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

					if err := client.DeleteFirewall(context.Background(), firewall.ID); err != nil {
						t.Fatalf("failed to delete firewall: %s", err)
					}
				},
				Config: acceptance.AccTestWithProvider(tmpl.Basic(t, name, devicePrefix),
					map[string]interface{}{
						acceptance.SkipInstanceReadyPollKey: true,
					}),
				Check: resource.ComposeTestCheckFunc(
					checkFirewallExists(testFirewallResName, &firewall),
					resource.TestCheckResourceAttr(testFirewallResName, "label", name),
					resource.TestCheckResourceAttr(testFirewallResName, "disabled", "false"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound_policy", "DROP"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.action", "ACCEPT"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ports", "80"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv4.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv6.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ipv6.0", "::/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound_policy", "DROP"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.action", "ACCEPT"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ports", "80"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ipv4.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ipv6.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "outbound.0.ipv6.0", "2001:db8::/32"),
					resource.TestCheckResourceAttr(testFirewallResName, "devices.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "devices.0.type", "linode"),
					resource.TestCheckResourceAttr(testFirewallResName, "linodes.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "tags.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "tags.0", "test"),
				),
			},
		},
	})
}

func checkFirewallExists(name string, firewall *linodego.Firewall) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		found, err := client.GetFirewall(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of Firewall %s: %s", rs.Primary.Attributes["label"], err)
		}

		*firewall = *found

		return nil
	}
}
