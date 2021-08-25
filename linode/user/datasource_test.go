package user_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

func TestAccDataSourceUser_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_user.user"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: profileConfigBasic() + dataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "username"),
					resource.TestCheckResourceAttrSet(resourceName, "email"),
				),
			},
			{
				Config:      dataSourceConfigNoUser(),
				ExpectError: regexp.MustCompile(" was not found"),
			},
		},
	})
}

func profileConfigBasic() string {
	return `data "linode_profile" "user" {}`
}

func dataSourceConfigBasic() string {
	return `
		data "linode_user" "user" {
			username = "${data.linode_profile.user.username}"
		}`
}

func dataSourceConfigNoUser() string {
	return `
		data "linode_user" "user" {
			username = "does-not-exist"
		}`
}
