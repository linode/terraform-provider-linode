package linode

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLinodeProfile_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_profile.user"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLinodeProfileBasic(),
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
					resource.TestCheckResourceAttr(resourceName, "referrals.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "referrals.0.code"),
					resource.TestCheckResourceAttrSet(resourceName, "referrals.0.url"),
					resource.TestCheckResourceAttrSet(resourceName, "referrals.0.total"),
					resource.TestCheckResourceAttrSet(resourceName, "referrals.0.credit"),
					resource.TestCheckResourceAttrSet(resourceName, "referrals.0.completed"),
					resource.TestCheckResourceAttrSet(resourceName, "referrals.0.pending"),
				),
			},
		},
	})
}

func testDataSourceLinodeProfileBasic() string {
	return `data "linode_profile" "user" {}`
}
