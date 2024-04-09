//go:build integration || stackscript

package stackscript_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/stackscript/tmpl"
)

func init() {
	resource.AddTestSweepers("linode_stackscript", &resource.Sweeper{
		Name: "linode_stackscript",
		F:    sweep,
	})
}

func sweep(prefix string) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	listOpts := acceptance.SweeperListOptions(prefix, "label")
	stackscripts, err := client.ListStackscripts(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("Error getting stackscripts: %s", err)
	}
	for _, stackscript := range stackscripts {
		if !acceptance.ShouldSweep(prefix, stackscript.Label) {
			continue
		}
		err := client.DeleteStackscript(context.Background(), stackscript.ID)
		if err != nil {
			return fmt.Errorf("Error destroying %s during sweep: %s", stackscript.Label, err)
		}
	}

	return nil
}

func TestAccResourceStackscript_basic_smoke(t *testing.T) {
	t.Parallel()

	resName := "linode_stackscript.foobar"
	stackscriptName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkStackscriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, stackscriptName),
				Check: resource.ComposeTestCheckFunc(
					checkStackscriptExists,
					resource.TestCheckResourceAttr(resName, "description", "tf_test stackscript"),
					resource.TestCheckResourceAttr(resName, "rev_note", "initial"),
					acceptance.CheckListContains(resName, "images", "linode/ubuntu18.04"),
					resource.TestCheckResourceAttr(resName, "label", stackscriptName),
				),
			},

			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"created", "updated"}, // Ignore strict comparison for these attributes
			},
		},
	})
}

func TestAccResourceStackscript_update(t *testing.T) {
	t.Parallel()

	stackscriptName := acctest.RandomWithPrefix("tf_test")
	resName := "linode_stackscript.foobar"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkStackscriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, stackscriptName),
				Check: resource.ComposeTestCheckFunc(
					checkStackscriptExists,
					resource.TestCheckResourceAttr(resName, "description", "tf_test stackscript"),
					resource.TestCheckResourceAttr(resName, "rev_note", "initial"),
					acceptance.CheckListContains(resName, "images", "linode/ubuntu18.04"),
					resource.TestCheckResourceAttr(resName, "label", stackscriptName),
				),
			},
			{
				Config: tmpl.Basic(t, stackscriptName+"_renamed"),
				Check: resource.ComposeTestCheckFunc(
					checkStackscriptExists,
					resource.TestCheckResourceAttr(resName, "description", "tf_test stackscript"),
					resource.TestCheckResourceAttr(resName, "rev_note", "initial"),
					acceptance.CheckListContains(resName, "images", "linode/ubuntu18.04"),
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

func TestAccResourceStackscript_codeChange(t *testing.T) {
	t.Parallel()

	stackscriptName := acctest.RandomWithPrefix("tf_test")
	resName := "linode_stackscript.foobar"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkStackscriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, stackscriptName),
				Check: resource.ComposeTestCheckFunc(
					checkStackscriptExists,
					resource.TestCheckResourceAttr(resName, "description", "tf_test stackscript"),
					resource.TestCheckResourceAttr(resName, "rev_note", "initial"),
					acceptance.CheckListContains(resName, "images", "linode/ubuntu18.04"),
					resource.TestCheckResourceAttr(resName, "script", "#!/bin/bash\necho hello\n"),
					resource.TestCheckResourceAttr(resName, "user_defined_fields.#", "0"),
					resource.TestCheckResourceAttr(resName, "label", stackscriptName),
				),
			},
			{
				Config: tmpl.CodeChange(t, stackscriptName),
				Check: resource.ComposeTestCheckFunc(
					checkStackscriptExists,
					resource.TestCheckResourceAttr(resName, "description", "tf_test stackscript"),
					resource.TestCheckResourceAttr(resName, "rev_note", "second"),
					resource.TestCheckResourceAttr(resName, "script", "#!/bin/bash\n# <UDF name=\"hasudf\" label=\"a label\" example=\"an example\" default=\"a default\">\necho bye\n"),
					acceptance.CheckListContains(resName, "images", "linode/ubuntu18.04"),
					acceptance.CheckListContains(resName, "images", "linode/ubuntu16.04lts"),
					resource.TestCheckResourceAttr(resName, "user_defined_fields.#", "1"),
					resource.TestCheckResourceAttr(resName, "user_defined_fields.0.name", "hasudf"),
					resource.TestCheckResourceAttr(resName, "user_defined_fields.0.label", "a label"),
					resource.TestCheckResourceAttr(resName, "user_defined_fields.0.default", "a default"),
					resource.TestCheckResourceAttr(resName, "user_defined_fields.0.example", "an example"),
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

func checkStackscriptExists(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

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

func checkStackscriptDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
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
			return fmt.Errorf("error requesting Linode Stackscript with id %d: %s", id, apiErr)
		}
	}

	return nil
}
