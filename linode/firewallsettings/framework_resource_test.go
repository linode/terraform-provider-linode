//go:build integration || firewallsettings

package firewallsettings_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/firewallsettings/tmpl"
)

var (
	originalFirewallSettings *linodego.FirewallSettings
	testFirewallID           int
)

const (
	resourceName     = "test"
	resourceFullName = "linode_firewall_settings." + resourceName
)

func init() {
	client, err := acceptance.GetTestClient()
	if err != nil {
		log.Fatalf("failed to get client: %s", err)
	}

	originalFirewallSettings, err = client.GetFirewallSettings(context.Background())
	if err != nil {
		log.Fatalf("failed to get firewall settings: %s", err)
	}

	testFirewall, err := client.CreateFirewall(context.Background(), linodego.FirewallCreateOptions{
		Label: acctest.RandomWithPrefix("tf_test"),
		Rules: linodego.FirewallRuleSet{
			InboundPolicy:  "DROP",
			OutboundPolicy: "ACCEPT",
		},
	})
	if err != nil {
		log.Fatalf("failed to setup test firewall: %s", err)
	}

	testFirewallID = testFirewall.ID
	time.Sleep(2 * time.Second) // Wait for the firewall to be fully created
}

func TestAccResourceFirewallSettings_basic(t *testing.T) {
	t.Parallel()
	t.Cleanup(func() {
		client, err := acceptance.GetTestClient()
		if err != nil {
			log.Fatalf("failed to get client: %s", err)
		}

		if originalFirewallSettings != nil {
			_, err = client.UpdateFirewallSettings(context.Background(), linodego.FirewallSettingsUpdateOptions{
				DefaultFirewallIDs: &linodego.DefaultFirewallIDsOptions{
					Linode:          linodego.Pointer(originalFirewallSettings.DefaultFirewallIDs.Linode),
					NodeBalancer:    linodego.Pointer(originalFirewallSettings.DefaultFirewallIDs.NodeBalancer),
					PublicInterface: linodego.Pointer(originalFirewallSettings.DefaultFirewallIDs.PublicInterface),
					VPCInterface:    linodego.Pointer(originalFirewallSettings.DefaultFirewallIDs.VPCInterface),
				},
			})
			if err != nil {
				log.Fatalf("failed to restore original firewall settings: %s", err)
			}
		}

		err = client.DeleteFirewall(context.Background(), testFirewallID)
		if err != nil {
			log.Fatalf("failed to delete test firewall: %s", err)
		}
	})
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, resourceName, testFirewallID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceFullName, tfjsonpath.New("default_firewall_ids"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(
						resourceFullName, tfjsonpath.New("default_firewall_ids").AtMapKey("linode"), knownvalue.Int64Exact(int64(testFirewallID)),
					),
					statecheck.ExpectKnownValue(
						resourceFullName, tfjsonpath.New("default_firewall_ids").AtMapKey("nodebalancer"), knownvalue.Int64Exact(int64(testFirewallID)),
					),
					statecheck.ExpectKnownValue(
						resourceFullName, tfjsonpath.New("default_firewall_ids").AtMapKey("public_interface"), knownvalue.Int64Exact(int64(testFirewallID)),
					),
					statecheck.ExpectKnownValue(
						resourceFullName, tfjsonpath.New("default_firewall_ids").AtMapKey("vpc_interface"), knownvalue.Int64Exact(int64(testFirewallID)),
					),
				},
			},
		},
	})
}
