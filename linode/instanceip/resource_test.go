package instanceip_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/instanceip/tmpl"
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
				Config: tmpl.Basic(t, name, true),
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
				Config: tmpl.Basic(t, name, true),
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
				Config: tmpl.NoBoot(t, name, true),
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
				Config: tmpl.NoBoot(t, name, true),
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
				Config: tmpl.Basic(t, name, false),
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
				Config: tmpl.Basic(t, name, false),
			},
		},
	})
}
