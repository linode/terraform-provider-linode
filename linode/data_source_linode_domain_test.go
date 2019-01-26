package linode

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceLinodeDomain(t *testing.T) {
	t.Parallel()

	domainID := "1234567"
	resourceName := "data.linode_domain.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLinodeDomain(domainID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", domainID),
				),
			},
		},
	})
}

func testDataSourceLinodeDomain(domainID string) string {
	return fmt.Sprintf(`
data "linode_domain" "foobar" {
	id = "%s"
}`, domainID)
}
