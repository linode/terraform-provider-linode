package sshkey_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

func TestAccDataSourceSSHKey_basic(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf_test")
	// resourceName := "data.linode_sshkey.foobar"

	// TODO(ellisbenjamin) -- This test passes only because of the Destroy: true statement and needs attention.
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config:  resourceConfigBasic(label, acceptance.PublicKeyMaterial),
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
				Config:      dataSourceConfigBasic(label, acceptance.PublicKeyMaterial),
				ExpectError: regexp.MustCompile(label + " was not found"),
			},
		},
	})
}

func dataSourceConfigBasic(label, sshKey string) string {
	return fmt.Sprintf(`
data "linode_sshkey" "foobar" {
	label = "%s"
}

`, label)
}
