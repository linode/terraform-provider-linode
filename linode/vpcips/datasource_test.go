//go:build integration || vpcips

package vpcips_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckResourceAttrGreaterThan(resourceNameAll, "vpc_ips.#", 0),
					resource.TestCheckResourceAttr(resourceNameScoped, "vpc_ips.#", "1"),

					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.0.address"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.0.gateway"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.0.linode_id"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.0.prefix"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.0.region"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.0.subnet_mask"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.0.subnet_id"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.0.config_id"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.0.interface_id"),
					resource.TestCheckNoResourceAttr(resourceNameAll, "vpc_ips.0.address_range"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.0.vpc_id"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.0.active"),

					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.1.address"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.1.gateway"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.1.linode_id"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.1.prefix"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.1.region"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.1.subnet_mask"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.1.subnet_id"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.1.config_id"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.1.interface_id"),
					resource.TestCheckNoResourceAttr(resourceNameAll, "vpc_ips.1.address_range"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.1.vpc_id"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.1.active"),

					resource.TestCheckResourceAttrSet(resourceNameScoped, "vpc_ips.0.address"),
					resource.TestCheckResourceAttrSet(resourceNameScoped, "vpc_ips.0.gateway"),
					resource.TestCheckResourceAttrSet(resourceNameScoped, "vpc_ips.0.linode_id"),
					resource.TestCheckResourceAttrSet(resourceNameScoped, "vpc_ips.0.prefix"),
					resource.TestCheckResourceAttrSet(resourceNameScoped, "vpc_ips.0.region"),
					resource.TestCheckResourceAttrSet(resourceNameScoped, "vpc_ips.0.subnet_mask"),
					resource.TestCheckResourceAttrSet(resourceNameScoped, "vpc_ips.0.subnet_id"),
					resource.TestCheckResourceAttrSet(resourceNameScoped, "vpc_ips.0.config_id"),
					resource.TestCheckResourceAttrSet(resourceNameScoped, "vpc_ips.0.interface_id"),
					resource.TestCheckNoResourceAttr(resourceNameScoped, "vpc_ips.0.address_range"),
					resource.TestCheckResourceAttrSet(resourceNameScoped, "vpc_ips.0.vpc_id"),
					resource.TestCheckResourceAttrSet(resourceNameScoped, "vpc_ips.0.active"),
				),
			},
		},
	})
}

func TestAccDataSourceVPCIPs_dualStack(t *testing.T) {
	// TODO (VPC Dual Stack): Finish test after interfaces readiness.
	t.Skip("TODO (VPC Dual Stack): Finish test after interfaces readiness.")

	t.Parallel()

	const (
		resourceNameScoped = "data.linode_vpc_ips.scoped"
		resourceNameAll    = "data.linode_vpc_ips.all"
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
				Config: tmpl.DataDualStack(t, vpcLabel, testRegion, "10.0.0.0/24"),
			},
			{
				Config: tmpl.DataDualStack(t, vpcLabel, testRegion, "10.0.0.0/24"),
				Check: resource.ComposeTestCheckFunc(
					// acceptance.CheckResourceAttrGreaterThan(resourceNameAll, "vpc_ips.#", 0),
					// acceptance.CheckResourceAttrGreaterThan(resourceNameScoped, "vpc_ips.#", 0),

					resource.TestCheckNoResourceAttr(resourceNameAll, "vpc_ips.0.address"),
					resource.TestCheckResourceAttr(resourceNameAll, "vpc_ips.0.gateway", ""),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.0.linode_id"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.0.prefix"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.0.region"),
					resource.TestCheckResourceAttr(resourceNameAll, "vpc_ips.0.subnet_mask", ""),
					resource.TestCheckResourceAttr(resourceNameAll, "vpc_ips.0.nat_1_1", ""),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.0.subnet_id"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.0.config_id"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.0.interface_id"),
					resource.TestCheckNoResourceAttr(resourceNameAll, "vpc_ips.0.address_range"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.0.vpc_id"),
					resource.TestCheckResourceAttrSet(resourceNameAll, "vpc_ips.0.active"),

					// resource.TestCheckNoResourceAttr(resourceNameScoped, "vpc_ips.0.address"),
					// resource.TestCheckResourceAttr(resourceNameScoped, "vpc_ips.0.gateway", ""),
					resource.TestCheckResourceAttrSet(resourceNameScoped, "vpc_ips.0.linode_id"),
					resource.TestCheckResourceAttrSet(resourceNameScoped, "vpc_ips.0.prefix"),
					resource.TestCheckResourceAttrSet(resourceNameScoped, "vpc_ips.0.region"),
					resource.TestCheckResourceAttr(resourceNameScoped, "vpc_ips.0.subnet_mask", ""),
					resource.TestCheckResourceAttr(resourceNameScoped, "vpc_ips.0.nat_1_1", ""),
					resource.TestCheckResourceAttrSet(resourceNameScoped, "vpc_ips.0.subnet_id"),
					resource.TestCheckResourceAttrSet(resourceNameScoped, "vpc_ips.0.config_id"),
					resource.TestCheckResourceAttrSet(resourceNameScoped, "vpc_ips.0.interface_id"),
					resource.TestCheckNoResourceAttr(resourceNameScoped, "vpc_ips.0.address_range"),
					resource.TestCheckResourceAttrSet(resourceNameScoped, "vpc_ips.0.vpc_id"),
					resource.TestCheckResourceAttrSet(resourceNameScoped, "vpc_ips.0.active"),
				),
			},
		},
	})
}
