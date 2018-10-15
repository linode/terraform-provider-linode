package linode

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestDataSourceLinodeSSHKey(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf_test")
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("linode@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	resourceName := "data.linode_sshkey.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLinodeSSHKey(label, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ssh_key", publicKeyMaterial),
					resource.TestCheckResourceAttr(resourceName, "label", label),
				),
			},
		},
	})
}

func testDataSourceLinodeSSHKey(label, sshKey string) string {
	return fmt.Sprintf(`
data "linode_sshkey" "foobar" {
	label = "%s"
}

resource "linode_sshkey" "foobar" {
	label = "%s"
	ssh_key = "%s"
}
`, label, label, sshKey)
}
