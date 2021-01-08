package linode

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLinodeSSHKey_basic(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf_test")
	// resourceName := "data.linode_sshkey.foobar"

	// TODO(ellisbenjamin) -- This test passes only because of the Destroy: true statement and needs attention.
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config:  testAccCheckLinodeSSHKeyConfigBasic(label, publicKeyMaterial),
				Destroy: true,
			},
			// {
			// 	Config: testAccCheckLinodeSSHKeyConfigBasic(label, publicKeyMaterial) + testDataSourceLinodeSSHKey(label, publicKeyMaterial),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestCheckResourceAttr(resourceName, "ssh_key", publicKeyMaterial),
			// 		resource.TestCheckResourceAttr(resourceName, "label", label),
			// 	),
			// },
			{
				Config:      testDataSourceLinodeSSHKeyBasic(label, publicKeyMaterial),
				ExpectError: regexp.MustCompile(label + " was not found"),
			},
		},
	})
}

func testDataSourceLinodeSSHKeyBasic(label, sshKey string) string {
	return fmt.Sprintf(`
data "linode_sshkey" "foobar" {
	label = "%s"
}

`, label)
}
