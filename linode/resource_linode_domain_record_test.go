package linode

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
)

func TestAccLinodeDomainRecord_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_domain_record.foobar"
	domainRecordName := acctest.RandomWithPrefix("tf-test-")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeDomainRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeDomainRecordConfigBasic(domainRecordName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeDomainRecordExists,
					resource.TestCheckResourceAttr(resName, "name", domainRecordName),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateIDDomainRecord,
			},
		},
	})
}

func TestAccLinodeDomainRecord_ANoName(t *testing.T) {
	t.Parallel()

	resName := "linode_domain_record.foobar"
	domainName := acctest.RandomWithPrefix("tf-test-") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeDomainRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeDomainRecordConfigANoName(domainName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "record_type", "A"),
					resource.TestCheckResourceAttr(resName, "name", ""),
				),
			},
		},
	})
}

func TestAccLinodeDomainRecord_AAAANoName(t *testing.T) {
	t.Parallel()

	resName := "linode_domain_record.foobar"
	domainName := acctest.RandomWithPrefix("tf-test-") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeDomainRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeDomainRecordConfigAAAANoName(domainName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "record_type", "AAAA"),
					resource.TestCheckResourceAttr(resName, "name", ""),
				),
			},
		},
	})
}

func TestAccLinodeDomainRecord_CAANoName(t *testing.T) {
	t.Parallel()

	resName := "linode_domain_record.foobar"
	domainName := acctest.RandomWithPrefix("tf-test-") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeDomainRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeDomainRecordConfigCAANoName(domainName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "record_type", "CAA"),
					resource.TestCheckResourceAttr(resName, "name", ""),
				),
			},
		},
	})
}

func TestAccLinodeDomainRecord_SRV(t *testing.T) {
	t.Parallel()

	resName := "linode_domain_record.foobar"
	domainName := acctest.RandomWithPrefix("tftest") + ".example"
	expectedName := "_myservice._tcp"
	expectedTarget := "mysubdomain." + domainName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeDomainRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckLinodeDomainRecordConfigSRVNameSet(domainName),
				ExpectError: regexp.MustCompile(errLinodeDomainRecordSRVNameComputed),
			},
			{
				Config:      testAccCheckLinodeDomainRecordConfigSRVInvalidTarget(domainName),
				ExpectError: regexp.MustCompile("Target for SRV records must be the associated domain or a related FQDN."),
			},
			{
				Config: testAccCheckLinodeDomainRecordConfigSRVCorrected(domainName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", expectedName),
					resource.TestCheckResourceAttr(resName, "target", expectedTarget),
					resource.TestCheckResourceAttr(resName, "record_type", "SRV"),
					resource.TestCheckResourceAttr(resName, "protocol", "tcp"),
				),
			},
		},
	})
}

func testAccStateIDDomainRecord(s *terraform.State) (string, error) {
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

func TestAccLinodeDomainRecord_update(t *testing.T) {
	t.Parallel()

	domainRecordName := acctest.RandomWithPrefix("tf-test-")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeDomainRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeDomainRecordConfigBasic(domainRecordName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeDomainRecordExists,
					resource.TestCheckResourceAttr("linode_domain_record.foobar", "name", domainRecordName),
				),
			},
			{
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

func testAccCheckLinodeDomainRecordConfigSRVNoName(domainName string) string {
	return testAccCheckLinodeDomainConfigBasic(domainName) + fmt.Sprintf(`
resource "linode_domain_record" "foobar" {
	domain_id = "${linode_domain.foobar.id}"
	record_type = "SRV"
	target = "target.%s"
}`, domainName)
}

func testAccCheckLinodeDomainRecordConfigANoName(domainName string) string {
	return testAccCheckLinodeDomainConfigBasic(domainName) + `
resource "linode_domain_record" "foobar" {
	domain_id = "${linode_domain.foobar.id}"
	record_type = "A"
	target = "192.168.1.1"
}`
}

func testAccCheckLinodeDomainRecordConfigAAAANoName(domainName string) string {
	return testAccCheckLinodeDomainConfigBasic(domainName) + `
resource "linode_domain_record" "foobar" {
	domain_id = "${linode_domain.foobar.id}"
	record_type = "AAAA"
	target = "2400:3f00::22"
}`
}

func testAccCheckLinodeDomainRecordConfigCAANoName(domainName string) string {
	return testAccCheckLinodeDomainConfigBasic(domainName) + fmt.Sprintf(`
resource "linode_domain_record" "foobar" {
	domain_id = "${linode_domain.foobar.id}"
	record_type = "CAA"
	target = "target.%s"
	tag = "issue"
}`, domainName)
}

func testAccCheckLinodeDomainRecordConfigSRVNameSet(domainName string) string {
	return testAccCheckLinodeDomainConfigBasic(domainName) + fmt.Sprintf(`
resource "linode_domain_record" "foobar" {
	domain_id = "${linode_domain.foobar.id}"
	name = "named-%[1]s"
	record_type = "SRV"
	target      = "mysubdomain"
	service     = "myservice"
	protocol    = "tcp"
	port        = 1001
	priority    = 10
	weight      = 0
}`, domainName)
}

func testAccCheckLinodeDomainRecordConfigSRVInvalidTarget(domainName string) string {
	return testAccCheckLinodeDomainConfigBasic(domainName) + `
resource "linode_domain_record" "foobar" {
	domain_id = "${linode_domain.foobar.id}"
	record_type = "SRV"
	target      = "mysubdomain"
	service     = "myservice"
	protocol    = "tcp"
	port        = 1001
	priority    = 10
	weight      = 0
}`
}

func testAccCheckLinodeDomainRecordConfigSRVCorrected(domainName string) string {
	return testAccCheckLinodeDomainConfigBasic(domainName) + fmt.Sprintf(`
resource "linode_domain_record" "foobar" {
	domain_id = "${linode_domain.foobar.id}"
	record_type = "SRV"
	target      = "mysubdomain.%s"
	service     = "myservice"
	protocol    = "tcp"
	port        = 1001
	priority    = 10
	weight      = 0
}`, domainName)
}
