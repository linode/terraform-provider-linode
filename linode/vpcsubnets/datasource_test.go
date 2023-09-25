//go:build integration

package vpcsubnets_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/vpcsubnets/tmpl"
)

func TestAccDataSourceVPCSubnets_basic_smoke(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_vpc_subnets.foobar"
	vpcLabel := acctest.RandomWithPrefix("tf-test")
	testRegion := "us-east"

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
				),
			},
		},
	})
}

func TestAccDataSourceVPCSubnets_filterByLabel(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_vpc_subnets.foobar"
	vpcLabel := acctest.RandomWithPrefix("tf-test")
	testRegion := "us-east"

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
