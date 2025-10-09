//go:build integration || vpc

package vpc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/vpc/tmpl"
)

func TestAccDataSourceVPC_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_vpc.foo"
	vpcLabel := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, vpcLabel, testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "label", vpcLabel),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttrSet(resourceName, "created"),
					resource.TestCheckResourceAttrSet(resourceName, "updated"),
				),
			},
		},
	})
}

func TestAccDataSourceVPC_dualStack(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_vpc.foo"
	vpcLabel := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// TODO (VPC Dual Stack): Remove region hardcoding
				Config: tmpl.DataDualStack(t, vpcLabel, "no-osl-1"),
				// ... existing code ...
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("label"),
						knownvalue.StringExact(vpcLabel),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("description"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("region"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("created"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("updated"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("ipv6"),
						knownvalue.ListSizeExact(1),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("ipv6").AtSliceIndex(0).AtMapKey("range"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}
