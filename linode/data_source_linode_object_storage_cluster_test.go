package linode

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLinodeObjectStorageCluster_basic(t *testing.T) {
	t.Parallel()

	objectStorageClusterID := "us-east-1"
	region := "us-east"
	resourceName := "data.linode_object_storage_cluster.foobar"
	staticSiteDomain := "website-us-east-1.linodeobjects.com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLinodeObjectStorageClusterBasic(objectStorageClusterID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "region", region),
					resource.TestCheckResourceAttr(resourceName, "id", objectStorageClusterID),
					resource.TestCheckResourceAttr(resourceName, "static_site_domain", staticSiteDomain),
				),
			},
		},
	})
}

func testDataSourceLinodeObjectStorageClusterBasic(objectStorageClusterID string) string {
	return fmt.Sprintf(`
data "linode_object_storage_cluster" "foobar" {
    id = "%s"
}`, objectStorageClusterID)
}
