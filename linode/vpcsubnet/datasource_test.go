//go:build integration || vpcsubnet

package vpcsubnet_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/vpcsubnet/tmpl"
)

func TestAccDataSourceVPCSubnet_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_vpc_subnet.foo"
	subnetLabel := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, subnetLabel, "10.0.0.0/24", testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "label"),
					resource.TestCheckResourceAttrSet(resourceName, "ipv4"),
					resource.TestCheckResourceAttrSet(resourceName, "created"),
					resource.TestCheckResourceAttrSet(resourceName, "updated"),

					resource.TestCheckResourceAttr(resourceName, "linodes.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "linodes.0.id"),
					resource.TestCheckResourceAttr(resourceName, "linodes.0.interfaces.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "linodes.0.interfaces.0.id"),
					resource.TestCheckResourceAttr(resourceName, "linodes.0.interfaces.0.active", "false"),
				),
			},
		},
	})
}
