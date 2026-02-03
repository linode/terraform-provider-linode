package regionsvpcavailability_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					acceptance.CheckResourceAttrGreaterThan(resourceName, "regions_vpc_availability.#", 0),
					resource.TestCheckResourceAttrSet(resourceName, "regions_vpc_availability.0.available"),
					resource.TestCheckResourceAttrSet(resourceName, "regions_vpc_availability.0.id"),
					resource.TestCheckResourceAttr(resourceName, "regions_vpc_availability.0.available_ipv6_prefix_lengths.#", "0"),
				),
			},
		},
	})
}
