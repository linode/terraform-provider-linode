//go:build integration || vpcsubnets

package vpcsubnets_test

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/vpcsubnets/tmpl"
)

func TestSmokeTests_vpcsubnets(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{"TestAccDataSourceVPCSubnets_basic_smoke", TestAccDataSourceVPCSubnets_basic_smoke},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestAccDataSourceVPCSubnets_basic_smoke(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_vpc_subnets.foobar"
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
				Config: tmpl.DataBasic(t, vpcLabel, testRegion, "10.0.0.0/24"),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckResourceAttrGreaterThan(resourceName, "vpc_subnets.#", 0),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("vpc_subnets").AtSliceIndex(0).AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("vpc_subnets").AtSliceIndex(0).AtMapKey("label"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("vpc_subnets").AtSliceIndex(0).AtMapKey("ipv4"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("vpc_subnets").AtSliceIndex(0).AtMapKey("created"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("vpc_subnets").AtSliceIndex(0).AtMapKey("updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("vpc_subnets").AtSliceIndex(0).AtMapKey("linodes").AtSliceIndex(0).AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("vpc_subnets").AtSliceIndex(0).AtMapKey("linodes").AtSliceIndex(0).AtMapKey("interfaces").AtSliceIndex(0).AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("vpc_subnets").AtSliceIndex(0).AtMapKey("linodes").AtSliceIndex(0).AtMapKey("interfaces").AtSliceIndex(0).AtMapKey("active"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("vpc_subnets").AtSliceIndex(0).AtMapKey("linodes").AtSliceIndex(0).AtMapKey("interfaces").AtSliceIndex(0).AtMapKey("config_id"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccDataSourceVPCSubnets_filterByLabel(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_vpc_subnets.foobar"
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
				Config: tmpl.DataFilterLabel(t, vpcLabel, testRegion, "10.0.0.0/24"),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckResourceAttrGreaterThan(resourceName, "vpc_subnets.#", 0),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("vpc_subnets").AtSliceIndex(0).AtMapKey("label"), knownvalue.StringRegexp(regexp.MustCompile("tf-test"))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("vpc_subnets").AtSliceIndex(0).AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("vpc_subnets").AtSliceIndex(0).AtMapKey("ipv4"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("vpc_subnets").AtSliceIndex(0).AtMapKey("created"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("vpc_subnets").AtSliceIndex(0).AtMapKey("updated"), knownvalue.NotNull()),
				},
			},
		},
	})
}
