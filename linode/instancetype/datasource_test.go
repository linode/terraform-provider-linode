package instancetype_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/instancetype/tmpl"
)

func TestAccDataSourceLinodeInstanceType_basic(t *testing.T) {
	t.Parallel()

	instanceTypeID := "g6-standard-2"
	resourceName := "data.linode_instance_type.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, instanceTypeID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", instanceTypeID),
					resource.TestCheckResourceAttr(resourceName, "label", "Linode 4GB"),
					resource.TestCheckResourceAttr(resourceName, "disk", "81920"),
					resource.TestCheckResourceAttr(resourceName, "class", "standard"),
					resource.TestCheckResourceAttr(resourceName, "memory", "4096"),
					resource.TestCheckResourceAttr(resourceName, "vcpus", "2"),
					resource.TestCheckResourceAttr(resourceName, "network_out", "4000"),
					resource.TestCheckResourceAttr(resourceName, "price.0.hourly", "0.035999998450279236"),
					resource.TestCheckResourceAttr(resourceName, "price.0.monthly", "20"),
					resource.TestCheckResourceAttr(resourceName, "addons.0.backups.0.price.0.hourly", "0.00800000037997961"),
					resource.TestCheckResourceAttr(resourceName, "addons.0.backups.0.price.0.monthly", "5"),
				),
			},
		},
	})
}
