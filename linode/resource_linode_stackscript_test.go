package linode

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
)

func init() {
	resource.AddTestSweepers("linode_stackscript", &resource.Sweeper{
		Name: "linode_stackscript",
		F:    testSweepLinodeStackScript,
	})
}

func testSweepLinodeStackScript(prefix string) error {
	client, err := getClientForSweepers()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	listOpts := sweeperListOptions(prefix, "label")
	stackscripts, err := client.ListStackscripts(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("Error getting stackscripts: %s", err)
	}
	for _, stackscript := range stackscripts {
		if !shouldSweepAcceptanceTestResource(prefix, stackscript.Label) {
			continue
		}
		err := client.DeleteStackscript(context.Background(), stackscript.ID)

		if err != nil {
			return fmt.Errorf("Error destroying %s during sweep: %s", stackscript.Label, err)
		}
	}

	return nil
}

func TestAccLinodeStackscript_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_stackscript.foobar"
	var stackscriptName = acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeStackscriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeStackscriptBasic(stackscriptName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeStackscriptExists,
					resource.TestCheckResourceAttr(resName, "description", "tf_test stackscript"),
					resource.TestCheckResourceAttr(resName, "rev_note", "initial"),
					resource.TestCheckResourceAttr(resName, "images.0", "linode/ubuntu18.04"),
					resource.TestCheckResourceAttr(resName, "label", stackscriptName),
				),
			},

			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLinodeStackscript_update(t *testing.T) {
	t.Parallel()

	var stackscriptName = acctest.RandomWithPrefix("tf_test")
	var resName = "linode_stackscript.foobar"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeStackscriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeStackscriptBasic(stackscriptName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeStackscriptExists,
					resource.TestCheckResourceAttr(resName, "description", "tf_test stackscript"),
					resource.TestCheckResourceAttr(resName, "rev_note", "initial"),
					resource.TestCheckResourceAttr(resName, "images.0", "linode/ubuntu18.04"),
					resource.TestCheckResourceAttr(resName, "label", stackscriptName),
				),
			},
			{
				Config: testAccCheckLinodeStackscriptBasicRenamed(stackscriptName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeStackscriptExists,
					resource.TestCheckResourceAttr(resName, "description", "tf_test stackscript"),
					resource.TestCheckResourceAttr(resName, "rev_note", "initial"),
					resource.TestCheckResourceAttr(resName, "images.0", "linode/ubuntu18.04"),
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s_renamed", stackscriptName)),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLinodeStackscript_codeChange(t *testing.T) {
	t.Parallel()

	var stackscriptName = acctest.RandomWithPrefix("tf_test")
	var resName = "linode_stackscript.foobar"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeStackscriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeStackscriptBasic(stackscriptName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeStackscriptExists,
					resource.TestCheckResourceAttr(resName, "description", "tf_test stackscript"),
					resource.TestCheckResourceAttr(resName, "rev_note", "initial"),
					resource.TestCheckResourceAttr(resName, "images.0", "linode/ubuntu18.04"),
					resource.TestCheckResourceAttr(resName, "script", "#!/bin/bash\necho hello\n"),
					resource.TestCheckResourceAttr(resName, "user_defined_fields.#", "0"),
					resource.TestCheckResourceAttr(resName, "label", stackscriptName),
				),
			},
			{
				Config: testAccCheckLinodeStackscriptCodeChange(stackscriptName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeStackscriptExists,
					resource.TestCheckResourceAttr(resName, "description", "tf_test stackscript"),
					resource.TestCheckResourceAttr(resName, "rev_note", "second"),
					resource.TestCheckResourceAttr(resName, "script", "#!/bin/bash\n# <UDF name=\"hasudf\" label=\"a label\" example=\"an example\" default=\"a default\">\necho bye\n"),
					resource.TestCheckResourceAttr(resName, "images.0", "linode/ubuntu18.04"),
					resource.TestCheckResourceAttr(resName, "images.1", "linode/ubuntu16.04lts"),
					resource.TestCheckResourceAttr(resName, "user_defined_fields.#", "1"),
					resource.TestCheckResourceAttr(resName, "user_defined_fields.0.name", "hasudf"),
					resource.TestCheckResourceAttr(resName, "user_defined_fields.0.label", "a label"),
					resource.TestCheckResourceAttr(resName, "user_defined_fields.0.default", "a default"),
					resource.TestCheckResourceAttr(resName, "user_defined_fields.0.example", "an example"),
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s", stackscriptName)),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLinodeStackscriptExists(s *terraform.State) error {
	client := testAccProvider.Meta().(linodego.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_stackscript" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		_, err = client.GetStackscript(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of Stackscript %s: %s", rs.Primary.Attributes["label"], err)
		}
	}

	return nil
}

func testAccCheckLinodeStackscriptDestroy(s *terraform.State) error {
	client, ok := testAccProvider.Meta().(linodego.Client)
	if !ok {
		return fmt.Errorf("Error getting Linode client")
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_stackscript" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}
		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)

		}

		_, err = client.GetStackscript(context.Background(), id)

		if err == nil {
			return fmt.Errorf("Linode Stackscript with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error requesting Linode Stackscript with id %d", id)
		}
	}

	return nil
}

func testAccCheckLinodeStackscriptBasic(stackscript string) string {
	return fmt.Sprintf(`
resource "linode_stackscript" "foobar" {
	label = "%s"
	script = <<EOF
#!/bin/bash
echo hello
EOF
	images = ["linode/ubuntu18.04"]
	description = "tf_test stackscript"
	rev_note = "initial"
}`, stackscript)
}

func testAccCheckLinodeStackscriptBasicRenamed(stackscript string) string {
	return fmt.Sprintf(`
resource "linode_stackscript" "foobar" {
	label = "%s_renamed"
	script = <<EOF
#!/bin/bash
echo hello
EOF
	images = ["linode/ubuntu18.04"]
	description = "tf_test stackscript"
	rev_note = "initial"
}`, stackscript)
}

func testAccCheckLinodeStackscriptCodeChange(stackscript string) string {
	return fmt.Sprintf(`
resource "linode_stackscript" "foobar" {
	label = "%s"
	script = <<EOF
#!/bin/bash
# <UDF name="hasudf" label="a label" example="an example" default="a default">
echo bye
EOF
	images = ["linode/ubuntu18.04", "linode/ubuntu16.04lts"]
	description = "tf_test stackscript"
	rev_note = "second"
}`, stackscript)
}
