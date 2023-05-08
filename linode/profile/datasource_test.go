package profile_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/profile/tmpl"
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
					resource.TestCheckResourceAttrSet(resourceName, "referrals.total"),
					resource.TestCheckResourceAttrSet(resourceName, "referrals.credit"),
					resource.TestCheckResourceAttrSet(resourceName, "referrals.completed"),
					resource.TestCheckResourceAttrSet(resourceName, "referrals.pending"),
					resource.TestCheckResourceAttrSet(resourceName, "referrals.code"),
					resource.TestCheckResourceAttrSet(resourceName, "referrals.url"),
				),
			},
		},
	})
}
