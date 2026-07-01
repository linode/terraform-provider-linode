//go:build integration || vpcdefaultranges

package vpcdefaultranges_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/vpcdefaultranges/tmpl"
)

func TestAccDataSourceVPCDefaultRanges_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_vpc_default_ranges.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "default_ipv4_ranges.#"),
					resource.TestCheckResourceAttrSet(resourceName, "forbidden_ipv4_ranges.#"),
				),
			},
		},
	})
}
