//go:build integration || reservedip

package reservedip_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/reservedip/tmpl"
)

func TestAccDataSource_reservedIP(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_reserved_ip.test"
	region, _ := acceptance.GetRandomRegionWithCaps([]string{"linodes"}, "core")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "address"),
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttrSet(resourceName, "gateway"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_mask"),
					resource.TestCheckResourceAttrSet(resourceName, "prefix"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttrSet(resourceName, "public"),
					resource.TestCheckResourceAttrSet(resourceName, "rdns"),
					resource.TestCheckResourceAttrSet(resourceName, "linode_id"),
					resource.TestCheckResourceAttrSet(resourceName, "reserved"),
				),
			},
		},
	})
}
