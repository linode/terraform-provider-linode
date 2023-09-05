//go:build integration

package vpcsubnet_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/vpcsubnet/tmpl"
)

func TestAccDataSourceUser_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_vpc_subnet.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "label"),
					resource.TestCheckResourceAttrSet(resourceName, "ipv4"),
					resource.TestCheckResourceAttrSet(resourceName, "linodes"),
					resource.TestCheckResourceAttrSet(resourceName, "created"),
					resource.TestCheckResourceAttrSet(resourceName, "updated"),
				),
			},
		},
	})
}
