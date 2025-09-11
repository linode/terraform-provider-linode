//go:build integration || firewalls

package firewalls_test

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/firewalls/tmpl"
)

const testFirewallDataName = "data.linode_firewalls.test"

var testRegion string

func init() {
	region, err := acceptance.GetRandomRegionWithCaps([]string{"Linodes"}, "core")
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccDataSourceFirewalls_basic(t *testing.T) {
	t.Parallel()

	firewallName := acctest.RandomWithPrefix("tf_test")
	acceptance.RunTestWithRetries(t, 3, func(t *acceptance.WrappedT) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: tmpl.DataAll(t, firewallName, testRegion),
					Check: resource.ComposeTestCheckFunc(
						acceptance.CheckResourceAttrGreaterThan(testFirewallDataName, "firewalls.#", 0),
					),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(testFirewallDataName, tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("label"), knownvalue.NotNull()),
						statecheck.ExpectKnownValue(testFirewallDataName, tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("tags"), knownvalue.NotNull()),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("created"),
							knownvalue.NotNull(),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("updated"),
							knownvalue.NotNull(),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("inbound_policy"),
							knownvalue.NotNull(),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("outbound_policy"),
							knownvalue.NotNull(),
						),
					},
				},
				{
					Config: tmpl.DataFilter(t, firewallName, testRegion),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("label"),
							knownvalue.StringExact(firewallName),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("inbound_policy"),
							knownvalue.StringExact("DROP"),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("outbound_policy"),
							knownvalue.StringExact("ACCEPT"),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("disabled"),
							knownvalue.Bool(false),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("status"),
							knownvalue.StringExact("enabled"),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("tags"),
							knownvalue.SetSizeExact(2),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("created"),
							knownvalue.NotNull(),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("updated"),
							knownvalue.NotNull(),
						),

						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("inbound").AtSliceIndex(0).AtMapKey("label"),
							knownvalue.StringExact("allow-http"),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("inbound").AtSliceIndex(0).AtMapKey("action"),
							knownvalue.StringExact("ACCEPT"),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("inbound").AtSliceIndex(0).AtMapKey("protocol"),
							knownvalue.StringExact("TCP"),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("inbound").AtSliceIndex(0).AtMapKey("ports"),
							knownvalue.StringExact("80"),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("inbound").AtSliceIndex(0).AtMapKey("ipv4"),
							knownvalue.SetSizeExact(1),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("inbound").AtSliceIndex(0).AtMapKey("ipv4").AtSliceIndex(0),
							knownvalue.StringExact("0.0.0.0/0"),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("inbound").AtSliceIndex(0).AtMapKey("ipv6"),
							knownvalue.SetSizeExact(1),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("inbound").AtSliceIndex(0).AtMapKey("ipv6").AtSliceIndex(0),
							knownvalue.StringExact("::/0"),
						),

						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("outbound").AtSliceIndex(0).AtMapKey("label"),
							knownvalue.StringExact("reject-http"),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("outbound").AtSliceIndex(0).AtMapKey("action"),
							knownvalue.StringExact("DROP"),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("outbound").AtSliceIndex(0).AtMapKey("protocol"),
							knownvalue.StringExact("TCP"),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("outbound").AtSliceIndex(0).AtMapKey("ports"),
							knownvalue.StringExact("80"),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("outbound").AtSliceIndex(0).AtMapKey("ipv4"),
							knownvalue.SetSizeExact(1),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("outbound").AtSliceIndex(0).AtMapKey("ipv4").AtSliceIndex(0),
							knownvalue.StringExact("0.0.0.0/0"),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("outbound").AtSliceIndex(0).AtMapKey("ipv6"),
							knownvalue.SetSizeExact(1),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("outbound").AtSliceIndex(0).AtMapKey("ipv6").AtSliceIndex(0),
							knownvalue.StringExact("::/0"),
						),

						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("devices"),
							knownvalue.SetSizeExact(2),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("devices").AtSliceIndex(0).AtMapKey("label"),
							knownvalue.NotNull(),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("devices").AtSliceIndex(0).AtMapKey("type"),
							knownvalue.NotNull(),
						),

						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("linodes"),
							knownvalue.SetSizeExact(1),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("nodebalancers"),
							knownvalue.SetSizeExact(1),
						),
						statecheck.ExpectKnownValue(
							testFirewallDataName,
							tfjsonpath.New("firewalls").AtSliceIndex(0).AtMapKey("interfaces"),
							knownvalue.SetSizeExact(0),
						),
					},
				},
			},
		})
	})
}
