// +build ignore

package linode

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/chiefy/linodego"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccLinodeTemplateBasic(t *testing.T) {
	t.Parallel()

	resName := "linode_template.foobar"
	var templateName = acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeTemplateConfigBasic(templateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeTemplateExists,
					resource.TestCheckResourceAttr(resName, "label", templateName),
				),
			},

			resource.TestStep{
				ResourceName: resName,
				ImportState:  true,
			},
		},
	})
}

func TestAccLinodeTemplateUpdate(t *testing.T) {
	t.Parallel()

	var templateName = acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeTemplateConfigBasic(templateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeTemplateExists,
					resource.TestCheckResourceAttr("linode_template.foobar", "label", templateName),
				),
			},
			resource.TestStep{
				Config: testAccCheckLinodeTemplateConfigUpdates(templateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeTemplateExists,
					resource.TestCheckResourceAttr("linode_template.foobar", "label", fmt.Sprintf("%s_renamed", templateName)),
				),
			},
		},
	})
}

func testAccCheckLinodeTemplateExists(s *terraform.State) error {
	client := testAccProvider.Meta().(linodego.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_template" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)

		_, err = client.GetTemplate(id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of Template %s: %s", rs.Primary.Attributes["label"], err)
		}
	}

	return nil
}

func testAccCheckLinodeTemplateDestroy(s *terraform.State) error {
	client, ok := testAccProvider.Meta().(linodego.Client)
	if !ok {
		return fmt.Errorf("Error getting Linode client")
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_template" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}
		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)

		}

		_, err = client.GetTemplate(id)

		if err == nil {
			return fmt.Errorf("Linode Template with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error requesting Linode Template with id %d", id)
		}
	}

	return nil
}

func testAccCheckLinodeTemplateConfigBasic(template string) string {
	return fmt.Sprintf(`
resource "linode_template" "foobar" {
	label = "%s"
}`, template)
}

func testAccCheckLinodeTemplateConfigUpdates(template string) string {
	return fmt.Sprintf(`
resource "linode_template" "foobar" {
	label = "%s_renamed"
}`, template)
}
