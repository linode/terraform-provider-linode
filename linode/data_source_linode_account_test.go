package linode

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLinodeAccount_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_account.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLinodeAccountBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "email"),
					resource.TestCheckResourceAttrSet(resourceName, "first_name"),
					resource.TestCheckResourceAttrSet(resourceName, "last_name"),
					resource.TestCheckResourceAttrSet(resourceName, "company"),
					resource.TestCheckResourceAttrSet(resourceName, "address_1"),
					resource.TestCheckResourceAttrSet(resourceName, "address_2"),
					resource.TestCheckResourceAttrSet(resourceName, "phone"),
					resource.TestCheckResourceAttrSet(resourceName, "city"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "country"),
					resource.TestCheckResourceAttrSet(resourceName, "zip"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "balance"),
				),
			},
		},
	})
}

func testDataSourceLinodeAccountBasic() string {
	return `data "linode_account" "foo" {}`
}
