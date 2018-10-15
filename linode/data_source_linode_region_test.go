package linode

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestDataSourceLinodeRegion(t *testing.T) {
	t.Parallel()

	country := "us"
	regionID := "us-east"
	resourceName := "data.linode_region.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLinodeRegion(regionID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "country", country),
					resource.TestCheckResourceAttr(resourceName, "id", regionID),
				),
			},
		},
	})
}

func testDataSourceLinodeRegion(regionID string) string {
	return fmt.Sprintf(`
data "linode_region" "foobar" {
	id = "%s"
}`, regionID)
}
