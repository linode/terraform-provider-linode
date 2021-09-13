package instanceip_test

import (
	"fmt"
	"testing"

	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testInstanceIPResName = "linode_instance_ip.test"

func TestAccInstanceIP_basic(t *testing.T) {
	t.Parallel()

	var instance linodego.Instance

	name := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasic(name, true),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists("linode_instance.foobar", &instance),
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
					acceptance.AssertInstanceReboot(t, true, &instance)
				},
				Config: resourceConfigBasic(name, true),
			},
		},
	})
}

func TestAccInstanceIP_noboot(t *testing.T) {
	t.Parallel()

	var instance linodego.Instance

	name := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigNoBoot(name, true),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists("linode_instance.foobar", &instance),
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
				Config: resourceConfigNoBoot(name, true),
				PreConfig: func() {
					acceptance.AssertInstanceReboot(t, false, &instance)
				},
			},
		},
	})
}

func TestAccInstanceIP_noApply(t *testing.T) {
	t.Parallel()

	var instance linodego.Instance

	name := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasic(name, false),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists("linode_instance.foobar", &instance),
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
					acceptance.AssertInstanceReboot(t, false, &instance)
				},
				Config: resourceConfigBasic(name, false),
			},
		},
	})
}

func instanceConfigBasic(label string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	group = "tf_test"
	type = "g6-nanode-1"
	region = "us-east"
	image = "linode/alpine3.14"
}`, label)
}

func resourceConfigBasic(label string, applyImmediately bool) string {
	return instanceConfigBasic(label) + fmt.Sprintf(`
resource "linode_instance_ip" "test" {
	linode_id = linode_instance.foobar.id
	public = true
	apply_immediately = %t
}`, applyImmediately)
}

func resourceConfigNoBoot(label string, applyImmediately bool) string {
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
