package linode

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLinodeUser_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_user.user"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLinodeProfileBasic() + testDataSourceLinodeUserBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.linode_user.user", "username", resourceName, "username"),
					resource.TestCheckResourceAttrSet(resourceName, "email"),
				),
			},
			{
				Config:      testDataSourceLinodeUserDoesNotExist(),
				ExpectError: regexp.MustCompile(" was not found"),
			},
		},
	})
}

func testDataSourceLinodeUserBasic() string {
	return `
		data "linode_user" "user" {
			username = "${data.linode_profile.user.username}"
		}`
}

func testDataSourceLinodeUserDoesNotExist() string {
	return `
		data "linode_user" "user" {
			username = "does-not-exist"
		}`
}
