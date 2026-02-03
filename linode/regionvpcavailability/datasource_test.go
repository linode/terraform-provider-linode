package regionvpcavailability_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/regionvpcavailability/tmpl"
)

func TestAccDataSourceRegionVPCAvailability_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_region_vpc_availability.test"
	regionID := "nl-ams"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, regionID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "available"),
					resource.TestCheckResourceAttr(resourceName, "available_ipv6_prefix_lengths.#", "0"),
				),
			},
			{
				Config:      tmpl.DataNoRegion(t),
				ExpectError: regexp.MustCompile(`\[404\] Not found`),
			},
		},
	})
}
