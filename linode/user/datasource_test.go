package user_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/user/tmpl"
)

func TestAccDataSourceUser_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_user.user"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "username"),
					resource.TestCheckResourceAttrSet(resourceName, "email"),
				),
			},
			{
				Config:      tmpl.DataNoUser(t),
				ExpectError: regexp.MustCompile(" was not found"),
			},
		},
	})
}
