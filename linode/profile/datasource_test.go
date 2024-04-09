//go:build integration || profile

package profile_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/profile/tmpl"
)

func TestAccDataSourceProfile_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_profile.user"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "email"),
					resource.TestCheckResourceAttrSet(resourceName, "timezone"),
					resource.TestCheckResourceAttrSet(resourceName, "username"),
					resource.TestCheckResourceAttrSet(resourceName, "email_notifications"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_whitelist_enabled"),
					resource.TestCheckResourceAttrSet(resourceName, "lish_auth_method"),
					resource.TestCheckResourceAttrSet(resourceName, "authorized_keys.#"),
					resource.TestCheckResourceAttrSet(resourceName, "restricted"),
					resource.TestCheckResourceAttrSet(resourceName, "two_factor_auth"),
					resource.TestCheckResourceAttrSet(resourceName, "referrals.0.total"),
					resource.TestCheckResourceAttrSet(resourceName, "referrals.0.credit"),
					resource.TestCheckResourceAttrSet(resourceName, "referrals.0.completed"),
					resource.TestCheckResourceAttrSet(resourceName, "referrals.0.pending"),
					resource.TestCheckResourceAttrSet(resourceName, "referrals.0.code"),
					resource.TestCheckResourceAttrSet(resourceName, "referrals.0.url"),
				),
			},
		},
	})
}
