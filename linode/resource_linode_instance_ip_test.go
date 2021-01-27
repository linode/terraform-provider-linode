package linode

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testInstanceIPResName = "linode_instance_ip.test"

func TestAccLinodeInstanceIP_basic(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeInstanceIPBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testInstanceIPResName, "address"),
					resource.TestCheckResourceAttrSet(testInstanceIPResName, "gateway"),
					resource.TestCheckResourceAttrSet(testInstanceIPResName, "prefix"),
					resource.TestCheckResourceAttrSet(testInstanceIPResName, "rdns"),
					resource.TestCheckResourceAttrSet(testInstanceIPResName, "subnet_mask"),
					resource.TestCheckResourceAttr(testInstanceIPResName, "region", "us-east"),
					resource.TestCheckResourceAttr(testInstanceIPResName, "type", "ipv4"),
				),
			},
		},
	})
}

func testAccCheckLinodeInstanceIPInstance(label string) string {
	return fmt.Sprintf(`
resource "linode_instance" "%[1]s" {
	label = "%[1]s"
	group = "tf_test"
	type = "g6-nanode-1"
	region = "us-east"
	disk {
		label = "disk"
		image = "linode/alpine3.11"
		root_pass = "b4d_p4s5"
		authorized_keys = ["%[2]s"]
		size = 3000
	}
}`, label, publicKeyMaterial)
}

func testAccCheckLinodeInstanceIPBasic(label string) string {
	return testAccCheckLinodeInstanceIPInstance(label) + fmt.Sprintf(`
resource "linode_instance_ip" "test" {
	linode_id = linode_instance.%s.id
	public = true
}`, label)
}
