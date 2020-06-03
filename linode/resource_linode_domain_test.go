package linode

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
)

func init() {
	resource.AddTestSweepers("linode_domain", &resource.Sweeper{
		Name: "linode_domain",
		F:    testSweepLinodeDomain,
	})

}

func testSweepLinodeDomain(prefix string) error {
	client, err := getClientForSweepers()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	listOpts := sweeperListOptions(prefix, "domain")
	domains, err := client.ListDomains(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("Error getting domains: %s", err)
	}
	for _, domain := range domains {
		if !shouldSweepAcceptanceTestResource(prefix, domain.Domain) {
			continue
		}
		err := client.DeleteDomain(context.Background(), domain.ID)

		if err != nil {
			return fmt.Errorf("Error destroying %s during sweep: %s", domain.Domain, err)
		}
	}

	return nil
}

func TestAccLinodeDomain_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_domain.foobar"
	var domainName = acctest.RandomWithPrefix("tf-test") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeDomainConfigBasic(domainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeDomainExists,
					resource.TestCheckResourceAttr(resName, "domain", domainName),
					resource.TestCheckResourceAttrSet(resName, "type"),
					resource.TestCheckResourceAttrSet(resName, "soa_email"),
					resource.TestCheckResourceAttrSet(resName, "description"),
					resource.TestCheckResourceAttrSet(resName, "retry_sec"),
					resource.TestCheckResourceAttrSet(resName, "expire_sec"),
					resource.TestCheckResourceAttrSet(resName, "status"),
					resource.TestCheckNoResourceAttr(resName, "master_ips"),
					resource.TestCheckNoResourceAttr(resName, "axfr_ips"),
					resource.TestCheckResourceAttr(resName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resName, "tags.4106436895", "tf_test"),
				),
			},

			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLinodeDomain_update(t *testing.T) {
	t.Parallel()

	var domainName = acctest.RandomWithPrefix("tf-test") + ".example"
	var resName = "linode_domain.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeDomainConfigBasic(domainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeDomainExists,
					resource.TestCheckResourceAttr(resName, "domain", domainName),
				),
			},
			{
				Config: testAccCheckLinodeDomainConfigUpdates(domainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeDomainExists,
					resource.TestCheckResourceAttr(resName, "domain", fmt.Sprintf("renamed-%s", domainName)),
					resource.TestCheckResourceAttr(resName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resName, "tags.4106436895", "tf_test"),
					resource.TestCheckResourceAttr(resName, "tags.2667398925", "tf_test_2"),
				),
			},
		},
	})
}

func testAccCheckLinodeDomainExists(s *terraform.State) error {
	client := testAccProvider.Meta().(linodego.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_domain" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		_, err = client.GetDomain(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of Domain %s: %s", rs.Primary.Attributes["domain"], err)
		}
	}

	return nil
}

func testAccCheckLinodeDomainDestroy(s *terraform.State) error {
	client, ok := testAccProvider.Meta().(linodego.Client)
	if !ok {
		return fmt.Errorf("Error getting Linode client")
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_domain" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}
		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)

		}

		_, err = client.GetDomain(context.Background(), id)

		if err == nil {
			return fmt.Errorf("Linode Domain with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error requesting Linode Domain with id %d", id)
		}
	}

	return nil
}

func testAccCheckLinodeDomainConfigBasic(domain string) string {
	return fmt.Sprintf(`
resource "linode_domain" "foobar" {
	domain = "%s"
	type = "master"
	status = "active"
	soa_email = "example@%s"
	description = "tf-testing"
	tags = ["tf_test"]
}`, domain, domain)
}

func testAccCheckLinodeDomainConfigUpdates(domain string) string {
	return fmt.Sprintf(`
resource "linode_domain" "foobar" {
	domain = "renamed-%s"
	type = "master"
	status = "active"
	soa_email = "example@%s"
	description = "tf-testing"
	tags = ["tf_test", "tf_test_2"]
}`, domain, domain)
}
