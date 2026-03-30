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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "address"),
					resource.TestCheckResourceAttr(testResourceName, "region", testRegion),
					resource.TestCheckResourceAttr(testResourceName, "reserved", "true"),
					resource.TestCheckResourceAttr(testResourceName, "public", "true"),
					resource.TestCheckResourceAttr(testResourceName, "type", "ipv4"),
					resource.TestCheckResourceAttrSet(testResourceName, "gateway"),
					resource.TestCheckResourceAttrSet(testResourceName, "subnet_mask"),
					resource.TestCheckResourceAttrSet(testResourceName, "prefix"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
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
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags"},
			},
		},
	})
}

func TestAccResourceReservedIP_withTags(t *testing.T) {
	t.Parallel()

	tags := []string{"tf-test", "reserved"}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.WithTags(t, testRegion, tags),
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
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags"},
			},
		},
	})
}
