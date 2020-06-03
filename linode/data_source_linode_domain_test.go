package linode

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLinodeDomain_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_domain.foobar"
	domainName := acctest.RandomWithPrefix("tf-test") + ".example"

	// TODO(ellisbenjamin) -- This test passes only because of the Destroy: true statement and needs attention.

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeDomainConfigBasic(domainName),
			},
			{
				Config: testAccCheckLinodeDomainConfigBasic(domainName) + testDataSourceLinodeDomainBasic(domainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "domain", domainName),
					resource.TestCheckResourceAttr(resourceName, "type", "master"),
					resource.TestCheckResourceAttr(resourceName, "description", "tf-testing"),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "tags.4106436895", "tf_test"),
					resource.TestCheckResourceAttr(resourceName, "soa_email", "example@"+domainName),
					resource.TestCheckResourceAttrSet(resourceName, "retry_sec"),
					resource.TestCheckResourceAttrSet(resourceName, "expire_sec"),
				),
			},
			{
				Config: testAccCheckLinodeDomainConfigBasic(domainName) + testDataSourceLinodeDomainByID(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "domain", domainName),
				),
				Destroy: true,
			},
			{
				Config:      testDataSourceLinodeDomainBasic(domainName),
				ExpectError: regexp.MustCompile(domainName + " was not found"),
			},
		},
	})
}

func testDataSourceLinodeDomainBasic(domainName string) string {
	return fmt.Sprintf(`
data "linode_domain" "foobar" {
	domain = "%s"
}`, domainName)
}

func testDataSourceLinodeDomainByID() string {
	return `
data "linode_domain" "foobar" {
	id = "${linode_domain.foobar.id}"
}`
}
