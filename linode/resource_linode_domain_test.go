package linode

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/chiefy/linodego"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccLinodeDomainBasic(t *testing.T) {
	t.Parallel()

	resName := "linode_domain.foobar"
	var domainName = fmt.Sprintf("tf-test-%s.example", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeDomainDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeDomainConfigBasic(domainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeDomainExists,
					resource.TestCheckResourceAttr(resName, "domain", domainName),
					resource.TestCheckResourceAttrSet(resName, "domain_type"),
					resource.TestCheckResourceAttrSet(resName, "soa_email"),
					resource.TestCheckResourceAttrSet(resName, "status"),
				),
			},

			resource.TestStep{
				ResourceName: resName,
				ImportState:  true,
			},
		},
	})
}

func TestAccLinodeDomainUpdate(t *testing.T) {
	t.Parallel()

	var domainName = fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeDomainDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeDomainConfigBasic(domainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeDomainExists,
					resource.TestCheckResourceAttr("linode_domain.foobar", "domain", domainName),
				),
			},
			resource.TestStep{
				Config: testAccCheckLinodeDomainConfigUpdates(domainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeDomainExists,
					resource.TestCheckResourceAttr("linode_domain.foobar", "domain", fmt.Sprintf("renamed-%s", domainName)),
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
		return fmt.Errorf("Failed to get Linode client")
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_domain" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Failed parsing %v to int", rs.Primary.ID)
		}
		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)

		}

		_, err = client.GetDomain(context.Background(), id)

		if err == nil {
			return fmt.Errorf("Linode Domain with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Failed to request Linode Domain with id %d", id)
		}
	}

	return nil
}

func testAccCheckLinodeDomainConfigBasic(domain string) string {
	return fmt.Sprintf(`
resource "linode_domain" "foobar" {
	domain = "%s"
	soa_email = "example@%s"
}`, domain, domain)
}

func testAccCheckLinodeDomainConfigUpdates(domain string) string {
	return fmt.Sprintf(`
resource "linode_domain" "foobar" {
	domain = "renamed-%s"
	soa_email = "example@%s"
}`, domain, domain)
}
