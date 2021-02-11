package linode

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
)

const testUserResName = "linode_user.test"

func testAccCheckLinodeUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_user" {
			continue
		}

		username := rs.Primary.ID
		_, err := client.GetUser(context.TODO(), username)

		if err == nil {
			return fmt.Errorf("should not find user %s existing after delete", username)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error getting user %s: %s", username, err)
		}
	}
	return nil
}

func testAccCheckLinodeUserConfigBasic(username, email string, restricted bool) string {
	return fmt.Sprintf(`
resource "linode_user" "test" {
	username = "%s"
	email = "%s"
	restricted = %t
}`, username, email, restricted)
}

func TestAccLinodeUser_basic(t *testing.T) {
	t.Parallel()

	username := acctest.RandomWithPrefix("tf-test")
	email := username + "@example.com"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeUserConfigBasic(username, email, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testUserResName, "email", email),
					resource.TestCheckResourceAttr(testUserResName, "username", username),
					resource.TestCheckResourceAttr(testUserResName, "restricted", "true"),
					resource.TestCheckResourceAttr(testUserResName, "ssh_keys.#", "0"),
					resource.TestCheckResourceAttr(testUserResName, "tfa_enabled", "false"),
				),
			},
		},
	})
}

func TestAccLinodeUser_updates(t *testing.T) {
	t.Parallel()

	username := acctest.RandomWithPrefix("tf-test")
	updatedUsername := acctest.RandomWithPrefix("tf-test")
	email := username + "@example.com"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeUserConfigBasic(username, email, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testUserResName, "email", email),
					resource.TestCheckResourceAttr(testUserResName, "username", username),
					resource.TestCheckResourceAttr(testUserResName, "restricted", "false"),
					resource.TestCheckResourceAttr(testUserResName, "ssh_keys.#", "0"),
					resource.TestCheckResourceAttr(testUserResName, "tfa_enabled", "false"),
				),
			},
			{
				Config: testAccCheckLinodeUserConfigBasic(updatedUsername, email, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testUserResName, "email", email),
					resource.TestCheckResourceAttr(testUserResName, "username", updatedUsername),
					resource.TestCheckResourceAttr(testUserResName, "restricted", "true"),
					resource.TestCheckResourceAttr(testUserResName, "ssh_keys.#", "0"),
					resource.TestCheckResourceAttr(testUserResName, "tfa_enabled", "false"),
				),
			},
		},
	})
}
