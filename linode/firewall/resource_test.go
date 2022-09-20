package firewall_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/firewall/tmpl"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

const testFirewallResName = "linode_firewall.test"

var (
	testProviders map[string]*schema.Provider
	testProvider  *schema.Provider
	testRegion    string
)

func init() {
	resource.AddTestSweepers("linode_firewall", &resource.Sweeper{
		Name: "linode_firewall",
		F:    sweep,
	})

	provider, providerMap := acceptance.CreateTestProvider()
	acceptance.ModifyProviderMeta(provider, func(ctx context.Context, config *helper.ProviderMeta) error {
		config.Config.SkipInstanceReadyPoll = true
		config.Config.SkipInstanceDeletePoll = true
		return nil
	})

	testProvider = provider
	testProviders = providerMap

	region, err := acceptance.GetRandomRegionWithCaps([]string{"Cloud Firewall"})
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
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
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, name, devicePrefix, testRegion),
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
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Minimum(t, name),
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
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.MultipleRules(t, name, devicePrefix, testRegion),
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
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: testProviders,
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
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, name, devicePrefix, testRegion),
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
				Config: tmpl.Updates(t, newName, devicePrefix, testRegion),
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
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, name, devicePrefix, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckFirewallExists(testProvider, testFirewallResName, &firewall),
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
					client := testProvider.Meta().(*helper.ProviderMeta).Client

					if err := client.DeleteFirewall(context.Background(), firewall.ID); err != nil {
						t.Fatalf("failed to delete firewall: %s", err)
					}
				},
				Config: tmpl.Basic(t, name, devicePrefix, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckFirewallExists(testProvider, testFirewallResName, &firewall),
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

func TestAccLinodeFirewall_emptyIPv6(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.NoIPv6(t, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testFirewallResName, "label", name),
					resource.TestCheckResourceAttr(testFirewallResName, "disabled", "false"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.#", "1"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testFirewallResName, "inbound.0.ports", "80"),
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
