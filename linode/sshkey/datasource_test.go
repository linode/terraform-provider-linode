package sshkey_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/sshkey/tmpl"
)

func TestAccDataSourceSSHKey_basic(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf_test")
	// resourceName := "data.linode_sshkey.foobar"

	// TODO(ellisbenjamin) -- This test passes only because of the Destroy: true statement and needs attention.
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config:  tmpl.Basic(t, label, acceptance.PublicKeyMaterial),
				Destroy: true,
			},
			// {
			// 	Config: resourceConfigBasic(label, acceptance.PublicKeyMaterial) + testDataSourceLinodeSSHKey(label, acceptance.PublicKeyMaterial),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestCheckResourceAttr(resourceName, "ssh_key", acceptance.PublicKeyMaterial),
			// 		resource.TestCheckResourceAttr(resourceName, "label", label),
			// 	),
			// },
			{
				Config:      tmpl.DataBasic(t, label),
				ExpectError: regexp.MustCompile(label + " was not found"),
			},
		},
	})
}
