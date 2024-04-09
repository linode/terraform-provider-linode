//go:build integration || users

package users_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/users/tmpl"
)

func TestAccDataSourceUsers_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_users.user"
	username := acctest.RandomWithPrefix("tf-test")
	email := username + "@example.com"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, username, email),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "users.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "users.0.username"),
					resource.TestCheckResourceAttrSet(resourceName, "users.0.email"),
					resource.TestCheckResourceAttrSet(resourceName, "users.0.tfa_enabled"),
				),
			},
		},
	})
}

func TestAccDataSourceUsers_clientFilter(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_users.user"
	username := acctest.RandomWithPrefix("tf-test")
	email := username + "@example.com"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataClientFilter(t, username, email),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "users.0.email", email),
					resource.TestCheckResourceAttr(resourceName, "users.0.restricted", "true"),
					resource.TestCheckResourceAttr(resourceName, "users.0.global_grants.#", "1"),
					resource.TestCheckNoResourceAttr(resourceName, "users.0.verified_phone_number"),
					resource.TestCheckNoResourceAttr(resourceName, "users.0.password_created"),
					resource.TestCheckResourceAttrSet(resourceName, "users.0.domain_grant.#"),
					resource.TestCheckResourceAttrSet(resourceName, "users.0.firewall_grant.#"),
					resource.TestCheckResourceAttrSet(resourceName, "users.0.image_grant.#"),
					resource.TestCheckResourceAttrSet(resourceName, "users.0.linode_grant.#"),
					resource.TestCheckResourceAttrSet(resourceName, "users.0.longview_grant.#"),
					resource.TestCheckResourceAttrSet(resourceName, "users.0.nodebalancer_grant.#"),
					resource.TestCheckResourceAttrSet(resourceName, "users.0.stackscript_grant.#"),
					resource.TestCheckResourceAttrSet(resourceName, "users.0.volume_grant.#"),
				),
			},
		},
	})
}

func TestAccDataSourceUsers_substring(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_users.user"
	username := acctest.RandomWithPrefix("tf-test")
	email := username + "@example.com"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataSubstring(t, username, email),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "2"),
					acceptance.CheckResourceAttrContains(resourceName, "users.0.username", username),
				),
			},
		},
	})
}
