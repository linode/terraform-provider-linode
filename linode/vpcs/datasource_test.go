//go:build integration || vpcs

package vpcs_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/vpcs/tmpl"
)

func TestAccDataSourceVPCs_basic_smoke(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_vpcs.foobar"
	vpcLabel := acctest.RandomWithPrefix("tf-test")
	testRegion, err := acceptance.GetRandomRegionWithCaps([]string{"VPCs"})
	if err != nil {
		t.Error(fmt.Errorf("failed to get region with VPC capability: %w", err))
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, vpcLabel, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckResourceAttrGreaterThan(resourceName, "vpcs.#", 0),
					resource.TestCheckResourceAttrSet(resourceName, "vpcs.0.label"),
					resource.TestCheckResourceAttrSet(resourceName, "vpcs.0.description"),
					resource.TestCheckResourceAttrSet(resourceName, "vpcs.0.region"),
					resource.TestCheckResourceAttrSet(resourceName, "vpcs.0.created"),
					resource.TestCheckResourceAttrSet(resourceName, "vpcs.0.updated"),
				),
			},
		},
	})
}

func TestAccDataSourceVPCs_filterByLabel(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_vpcs.foobar"
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
				Config: tmpl.DataFilterLabel(t, vpcLabel, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckResourceAttrGreaterThan(resourceName, "vpcs.#", 0),
					acceptance.CheckResourceAttrContains(resourceName, "vpcs.0.label", "tf-test"),
					resource.TestCheckResourceAttrSet(resourceName, "vpcs.0.region"),
					resource.TestCheckResourceAttrSet(resourceName, "vpcs.0.created"),
					resource.TestCheckResourceAttrSet(resourceName, "vpcs.0.updated"),
				),
			},
		},
	})
}
