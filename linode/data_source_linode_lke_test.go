package linode

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceLinodeLKE(t *testing.T) {
	t.Parallel()
	resourceName := "data.linode_lke.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeLKEConfigBasic("k8s-acc") + testDataSourceLinodeLKE(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "label", "k8s-acc"),
					resource.TestCheckResourceAttr(resourceName, "region", "us-central"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.16"),
				),
			},
		},
	})
}

// TODO(sgmac): test passes, destroy leaves linode nodes behind even though
// the cluster is destroyed.
func testAccCheckLinodeLKEConfigBasic(label string) string {
	return fmt.Sprintf(`
resource "linode_lke" "foobar" {
	label = "%s"
	region = "us-central"
	version = "1.16"
	node_pools = [
		{ "count" = 3, "type" = "g6-standard-2"}
	]
}`, label)
}
func testDataSourceLinodeLKE() string {
	return fmt.Sprintf(`
data "linode_lke" "foobar" {
	id = "${linode_lke.foobar.id}"
}`)
}
