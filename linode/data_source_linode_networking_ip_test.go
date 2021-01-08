package linode

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLinodeNetworkingIP_basic(t *testing.T) {
	t.Parallel()

	resourceName := "linode_instance.foobar"
	dataResourceName := "data.linode_networking_ip.foobar"

	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLinodeNetworkingIPBasic(label),
			},
			{
				Config: testDataSourceLinodeNetworkingIPBasic(label),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataResourceName, "address", resourceName, "ip_address"),
					resource.TestCheckResourceAttrPair(dataResourceName, "linode_id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataResourceName, "region", resourceName, "region"),
					resource.TestMatchResourceAttr(dataResourceName, "gateway", regexp.MustCompile(`\.1$`)),
					resource.TestCheckResourceAttr(dataResourceName, "type", "ipv4"),
					resource.TestCheckResourceAttr(dataResourceName, "public", "true"),
					resource.TestCheckResourceAttr(dataResourceName, "prefix", "24"),
					resource.TestMatchResourceAttr(dataResourceName, "rdns", regexp.MustCompile(`\.members\.linode\.com$`)),
				),
			},
		},
	})
}

func testDataSourceLinodeNetworkingIPBasic(label string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	group = "tf_test"
	image = "linode/containerlinux"
	type = "g6-standard-1"
	region = "us-east"
}

data "linode_networking_ip" "foobar" {
	address = "${linode_instance.foobar.ip_address}"
}`, label)
}
