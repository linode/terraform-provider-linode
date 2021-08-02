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

func TestAccLinodeInstanceIP_noboot(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: accTestWithProvider(
					testAccCheckLinodeInstanceIPInstanceNoBoot(name),
					map[string]interface{}{
						providerKeySkipInstanceReadyPoll: true,
					}),
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
        image = "linode/alpine3.14"
}`, label, publicKeyMaterial)
}

func testAccCheckLinodeInstanceIPBasic(label string) string {
	return testAccCheckLinodeInstanceIPInstance(label) + fmt.Sprintf(`
resource "linode_instance_ip" "test" {
	linode_id = linode_instance.%s.id
	public = true
}`, label)
}

func testAccCheckLinodeInstanceIPInstanceNoBoot(label string) string {
	return fmt.Sprintf(`
resource "linode_instance" "%[1]s" {
	label = "%[1]s"
	group = "tf_test"
	type = "g6-nanode-1"
	region = "us-east"
}

resource "linode_instance_ip" "test" {
	linode_id = linode_instance.%[1]s.id
	public = true
}`, label, publicKeyMaterial)
}
