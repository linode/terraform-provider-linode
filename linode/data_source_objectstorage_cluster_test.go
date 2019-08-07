package linode

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceLinodeObjectStorageCluster(t *testing.T) {
	t.Parallel()

	objectStorageClusterID := "us-east-1"
	region := "us-east"
	resourceName := "data.linode_objectstorage_cluster.foobar"
	staticSiteDomain := "website-us-east-1.linodeobjects.com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLinodeObjectStorageCluster(objectStorageClusterID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "region", region),
					resource.TestCheckResourceAttr(resourceName, "id", objectStorageClusterID),
					resource.TestCheckResourceAttr(resourceName, "static_site_domain", staticSiteDomain),
				),
			},
		},
	})
}

func testDataSourceLinodeObjectStorageCluster(objectStorageClusterID string) string {
	return fmt.Sprintf(`
data "linode_objectstorage_cluster" "foobar" {
    id = "%s"
}`, objectStorageClusterID)
}
