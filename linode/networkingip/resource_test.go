//go:build integration || networkingip

package networkingip_test

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/networkingip/tmpl"
)

func init() {
	region, err := acceptance.GetRandomRegionWithCaps([]string{"linodes"}, "core")
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
