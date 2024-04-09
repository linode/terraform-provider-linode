//go:build integration || user

package user_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/user/tmpl"
)

const testUserResName = "linode_user.test"

func TestAccResourceUser_basic(t *testing.T) {
	t.Parallel()

	username := acctest.RandomWithPrefix("tf-test")
	email := username + "@example.com"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, username, email, true),
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

func TestAccResourceUser_updates(t *testing.T) {
	t.Parallel()

	username := acctest.RandomWithPrefix("tf-test")
	updatedUsername := acctest.RandomWithPrefix("tf-test")
	email := username + "@example.com"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, username, email, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testUserResName, "email", email),
					resource.TestCheckResourceAttr(testUserResName, "username", username),
					resource.TestCheckResourceAttr(testUserResName, "restricted", "false"),
					resource.TestCheckResourceAttr(testUserResName, "ssh_keys.#", "0"),
					resource.TestCheckResourceAttr(testUserResName, "tfa_enabled", "false"),
				),
			},
			{
				Config: tmpl.Basic(t, updatedUsername, email, true),
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

func TestAccResourceUser_grants(t *testing.T) {
	t.Parallel()

	username := acctest.RandomWithPrefix("tf-test")
	instance := acctest.RandomWithPrefix("tf-test")

	email := username + "@example.com"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Grants(t, username, email),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testUserResName, "global_grants.0.account_access", ""),
					resource.TestCheckResourceAttr(testUserResName, "global_grants.0.add_domains", "true"),
					resource.TestCheckResourceAttr(testUserResName, "global_grants.0.add_databases", "false"),
					resource.TestCheckResourceAttr(testUserResName, "global_grants.0.add_firewalls", "true"),
					resource.TestCheckResourceAttr(testUserResName, "global_grants.0.add_images", "false"),
					resource.TestCheckResourceAttr(testUserResName, "global_grants.0.add_linodes", "true"),
					resource.TestCheckResourceAttr(testUserResName, "global_grants.0.add_longview", "false"),
					resource.TestCheckResourceAttr(testUserResName, "global_grants.0.add_nodebalancers", "true"),
					resource.TestCheckResourceAttr(testUserResName, "global_grants.0.add_stackscripts", "false"),
					resource.TestCheckResourceAttr(testUserResName, "global_grants.0.add_volumes", "false"),
					resource.TestCheckResourceAttr(testUserResName, "global_grants.0.cancel_account", "false"),
					resource.TestCheckResourceAttr(testUserResName, "global_grants.0.longview_subscription", "false"),
					resource.TestCheckResourceAttr(testUserResName, "linode_grant.#", "0"),
				),
			},
			{
				Config: tmpl.GrantsUpdate(t, username, email, instance),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testUserResName, "global_grants.0.account_access", "read_only"),
					resource.TestCheckResourceAttr(testUserResName, "global_grants.0.add_domains", "false"),
					resource.TestCheckResourceAttr(testUserResName, "global_grants.0.add_databases", "false"),
					resource.TestCheckResourceAttr(testUserResName, "global_grants.0.add_firewalls", "false"),
					resource.TestCheckResourceAttr(testUserResName, "global_grants.0.add_images", "true"),
					resource.TestCheckResourceAttr(testUserResName, "global_grants.0.add_linodes", "true"),
					resource.TestCheckResourceAttr(testUserResName, "global_grants.0.add_longview", "false"),
					resource.TestCheckResourceAttr(testUserResName, "global_grants.0.add_nodebalancers", "false"),
					resource.TestCheckResourceAttr(testUserResName, "global_grants.0.add_stackscripts", "false"),
					resource.TestCheckResourceAttr(testUserResName, "global_grants.0.add_volumes", "false"),
					resource.TestCheckResourceAttr(testUserResName, "global_grants.0.cancel_account", "false"),
					resource.TestCheckResourceAttr(testUserResName, "global_grants.0.longview_subscription", "false"),
					resource.TestCheckResourceAttr(testUserResName, "linode_grant.#", "1"),
					resource.TestCheckResourceAttr(testUserResName, "linode_grant.0.permissions", "read_write"),
				),
			},
		},
	})
}

func checkUserDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
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
