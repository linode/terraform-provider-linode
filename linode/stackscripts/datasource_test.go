//go:build integration || stackscripts

package stackscripts_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/stackscripts/tmpl"
)

var basicStackScript = `#!/bin/bash
#<UDF name="name" label="Your name" example="Linus Torvalds" default="user">
# NAME=
echo "Hello, $NAME!"
`

func TestAccDataSourceStackscripts_basic_smoke(t *testing.T) {
	t.Parallel()

	stackScriptName := acctest.RandomWithPrefix("tf_test")

	resourceName := "data.linode_stackscripts.stackscript"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, stackScriptName, basicStackScript),
				Check: resource.ComposeTestCheckFunc(
					validateStackscript(resourceName, stackScriptName),
				),
			},
			{
				Config: tmpl.DataSubString(t, stackScriptName, basicStackScript),
				Check: resource.ComposeTestCheckFunc(
					validateStackscript(resourceName, stackScriptName),
				),
			},
			{
				Config: tmpl.DataLatest(t, stackScriptName, basicStackScript),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "stackscripts.#", "1"),
					validateStackscript(resourceName, stackScriptName),
				),
			},
			{
				Config: tmpl.DataClientFilter(t, stackScriptName, basicStackScript),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "stackscripts.#", "1"),
					validateStackscript(resourceName, stackScriptName),
				),
			},
		},
	})
}

func validateStackscript(resourceName, stackScriptName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "stackscripts.0.id"),
		resource.TestCheckResourceAttrSet(resourceName, "stackscripts.0.deployments_active"),
		resource.TestCheckResourceAttrSet(resourceName, "stackscripts.0.deployments_total"),
		resource.TestCheckResourceAttrSet(resourceName, "stackscripts.0.username"),
		resource.TestCheckResourceAttrSet(resourceName, "stackscripts.0.created"),
		resource.TestCheckResourceAttrSet(resourceName, "stackscripts.0.updated"),
		resource.TestCheckResourceAttr(resourceName, "stackscripts.0.label", stackScriptName),
		resource.TestCheckResourceAttr(resourceName, "stackscripts.0.description", "test"),
		resource.TestCheckResourceAttr(resourceName, "stackscripts.0.is_public", "false"),
		resource.TestCheckResourceAttr(resourceName, "stackscripts.0.rev_note", "initial"),
		resource.TestCheckResourceAttr(resourceName, "stackscripts.0.script", basicStackScript),
		resource.TestCheckResourceAttr(resourceName, "stackscripts.0.images.#", "2"),
		acceptance.CheckListContains(resourceName, "stackscripts.0.images", "linode/ubuntu18.04"),
		acceptance.CheckListContains(resourceName, "stackscripts.0.images", "linode/ubuntu16.04lts"),
		resource.TestCheckResourceAttr(resourceName, "stackscripts.0.user_defined_fields.#", "1"),
		resource.TestCheckResourceAttr(resourceName, "stackscripts.0.user_defined_fields.0.name", "name"),
		resource.TestCheckResourceAttr(resourceName, "stackscripts.0.user_defined_fields.0.label", "Your name"),
		resource.TestCheckResourceAttr(resourceName, "stackscripts.0.user_defined_fields.0.default", "user"),
		resource.TestCheckResourceAttr(resourceName, "stackscripts.0.user_defined_fields.0.example", "Linus Torvalds"),
	)
}
