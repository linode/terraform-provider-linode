package linode

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceLinodeLKECluster(t *testing.T) {
	t.Parallel()
	resourceName := "data.linode_lke.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeLKEClusterConfigBasic("k8s-acc") + testDataSourceLinodeLKECluster(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "label", "k8s-acc"),
					resource.TestCheckResourceAttr(resourceName, "region", "us-central"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.16"),
				),
			},
		},
	})
}

func testDataSourceLinodeLKECluster() string {
	return fmt.Sprintf(`
data "linode_lke_cluster" "foobar" {
	id = "${linode_lke.foobar.id}"
}`)
}
