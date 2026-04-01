//go:build integration || reservediptypes

package reservediptypes_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/reservediptypes/tmpl"
)

func TestAccDataSourceReservedIPTypes_basic(t *testing.T) {
	t.Parallel()

	dataSourceName := "data.linode_reserved_ip_types.foobar"

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
						knownvalue.NotNull(),
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
				},
			},
		},
	})
}
