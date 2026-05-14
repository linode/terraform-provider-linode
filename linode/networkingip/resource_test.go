//go:build integration || networkingip

package networkingip_test

import (
	"log"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/networkingip/tmpl"
)

func init() {
	region, err := acceptance.GetRandomRegionWithCaps([]string{linodego.CapabilityLinodes}, "core")
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccResourceNetworkingIP_ephemeral(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")

	resourceName := "linode_networking_ip.reserved_ip"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.NetworkingIPReservedAssigned(t, label, testRegion, 0, false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("reserved"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("public"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("type"),
						knownvalue.StringExact("ipv4"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("address"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("gateway"),
						knownvalue.NotNull(),
					),
					statecheck.CompareValuePairs(
						resourceName,
						tfjsonpath.New("linode_id"),
						"linode_instance.test[0]",
						tfjsonpath.New("id"),
						helper.TypeAgnosticComparer(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("prefix"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("region"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rdns"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("subnet_mask"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("vpc_nat_1_1"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("tags"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("assigned_entity"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func TestAccResourceNetworkingIP_reserved(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")

	resourceName := "linode_networking_ip.reserved_ip"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.NetworkingIPReservedAssigned(t, label, testRegion, 0, true),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("reserved"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("public"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("type"),
						knownvalue.StringExact("ipv4"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("address"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("gateway"),
						knownvalue.NotNull(),
					),
					statecheck.CompareValuePairs(
						resourceName,
						tfjsonpath.New("linode_id"),
						"linode_instance.test[0]",
						tfjsonpath.New("id"),
						helper.TypeAgnosticComparer(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("prefix"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("region"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rdns"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("subnet_mask"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("vpc_nat_1_1"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("tags"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("assigned_entity"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func TestAccResourceNetworkingIP_reservedEphemeralReassignment(t *testing.T) {
	t.Parallel()

	resName := "linode_networking_ip.reserved_ip"
	linodeLabel := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Create an assigned reserved IP
			{
				Config: tmpl.NetworkingIPReservedAssigned(
					t,
					linodeLabel,
					testRegion,
					0,
					true,
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs(
						resName,
						tfjsonpath.New("linode_id"),
						"linode_instance.test[0]",
						tfjsonpath.New("id"),
						helper.TypeAgnosticComparer(),
					),
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("reserved"),
						knownvalue.Bool(true),
					),
				},
			},
			// Make the IP ephemeral
			{
				Config: tmpl.NetworkingIPReservedAssigned(
					t,
					linodeLabel,
					testRegion,
					0,
					false,
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs(
						resName,
						tfjsonpath.New("linode_id"),
						"linode_instance.test[0]",
						tfjsonpath.New("id"),
						helper.TypeAgnosticComparer(),
					),
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("reserved"),
						knownvalue.Bool(false),
					),
				},
			},
			// Attempt to reassign the ephemeral IP; expect RequiresReplace
			{
				Config: tmpl.NetworkingIPReservedAssigned(
					t,
					linodeLabel,
					testRegion,
					1,
					false,
				),
				PlanOnly: true,
				// This plan is expected to trigger a RequiresReplace
				ExpectNonEmptyPlan: true,
			},
			// Convert back to a reserved IP, assign to second instance
			{
				Config: tmpl.NetworkingIPReservedAssigned(
					t,
					linodeLabel,
					testRegion,
					1,
					true,
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs(
						resName,
						tfjsonpath.New("linode_id"),
						"linode_instance.test[1]",
						tfjsonpath.New("id"),
						helper.TypeAgnosticComparer(),
					),
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("reserved"),
						knownvalue.Bool(true),
					),
				},
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"wait_for_available"},
			},
		},
	})
}

func TestAccResourceNetworkingIP_ephemeralToReservedConversion(t *testing.T) {
	t.Parallel()

	resName := "linode_networking_ip.reserved_ip"
	linodeLabel := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Step 1: Create an ephemeral IP
			{
				Config: tmpl.NetworkingIPReservedAssigned(
					t,
					linodeLabel,
					testRegion,
					0,
					false,
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("reserved"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("tags"),
						knownvalue.NotNull(),
					),
				},
			},
			// Step 2: Convert ephemeral → reserved (in-place, no replacement)
			{
				Config: tmpl.NetworkingIPReservedAssigned(
					t,
					linodeLabel,
					testRegion,
					0,
					true,
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("reserved"),
						knownvalue.Bool(true),
					),
					statecheck.CompareValuePairs(
						resName,
						tfjsonpath.New("linode_id"),
						"linode_instance.test[0]",
						tfjsonpath.New("id"),
						helper.TypeAgnosticComparer(),
					),
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("tags"),
						knownvalue.NotNull(),
					),
				},
			},
			// Step 3: Convert reserved → ephemeral (in-place, no replacement)
			{
				Config: tmpl.NetworkingIPReservedAssigned(
					t,
					linodeLabel,
					testRegion,
					0,
					false,
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("reserved"),
						knownvalue.Bool(false),
					),
					statecheck.CompareValuePairs(
						resName,
						tfjsonpath.New("linode_id"),
						"linode_instance.test[0]",
						tfjsonpath.New("id"),
						helper.TypeAgnosticComparer(),
					),
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("tags"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func TestAccResourceNetworkingIP_reservedUnassignedToEphemeral(t *testing.T) {
	t.Parallel()

	resName := "linode_networking_ip.reserved_ip"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create a reserved unassigned IP.
			{
				Config: tmpl.NetworkingIPReservedUnassigned(t, testRegion, true),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("reserved"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("linode_id"), knownvalue.Null()),
				},
			},
			// Step 2: Convert reserved → ephemeral. Because the IP is unassigned,
			// the API deletes it. The provider surfaces an error so the user knows
			// to remove the resource from their configuration.
			{
				Config:      tmpl.NetworkingIPReservedUnassigned(t, testRegion, false),
				ExpectError: regexp.MustCompile(`IP Address Deleted During Update`),
			},
		},
	})
}
