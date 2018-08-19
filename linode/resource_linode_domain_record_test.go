package linode

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/linode/linodego"
)

func TestAccLinodeDomainRecordBasic(t *testing.T) {
	t.Parallel()

	resName := "linode_domain_record.foobar"
	var domainRecordName = acctest.RandomWithPrefix("tf-test-")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeDomainRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeDomainRecordConfigBasic(domainRecordName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeDomainRecordExists,
					resource.TestCheckResourceAttr(resName, "name", domainRecordName),
				),
			},

			resource.TestStep{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateIdDomainRecord,
			},
		},
	})
}

func testAccStateIdDomainRecord(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_domain_record" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return "", fmt.Errorf("Error parsing ID %v to int", rs.Primary.ID)
		}
		domainID, err := strconv.Atoi(rs.Primary.Attributes["domain_id"])
		if err != nil {
			return "", fmt.Errorf("Error parsing domain_id %v to int", rs.Primary.Attributes["domain_id"])
		}
		return fmt.Sprintf("%d,%d", domainID, id), nil
	}

	return "", fmt.Errorf("Error finding linode_domain_record")
}

func TestAccLinodeDomainRecordUpdate(t *testing.T) {
	t.Parallel()

	var domainRecordName = acctest.RandomWithPrefix("tf-test-")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeDomainRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeDomainRecordConfigBasic(domainRecordName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeDomainRecordExists,
					resource.TestCheckResourceAttr("linode_domain_record.foobar", "name", domainRecordName),
				),
			},
			resource.TestStep{
				Config: testAccCheckLinodeDomainRecordConfigUpdates(domainRecordName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeDomainRecordExists,
					resource.TestCheckResourceAttr("linode_domain_record.foobar", "name", fmt.Sprintf("renamed-%s", domainRecordName)),
				),
			},
		},
	})
}

func testAccCheckLinodeDomainRecordExists(s *terraform.State) error {
	client := testAccProvider.Meta().(linodego.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_domain_record" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}
		domainID, err := strconv.Atoi(rs.Primary.Attributes["domain_id"])
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.Attributes["domain_id"])
		}
		_, err = client.GetDomainRecord(context.Background(), domainID, id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of DomainRecord %s: %s", rs.Primary.Attributes["name"], err)
		}
	}

	return nil
}

func testAccCheckLinodeDomainRecordDestroy(s *terraform.State) error {
	client, ok := testAccProvider.Meta().(linodego.Client)
	if !ok {
		return fmt.Errorf("Error getting Linode client")
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_domain_record" {
			continue
		}
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}
		domainID, err := strconv.Atoi(rs.Primary.Attributes["domain_id"])
		if err != nil {
			return fmt.Errorf("Error parsing domain_id %v to int", rs.Primary.Attributes["domain_id"])
		}

		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)

		}

		_, err = client.GetDomainRecord(context.Background(), domainID, id)

		if err == nil {
			return fmt.Errorf("Linode DomainRecord with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error requesting Linode DomainRecord with id %d", id)
		}
	}

	return nil
}

func testAccCheckLinodeDomainRecordConfigBasic(domainRecord string) string {
	return testAccCheckLinodeDomainConfigBasic(domainRecord+".example") + fmt.Sprintf(`
resource "linode_domain_record" "foobar" {
	domain_id = "${linode_domain.foobar.id}"
	name = "%s"
	record_type = "CNAME"
	target = "target.%s.example"
}`, domainRecord, domainRecord)
}

func testAccCheckLinodeDomainRecordConfigUpdates(domainRecord string) string {
	return testAccCheckLinodeDomainConfigBasic(domainRecord+".example") + fmt.Sprintf(`
resource "linode_domain_record" "foobar" {
	domain_id = "${linode_domain.foobar.id}"
	name = "renamed-%s"
	record_type = "CNAME"
	target = "target.%s.example"
}`, domainRecord, domainRecord)
}
