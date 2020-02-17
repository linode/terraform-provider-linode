package linode

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
					resource.TestCheckResourceAttr(resourceName, "label", "my_stackscript"),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "script", "#!/bin/bash\necho hello world\n"),
				),
			},
		},
	})
}

func testDataSourceLinodeStackScriptBasic() string {
	return `
resource "linode_stackscript" "stackscript" {
	label = "my_stackscript"
	script = <<EOF
#!/bin/bash
echo hello world
EOF
	images = ["linode/ubuntu18.04"]
	description = "test"
}

data "linode_stackscript" "stackscript" {
	id = linode_stackscript.stackscript.id
}`
}
