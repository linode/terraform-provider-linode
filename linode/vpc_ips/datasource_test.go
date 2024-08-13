//go:build integration || vpc_ips

package vpc_ips_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/vpc_ips/tmpl"
)

func TestAccDataSourceVPCIPs_basic(t *testing.T) {
	t.Parallel()

	resourceName_all_ips := "data.linode_vpc_ips.foobar"
	resourceName_vpc_ips := "data.linode_vpc_ips.barfoo"
	resourceName_filtered_ips := "data.linode_vpc_ips.foobar_filter"

	vpcLabel := acctest.RandomWithPrefix("tf-test")
	testRegion, err := acceptance.GetRandomRegionWithCaps([]string{"VPCs"})
	if err != nil {
		log.Fatal(fmt.Errorf("Error getting region: %s", err))
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, vpcLabel, testRegion, "10.0.0.0/24", "10.0.1.0/24"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName_all_ips, "vpc_ips.#", "2"),
					resource.TestCheckResourceAttr(resourceName_vpc_ips, "vpc_ips.#", "1"),
					resource.TestCheckResourceAttr(resourceName_filtered_ips, "vpc_ips.#", "1"),

					resource.TestCheckResourceAttrSet(resourceName_all_ips, "vpc_ips.0.address"),
					resource.TestCheckResourceAttrSet(resourceName_all_ips, "vpc_ips.0.gateway"),
					resource.TestCheckResourceAttrSet(resourceName_all_ips, "vpc_ips.0.linode_id"),
					resource.TestCheckResourceAttrSet(resourceName_all_ips, "vpc_ips.0.prefix"),
					resource.TestCheckResourceAttrSet(resourceName_all_ips, "vpc_ips.0.region"),
					resource.TestCheckResourceAttrSet(resourceName_all_ips, "vpc_ips.0.subnet_mask"),
					resource.TestCheckResourceAttrSet(resourceName_all_ips, "vpc_ips.0.nat_1_1"),
					resource.TestCheckResourceAttrSet(resourceName_all_ips, "vpc_ips.0.subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName_all_ips, "vpc_ips.0.config_id"),
					resource.TestCheckResourceAttrSet(resourceName_all_ips, "vpc_ips.0.interface_id"),
					resource.TestCheckNoResourceAttr(resourceName_all_ips, "vpc_ips.0.address_range"),
					resource.TestCheckResourceAttrSet(resourceName_all_ips, "vpc_ips.0.vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName_all_ips, "vpc_ips.0.active"),

					resource.TestCheckResourceAttrSet(resourceName_all_ips, "vpc_ips.1.address"),
					resource.TestCheckResourceAttrSet(resourceName_all_ips, "vpc_ips.1.gateway"),
					resource.TestCheckResourceAttrSet(resourceName_all_ips, "vpc_ips.1.linode_id"),
					resource.TestCheckResourceAttrSet(resourceName_all_ips, "vpc_ips.1.prefix"),
					resource.TestCheckResourceAttrSet(resourceName_all_ips, "vpc_ips.1.region"),
					resource.TestCheckResourceAttrSet(resourceName_all_ips, "vpc_ips.1.subnet_mask"),
					resource.TestCheckNoResourceAttr(resourceName_all_ips, "vpc_ips.1.nat_1_1"),
					resource.TestCheckResourceAttrSet(resourceName_all_ips, "vpc_ips.1.subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName_all_ips, "vpc_ips.1.config_id"),
					resource.TestCheckResourceAttrSet(resourceName_all_ips, "vpc_ips.1.interface_id"),
					resource.TestCheckNoResourceAttr(resourceName_all_ips, "vpc_ips.1.address_range"),
					resource.TestCheckResourceAttrSet(resourceName_all_ips, "vpc_ips.1.vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName_all_ips, "vpc_ips.1.active"),
				),
			},
		},
	})
}
