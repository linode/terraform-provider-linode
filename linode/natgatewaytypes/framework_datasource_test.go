//go:build integration || natgatewaytypes

package natgatewaytypes_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/natgatewaytypes/tmpl"
)

func TestAccDataSourceNATGatewayTypes_basic(t *testing.T) {
	t.Parallel()

	dataSourceName := "data.linode_natgateway_types.foobar"
	byLabelName := "data.linode_natgateway_types.by_label"
	orderedName := "data.linode_natgateway_types.ordered"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						dataSourceName,
						tfjsonpath.New("types"),
						knownvalue.ListSizeExact(1),
					),
					statecheck.ExpectKnownValue(
						dataSourceName,
						tfjsonpath.New("types").AtSliceIndex(0).AtMapKey("id"),
						knownvalue.StringExact("g1-natgateway"),
					),
					statecheck.ExpectKnownValue(
						dataSourceName,
						tfjsonpath.New("types").AtSliceIndex(0).AtMapKey("label"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dataSourceName,
						tfjsonpath.New("types").AtSliceIndex(0).AtMapKey("price").AtSliceIndex(0).AtMapKey("hourly"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dataSourceName,
						tfjsonpath.New("types").AtSliceIndex(0).AtMapKey("price").AtSliceIndex(0).AtMapKey("monthly"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dataSourceName,
						tfjsonpath.New("types").AtSliceIndex(0).AtMapKey("region_prices"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dataSourceName,
						tfjsonpath.New("types").AtSliceIndex(0).AtMapKey("transfer"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						byLabelName,
						tfjsonpath.New("types").AtSliceIndex(0).AtMapKey("label"),
						knownvalue.StringExact("NAT Gateway"),
					),
					statecheck.ExpectKnownValue(
						orderedName,
						tfjsonpath.New("types").AtSliceIndex(0).AtMapKey("id"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}
