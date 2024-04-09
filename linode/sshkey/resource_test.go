//go:build integration || sshkey

package sshkey_test

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
	"github.com/linode/terraform-provider-linode/v2/linode/sshkey/tmpl"
)

func init() {
	resource.AddTestSweepers("linode_sshkey", &resource.Sweeper{
		Name: "linode_sshkey",
		F:    sweep,
	})
}

func sweep(prefix string) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	listOpts := acceptance.SweeperListOptions(prefix, "label")
	sshkeys, err := client.ListSSHKeys(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("Error getting sshkeys: %s", err)
	}
	for _, sshkey := range sshkeys {
		if !acceptance.ShouldSweep(prefix, sshkey.Label) {
			continue
		}
		err := client.DeleteSSHKey(context.Background(), sshkey.ID)
		if err != nil {
			return fmt.Errorf("Error destroying %s during sweep: %s", sshkey.Label, err)
		}
	}

	return nil
}

func TestAccResourceSSHKey_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_sshkey.foobar"
	sshkeyName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, sshkeyName, acceptance.PublicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					checkSSHKeyExists,
					resource.TestCheckResourceAttr(resName, "label", sshkeyName),
					resource.TestCheckResourceAttr(resName, "ssh_key", acceptance.PublicKeyMaterial),
					resource.TestCheckResourceAttrSet(resName, "created"),
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

func TestAccResourceSSHKey_update(t *testing.T) {
	t.Parallel()
	resName := "linode_sshkey.foobar"
	sshkeyName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, sshkeyName, acceptance.PublicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					checkSSHKeyExists,
					resource.TestCheckResourceAttr(resName, "label", sshkeyName),
					resource.TestCheckResourceAttr(resName, "ssh_key", acceptance.PublicKeyMaterial),
					resource.TestCheckResourceAttrSet(resName, "created"),
				),
			},
			{
				Config: tmpl.Updates(t, sshkeyName, acceptance.PublicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					checkSSHKeyExists,
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s_renamed", sshkeyName)),
					resource.TestCheckResourceAttr(resName, "ssh_key", acceptance.PublicKeyMaterial),
					resource.TestCheckResourceAttrSet(resName, "created"),
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

func checkSSHKeyExists(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_sshkey" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		_, err = client.GetSSHKey(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of SSHKey %s: %s", rs.Primary.Attributes["label"], err)
		}
	}

	return nil
}

func checkSSHKeyDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_sshkey" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}
		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)
		}

		_, err = client.GetSSHKey(context.Background(), id)

		if err == nil {
			return fmt.Errorf("Linode SSH Key with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error requesting Linode SSH Key with id %d", id)
		}
	}

	return nil
}
