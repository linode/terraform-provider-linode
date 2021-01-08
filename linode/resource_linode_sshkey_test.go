package linode

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
)

func init() {
	resource.AddTestSweepers("linode_sshkey", &resource.Sweeper{
		Name: "linode_sshkey",
		F:    testSweepLinodeSSHKey,
	})
}

func testSweepLinodeSSHKey(prefix string) error {
	client, err := getClientForSweepers()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	listOpts := sweeperListOptions(prefix, "label")
	sshkeys, err := client.ListSSHKeys(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("Error getting sshkeys: %s", err)
	}
	for _, sshkey := range sshkeys {
		if !shouldSweepAcceptanceTestResource(prefix, sshkey.Label) {
			continue
		}
		err := client.DeleteSSHKey(context.Background(), sshkey.ID)

		if err != nil {
			return fmt.Errorf("Error destroying %s during sweep: %s", sshkey.Label, err)
		}
	}

	return nil
}

func TestAccLinodeSSHKey_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_sshkey.foobar"
	sshkeyName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeSSHKeyConfigBasic(sshkeyName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeSSHKeyExists,
					resource.TestCheckResourceAttr(resName, "label", sshkeyName),
					resource.TestCheckResourceAttr(resName, "ssh_key", publicKeyMaterial),
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

func TestAccLinodeSSHKey_update(t *testing.T) {
	t.Parallel()
	resName := "linode_sshkey.foobar"
	sshkeyName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeSSHKeyConfigBasic(sshkeyName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeSSHKeyExists,
					resource.TestCheckResourceAttr(resName, "label", sshkeyName),
					resource.TestCheckResourceAttr(resName, "ssh_key", publicKeyMaterial),
					resource.TestCheckResourceAttrSet(resName, "created"),
				),
			},
			{
				Config: testAccCheckLinodeSSHKeyConfigUpdates(sshkeyName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeSSHKeyExists,
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s_renamed", sshkeyName)),
					resource.TestCheckResourceAttr(resName, "ssh_key", publicKeyMaterial),
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

func testAccCheckLinodeSSHKeyExists(s *terraform.State) error {
	client := testAccProvider.Meta().(linodego.Client)

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

func testAccCheckLinodeSSHKeyDestroy(s *terraform.State) error {
	client, ok := testAccProvider.Meta().(linodego.Client)
	if !ok {
		return fmt.Errorf("Error getting Linode client")
	}
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

func testAccCheckLinodeSSHKeyConfigBasic(label, sshkey string) string {
	return fmt.Sprintf(`
resource "linode_sshkey" "foobar" {
	label = "%s"
	ssh_key = "%s"
}`, label, sshkey)
}

func testAccCheckLinodeSSHKeyConfigUpdates(label, sshkey string) string {
	return fmt.Sprintf(`
resource "linode_sshkey" "foobar" {
	label = "%s_renamed"
	ssh_key = "%s"
}`, label, sshkey)
}
