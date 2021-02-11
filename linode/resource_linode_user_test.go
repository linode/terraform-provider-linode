package linode

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testUserResName = "linode_user.test"

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
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
