//go:build integration || regionsvpcavailability

package regionsvpcavailability_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/regionsvpcavailability/tmpl"
)

func TestAccDataSourceRegionsVPCAvailability_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_regions_vpc_availability.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("regions_vpc_availability"), knownvalue.ListSizeExact(32)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("regions_vpc_availability"), knownvalue.ListPartial(map[int]knownvalue.Check{
						0: knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":        knownvalue.StringExact("nl-ams"),
							"available": knownvalue.Bool(true),
							"available_ipv6_prefix_lengths": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.Int64Exact(48),
								knownvalue.Int64Exact(52),
							}),
						}),
						31: knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":                            knownvalue.StringExact("us-west"),
							"available":                     knownvalue.Bool(false),
							"available_ipv6_prefix_lengths": knownvalue.ListSizeExact(0),
						}),
					})),
				},
			},
		},
	})
}
