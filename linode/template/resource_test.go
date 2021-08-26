// +build ignore

package template

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func init() {
	resource.AddTestSweepers("linode_template", &resource.Sweeper{
		Name: "linode_template",
		F:    sweep,
	})
}

func sweep(prefix string) error {
	client, err := acceptance.GetClientForSweepers()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	listOpts := acceptance.SweeperListOptions(prefix, "label")
	templates, err := client.ListTemplates(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("Error getting templates: %s", err)
	}
	for _, template := range templates {
		if !acceptance.ShouldSweep(prefix, template.Label) {
			continue
		}
		err := client.DeleteTemplate(context.Background(), template.ID)

		if err != nil {
			return fmt.Errorf("Error destroying %s during sweep: %s", template.Label, err)
		}
	}

	return nil
}

func TestAccResourceTemplate_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_template.foobar"
	var templateName = acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasic(templateName),
				Check: resource.ComposeTestCheckFunc(
					checkTemplateExists,
					resource.TestCheckResourceAttr(resName, "label", templateName),
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

func TestAccResourceTemplate_update(t *testing.T) {
	t.Parallel()

	var templateName = acctest.RandomWithPrefix("tf_test")
	resName := "linode_template.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasic(templateName),
				Check: resource.ComposeTestCheckFunc(
					checkTemplateExists,
					resource.TestCheckResourceAttr(resName, "label", templateName),
				),
			},
			{
				Config: resourceConfigUpdates(templateName),
				Check: resource.ComposeTestCheckFunc(
					checkTemplateExists,
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s_renamed", templateName)),
				),
			},
		},
	})
}

func checkTemplateExists(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_template" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		_, err = client.GetTemplate(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of Template %s: %s", rs.Primary.Attributes["label"], err)
		}
	}

	return nil
}

func checkTemplateDestroy(s *terraform.State) error {
	client, ok := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
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

		_, err = client.GetTemplate(context.Background(), id)

		if err == nil {
			return fmt.Errorf("Linode Template with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error requesting Linode Template with id %d", id)
		}
	}

	return nil
}

func resourceConfigBasic(template string) string {
	return fmt.Sprintf(`
resource "linode_template" "foobar" {
	label = "%s"
}`, template)
}

func resourceConfigUpdates(template string) string {
	return fmt.Sprintf(`
resource "linode_template" "foobar" {
	label = "%s_renamed"
}`, template)
}
