package instancetypes_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/instancetypes/tmpl"
)

func TestAccDataSourceInstanceTypes_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_instance_types.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "types.0.id", "g6-standard-2"),
					resource.TestCheckResourceAttr(resourceName, "types.0.label", "Linode 4GB"),
					resource.TestCheckResourceAttr(resourceName, "types.0.class", "standard"),
					resource.TestCheckResourceAttrSet(resourceName, "types.0.disk"),
					resource.TestCheckResourceAttrSet(resourceName, "types.0.network_out"),
					resource.TestCheckResourceAttrSet(resourceName, "types.0.memory"),
					resource.TestCheckResourceAttrSet(resourceName, "types.0.transfer"),
					resource.TestCheckResourceAttrSet(resourceName, "types.0.vcpus"),
					resource.TestCheckResourceAttrSet(resourceName, "types.0.price.0.hourly"),
					resource.TestCheckResourceAttrSet(resourceName, "types.0.price.0.monthly"),
					resource.TestCheckResourceAttrSet(resourceName, "types.0.addons.0.backups.0.price.0.hourly"),
					resource.TestCheckResourceAttrSet(resourceName, "types.0.addons.0.backups.0.price.0.monthly"),

					// Ensure order is correctly appended to filter
					acceptance.CheckResourceAttrContains(resourceName, "id", "\"+order_by\":\"vcpus\""),
					acceptance.CheckResourceAttrContains(resourceName, "id", "\"+order\":\"desc\""),
				),
			},
		},
	})
}
