//go:build integration || reservedip

package reservedip_test

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/reservedip/tmpl"
)

const testResourceName = "linode_reserved_ip.test"

var testRegion string

func init() {
	region, err := acceptance.GetRandomRegionWithCaps([]string{linodego.CapabilityLinodes}, "core")
	if err != nil {
		log.Fatal(err)
	}
	testRegion = region
}

func TestAccResourceReservedIP_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, testRegion),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						testResourceName,
						tfjsonpath.New("address"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						testResourceName,
						tfjsonpath.New("region"),
						knownvalue.StringExact(testRegion),
					),
					statecheck.ExpectKnownValue(
						testResourceName,
						tfjsonpath.New("reserved"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						testResourceName,
						tfjsonpath.New("public"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						testResourceName,
						tfjsonpath.New("type"),
						knownvalue.StringExact("ipv4"),
					),
					statecheck.ExpectKnownValue(
						testResourceName,
						tfjsonpath.New("gateway"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						testResourceName,
						tfjsonpath.New("subnet_mask"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						testResourceName,
						tfjsonpath.New("prefix"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						testResourceName,
						tfjsonpath.New("tags"),
						knownvalue.SetSizeExact(0),
					),
					statecheck.ExpectKnownValue(
						testResourceName,
						tfjsonpath.New("vpc_nat_1_1"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						testResourceName,
						tfjsonpath.New("assigned_entity"),
						knownvalue.Null(),
					),
				},
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceReservedIP_withTags(t *testing.T) {
	t.Parallel()

	tagsInitial := []string{"tf-test", "reserved"}
	tagsUpdated := []string{"tf-test", "updated"}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.WithTags(t, testRegion, tagsInitial),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						testResourceName,
						tfjsonpath.New("tags"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("tf-test"),
							knownvalue.StringExact("reserved"),
						}),
					),
				},
			},
			{
				Config: tmpl.WithTags(t, testRegion, tagsUpdated),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						testResourceName,
						tfjsonpath.New("tags"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("tf-test"),
							knownvalue.StringExact("updated"),
						}),
					),
				},
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
