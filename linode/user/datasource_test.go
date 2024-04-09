//go:build integration || user

package user_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/user/tmpl"
)

func TestAccDataSourceUser_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_user.user"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "username"),
					resource.TestCheckResourceAttrSet(resourceName, "email"),
					resource.TestCheckResourceAttrSet(resourceName, "tfa_enabled"),
				),
			},
			{
				Config:      tmpl.DataNoUser(t),
				ExpectError: regexp.MustCompile(" was not found"),
			},
		},
	})
}

func TestAccDataSourceUser_grants(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_user.test"

	username := acctest.RandomWithPrefix("tf-test")
	email := username + "@example.com"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataGrants(t, username, email),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "global_grants.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_grant.#"),
					resource.TestCheckResourceAttrSet(resourceName, "firewall_grant.#"),
					resource.TestCheckResourceAttrSet(resourceName, "image_grant.#"),
					resource.TestCheckResourceAttrSet(resourceName, "linode_grant.#"),
					resource.TestCheckResourceAttrSet(resourceName, "longview_grant.#"),
					resource.TestCheckResourceAttrSet(resourceName, "nodebalancer_grant.#"),
					resource.TestCheckResourceAttrSet(resourceName, "stackscript_grant.#"),
					resource.TestCheckResourceAttrSet(resourceName, "volume_grant.#"),
				),
			},
		},
	})
}
