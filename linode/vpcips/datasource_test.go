//go:build integration || vpcips

package vpcips_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/vpcips/tmpl"
)

func TestAccDataSourceVPCIPs_basic(t *testing.T) {
	t.Parallel()

	const (
		resourceNameAll    = "data.linode_vpc_ips.foobar"
		resourceNameScoped = "data.linode_vpc_ips.barfoo"
	)

	vpcLabel := acctest.RandomWithPrefix("tf-test")
	testRegion, err := acceptance.GetRandomRegionWithCaps([]string{"VPCs"}, "core")
	if err != nil {
		log.Fatal(fmt.Errorf("Error getting region: %s", err))
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, vpcLabel, testRegion, "10.0.0.0/24", "10.0.1.0/24"),
				ConfigStateChecks: []statecheck.StateCheck{
					// Index 0
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("gateway"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("linode_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("prefix"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("region"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("subnet_mask"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("subnet_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("config_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("interface_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("vpc_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("active"),
						knownvalue.NotNull(),
					),

					// index 1
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(1).AtMapKey("gateway"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(1).AtMapKey("linode_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(1).AtMapKey("prefix"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(1).AtMapKey("region"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(1).AtMapKey("subnet_mask"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(1).AtMapKey("subnet_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(1).AtMapKey("config_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(1).AtMapKey("interface_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(1).AtMapKey("vpc_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(1).AtMapKey("active"),
						knownvalue.NotNull(),
					),

					// scoped
					statecheck.ExpectKnownValue(
						resourceNameScoped,
						tfjsonpath.New("vpc_ips"),
						knownvalue.ListSizeExact(1),
					),
					statecheck.ExpectKnownValue(
						resourceNameScoped,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("address"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameScoped,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("gateway"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameScoped,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("linode_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameScoped,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("prefix"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameScoped,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("region"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameScoped,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("subnet_mask"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameScoped,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("subnet_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameScoped,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("config_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameScoped,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("interface_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameScoped,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("vpc_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameScoped,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("active"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	},
	)
}

func TestAccDataSourceVPCIPs_dualStack(t *testing.T) {
	t.Parallel()

	const (
		resourceNameScoped = "data.linode_vpc_ips.scoped"
		resourceNameAll    = "data.linode_vpc_ips.all"
	)

	vpcLabel := acctest.RandomWithPrefix("tf-test")

	// TODO (VPC Dual Stack): Remove region hardcoding
	targetRegion := "no-osl-1"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataDualStack(t, vpcLabel, targetRegion, "10.0.0.0/24"),
			},
			{
				Config: tmpl.DataDualStack(t, vpcLabel, targetRegion, "10.0.0.0/24"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("gateway"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("linode_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("prefix"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("region"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("subnet_mask"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("subnet_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("config_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("interface_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("vpc_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("active"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("ipv6_range"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("ipv6_range"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceNameAll,
						tfjsonpath.New("vpc_ips").AtSliceIndex(0).AtMapKey("ipv6_addresses"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}
