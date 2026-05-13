//go:build integration || networkingip

package networkingips_test

import (
	"log"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/networkingips/tmpl"
)

var testRegion string

func init() {
	region, err := acceptance.GetRandomRegionWithCaps([]string{linodego.CapabilityLinodes}, "core")
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccDataSourceNetworkingIP_list(t *testing.T) {
	t.Parallel()

	dataResourceName := "data.linode_networking_ips.list"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataList(t),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						dataResourceName,
						tfjsonpath.New("ip_addresses"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dataResourceName,
						tfjsonpath.New("ip_addresses").AtSliceIndex(0).AtMapKey("address"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dataResourceName,
						tfjsonpath.New("ip_addresses").AtSliceIndex(0).AtMapKey("region"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dataResourceName,
						tfjsonpath.New("ip_addresses").AtSliceIndex(0).AtMapKey("type"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dataResourceName,
						tfjsonpath.New("ip_addresses").AtSliceIndex(0).AtMapKey("tags"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func TestAccDataSourceNetworkingIP_filterReserved(t *testing.T) {
	t.Parallel()

	dataResourceName := "data.linode_networking_ips.filtered"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataFilterReserved(t),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dataResourceName, tfjsonpath.New("ip_addresses"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataResourceName, tfjsonpath.New("ip_addresses").AtSliceIndex(0).AtMapKey("reserved"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(dataResourceName, tfjsonpath.New("ip_addresses").AtSliceIndex(0).AtMapKey("address"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataResourceName, tfjsonpath.New("ip_addresses").AtSliceIndex(0).AtMapKey("interface_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue(dataResourceName, tfjsonpath.New("ip_addresses").AtSliceIndex(0).AtMapKey("region"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(
						dataResourceName,
						tfjsonpath.New("ip_addresses").AtSliceIndex(0).AtMapKey("gateway"),
						knownvalue.StringRegexp(regexp.MustCompile(`\.1$`)),
					),
					statecheck.ExpectKnownValue(
						dataResourceName,
						tfjsonpath.New("ip_addresses").AtSliceIndex(0).AtMapKey("type"),
						knownvalue.StringExact("ipv4"),
					),
					statecheck.ExpectKnownValue(dataResourceName, tfjsonpath.New("ip_addresses").AtSliceIndex(0).AtMapKey("public"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(dataResourceName, tfjsonpath.New("ip_addresses").AtSliceIndex(0).AtMapKey("prefix"), knownvalue.Int64Exact(24)),
					statecheck.ExpectKnownValue(dataResourceName, tfjsonpath.New("ip_addresses").AtSliceIndex(0).AtMapKey("subnet_mask"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(
						dataResourceName,
						tfjsonpath.New("ip_addresses").AtSliceIndex(0).AtMapKey("assigned_entity"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}
