//go:build integration || sshkeys

package sshkeys_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/sshkeys/tmpl"
)

func TestAccDataSourceSSHKeys_basic(t *testing.T) {
	t.Parallel()

	testSSHKeyDataName := "data.linode_sshkeys.keys"

	keyLabel := acctest.RandomWithPrefix("tf_test")
	keySSH := acceptance.PublicKeyMaterial

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataFilterEmpty(t, keyLabel, keySSH),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testSSHKeyDataName, "sshkeys.#", "0"),
				),
			},
			{
				Config: tmpl.DataFilter(t, keyLabel, keySSH),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testSSHKeyDataName, "sshkeys.#", "1"),
					resource.TestCheckResourceAttr(testSSHKeyDataName, "sshkeys.0.label", keyLabel+"-0"),
					resource.TestCheckResourceAttr(testSSHKeyDataName, "sshkeys.0.ssh_key", keySSH),
					resource.TestCheckResourceAttrSet(testSSHKeyDataName, "sshkeys.0.id"),
					resource.TestCheckResourceAttrSet(testSSHKeyDataName, "sshkeys.0.created"),
				),
			},
			{
				Config: tmpl.DataBasic(t, keyLabel, keySSH),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testSSHKeyDataName, "sshkeys.#", "1"),
					resource.TestCheckResourceAttr(testSSHKeyDataName, "sshkeys.0.label", keyLabel+"-0"),
					resource.TestCheckResourceAttr(testSSHKeyDataName, "sshkeys.0.ssh_key", keySSH),
					resource.TestCheckResourceAttrSet(testSSHKeyDataName, "sshkeys.0.id"),
					resource.TestCheckResourceAttrSet(testSSHKeyDataName, "sshkeys.0.created"),
				),
			},
			{
				Config: tmpl.DataAll(t, keyLabel, keySSH),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckResourceAttrGreaterThan(testSSHKeyDataName, "sshkeys.#", 1),
					resource.TestCheckResourceAttrSet(testSSHKeyDataName, "sshkeys.0.label"),
					resource.TestCheckResourceAttrSet(testSSHKeyDataName, "sshkeys.0.ssh_key"),
					resource.TestCheckResourceAttrSet(testSSHKeyDataName, "sshkeys.0.id"),
					resource.TestCheckResourceAttrSet(testSSHKeyDataName, "sshkeys.0.created"),
				),
			},
		},
	})
}
