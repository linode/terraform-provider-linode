//go:build integration || vpcsubnets

package vpcsubnets_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/vpcsubnets/tmpl"
)

func TestAccDataSourceVPCSubnets_basic_smoke(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_vpc_subnets.foobar"
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
				Config: tmpl.DataBasic(t, vpcLabel, testRegion, "10.0.0.0/24"),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckResourceAttrGreaterThan(resourceName, "vpc_subnets.#", 0),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_subnets.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_subnets.0.label"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_subnets.0.ipv4"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_subnets.0.created"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_subnets.0.updated"),

					resource.TestCheckResourceAttrSet(resourceName, "vpc_subnets.0.linodes.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_subnets.0.linodes.0.interfaces.0.id"),
					resource.TestCheckResourceAttr(resourceName, "vpc_subnets.0.linodes.0.interfaces.0.active", "false"),
				),
			},
		},
	})
}

func TestAccDataSourceVPCSubnets_filterByLabel(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_vpc_subnets.foobar"
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
				Config: tmpl.DataFilterLabel(t, vpcLabel, testRegion, "10.0.0.0/24"),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckResourceAttrGreaterThan(resourceName, "vpc_subnets.#", 0),
					acceptance.CheckResourceAttrContains(resourceName, "vpc_subnets.0.label", "tf-test"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_subnets.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_subnets.0.ipv4"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_subnets.0.created"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_subnets.0.updated"),
				),
			},
		},
	})
}
