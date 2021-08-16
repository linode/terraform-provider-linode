package linode

import (
	"fmt"
	"testing"

	"github.com/linode/linodego"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testInstanceIPResName = "linode_instance_ip.test"

func TestAccLinodeInstanceIP_basic(t *testing.T) {
	t.Parallel()

	var instance linodego.Instance

	name := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeInstanceIPBasic(name, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeInstanceExists("linode_instance.foobar", &instance),
					resource.TestCheckResourceAttrSet(testInstanceIPResName, "address"),
					resource.TestCheckResourceAttrSet(testInstanceIPResName, "gateway"),
					resource.TestCheckResourceAttrSet(testInstanceIPResName, "prefix"),
					resource.TestCheckResourceAttrSet(testInstanceIPResName, "rdns"),
					resource.TestCheckResourceAttrSet(testInstanceIPResName, "subnet_mask"),
					resource.TestCheckResourceAttr(testInstanceIPResName, "region", "us-east"),
					resource.TestCheckResourceAttr(testInstanceIPResName, "type", "ipv4"),
				),
			},
			{
				PreConfig: func() {
					testAccAssertReboot(t, true, &instance)
				},
				Config: testAccCheckLinodeInstanceIPBasic(name, true),
			},
		},
	})
}

func TestAccLinodeInstanceIP_noboot(t *testing.T) {
	t.Parallel()

	var instance linodego.Instance

	name := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeInstanceIPInstanceNoBoot(name, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeInstanceExists("linode_instance.foobar", &instance),
					resource.TestCheckResourceAttrSet(testInstanceIPResName, "address"),
					resource.TestCheckResourceAttrSet(testInstanceIPResName, "gateway"),
					resource.TestCheckResourceAttrSet(testInstanceIPResName, "prefix"),
					resource.TestCheckResourceAttrSet(testInstanceIPResName, "rdns"),
					resource.TestCheckResourceAttrSet(testInstanceIPResName, "subnet_mask"),
					resource.TestCheckResourceAttr(testInstanceIPResName, "region", "us-east"),
					resource.TestCheckResourceAttr(testInstanceIPResName, "type", "ipv4"),
				),
			},
			{
				Config: testAccCheckLinodeInstanceIPInstanceNoBoot(name, true),
				PreConfig: func() {
					testAccAssertReboot(t, false, &instance)
				},
			},
		},
	})
}

func TestAccLinodeInstanceIP_noApply(t *testing.T) {
	t.Parallel()

	var instance linodego.Instance

	name := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeInstanceIPBasic(name, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeInstanceExists("linode_instance.foobar", &instance),
					resource.TestCheckResourceAttrSet(testInstanceIPResName, "address"),
					resource.TestCheckResourceAttrSet(testInstanceIPResName, "gateway"),
					resource.TestCheckResourceAttrSet(testInstanceIPResName, "prefix"),
					resource.TestCheckResourceAttrSet(testInstanceIPResName, "rdns"),
					resource.TestCheckResourceAttrSet(testInstanceIPResName, "subnet_mask"),
					resource.TestCheckResourceAttr(testInstanceIPResName, "region", "us-east"),
					resource.TestCheckResourceAttr(testInstanceIPResName, "type", "ipv4"),
				),
			},
			{
				PreConfig: func() {
					testAccAssertReboot(t, false, &instance)
				},
				Config: testAccCheckLinodeInstanceIPBasic(name, false),
			},
		},
	})
}

func testAccCheckLinodeInstanceIPInstance(label string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	group = "tf_test"
	type = "g6-nanode-1"
	region = "us-east"
	image = "linode/alpine3.14"
}`, label)
}

func testAccCheckLinodeInstanceIPBasic(label string, applyImmediately bool) string {
	return testAccCheckLinodeInstanceIPInstance(label) + fmt.Sprintf(`
resource "linode_instance_ip" "test" {
	linode_id = linode_instance.foobar.id
	public = true
	apply_immediately = %t
}`, applyImmediately)
}

func testAccCheckLinodeInstanceIPInstanceNoBoot(label string, applyImmediately bool) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	group = "tf_test"
	type = "g6-nanode-1"
	region = "us-east"
}

resource "linode_instance_ip" "test" {
	linode_id = linode_instance.foobar.id
	public = true
	apply_immediately = %t
}
`, label, applyImmediately)
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
