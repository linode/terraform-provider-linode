package region_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

func TestAccDataSourceRegion_basic(t *testing.T) {
	t.Parallel()

	country := "us"
	regionID := "us-east"
	resourceName := "data.linode_region.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfigBasic(regionID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "country", country),
					resource.TestCheckResourceAttr(resourceName, "id", regionID),
				),
			},
		},
	})
}

func dataSourceConfigBasic(regionID string) string {
	return fmt.Sprintf(`
data "linode_region" "foobar" {
	id = "%s"
}`, regionID)
}
