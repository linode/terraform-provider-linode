package linode

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLinodeStackscript_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_stackscript.stackscript"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLinodeStackScriptBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "deployments_active"),
					resource.TestCheckResourceAttrSet(resourceName, "deployments_total"),
					resource.TestCheckResourceAttrSet(resourceName, "username"),
					resource.TestCheckResourceAttrSet(resourceName, "created"),
					resource.TestCheckResourceAttrSet(resourceName, "updated"),
					resource.TestCheckResourceAttr(resourceName, "label", "my_stackscript"),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "is_public", "false"),
					resource.TestCheckResourceAttr(resourceName, "rev_note", "initial"),
					resource.TestCheckResourceAttr(resourceName, "script", testDataSourceLinodeStackScriptBasicScript),
					resource.TestCheckResourceAttr(resourceName, "images.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "images.0", "linode/ubuntu18.04"),
					resource.TestCheckResourceAttr(resourceName, "images.1", "linode/ubuntu16.04lts"),
					resource.TestCheckResourceAttr(resourceName, "user_defined_fields.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "user_defined_fields.0.name", "name"),
					resource.TestCheckResourceAttr(resourceName, "user_defined_fields.0.label", "Your name"),
					resource.TestCheckResourceAttr(resourceName, "user_defined_fields.0.default", "user"),
					resource.TestCheckResourceAttr(resourceName, "user_defined_fields.0.example", "Linus Torvalds"),
				),
			},
		},
	})
}

var testDataSourceLinodeStackScriptBasicScript = `#!/bin/bash
#<UDF name="name" label="Your name" example="Linus Torvalds" default="user">
# NAME=
echo "Hello, $NAME!"
`

func testDataSourceLinodeStackScriptBasic() string {
	return fmt.Sprintf(`
resource "linode_stackscript" "stackscript" {
	label = "my_stackscript"
	script = <<EOF
%sEOF
	images = ["linode/ubuntu18.04", "linode/ubuntu16.04lts"]
	description = "test"
	rev_note = "initial"
}

data "linode_stackscript" "stackscript" {
	id = linode_stackscript.stackscript.id
}`, testDataSourceLinodeStackScriptBasicScript)
}
