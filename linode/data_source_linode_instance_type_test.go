package linode

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLinodeInstanceType_basic(t *testing.T) {
	t.Parallel()

	instanceTypeID := "g6-standard-2"
	resourceName := "data.linode_instance_type.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLinodeInstanceTypeBasic(instanceTypeID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", instanceTypeID),
					resource.TestCheckResourceAttr(resourceName, "label", "Linode 4GB"),
					resource.TestCheckResourceAttr(resourceName, "disk", "81920"),
					resource.TestCheckResourceAttr(resourceName, "class", "standard"),
					resource.TestCheckResourceAttr(resourceName, "memory", "4096"),
					resource.TestCheckResourceAttr(resourceName, "vcpus", "2"),
					resource.TestCheckResourceAttr(resourceName, "network_out", "4000"),
					resource.TestCheckResourceAttr(resourceName, "price.0.hourly", "0.029999999329447746"),
					resource.TestCheckResourceAttr(resourceName, "price.0.monthly", "20"),
					resource.TestCheckResourceAttr(resourceName, "addons.0.backups.0.price.0.hourly", "0.00800000037997961"),
					resource.TestCheckResourceAttr(resourceName, "addons.0.backups.0.price.0.monthly", "5"),
				),
			},
		},
	})
}

func testDataSourceLinodeInstanceTypeBasic(instanceTypeID string) string {
	return fmt.Sprintf(`
data "linode_instance_type" "foobar" {
	id = "%s"
}`, instanceTypeID)
}
