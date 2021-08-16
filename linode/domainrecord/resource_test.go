package domainrecord_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func TestAccResourceDomainRecord_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_domain_record.foobar"
	domainRecordName := acctest.RandomWithPrefix("tf-test-")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkDomainRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasic(domainRecordName),
				Check: resource.ComposeTestCheckFunc(
					checkDomainRecordExists,
					resource.TestCheckResourceAttr(resName, "name", domainRecordName),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateID,
			},
		},
	})
}

func TestAccResourceDomainRecord_roundedTTLSec(t *testing.T) {
	t.Parallel()

	resName := "linode_domain_record.foobar"
	domainRecordName := acctest.RandomWithPrefix("tf-test-")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkDomainRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigWithTTL(domainRecordName, 299),
				Check: resource.ComposeTestCheckFunc(
					checkDomainRecordExists,
					resource.TestCheckResourceAttr(resName, "name", domainRecordName),
					resource.TestCheckResourceAttr(resName, "ttl_sec", "300"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateID,
			},
		},
	})
}

func TestAccResourceDomainRecord_ANoName(t *testing.T) {
	t.Parallel()

	resName := "linode_domain_record.foobar"
	domainName := acctest.RandomWithPrefix("tf-test-") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkDomainRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigANoName(domainName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "record_type", "A"),
					resource.TestCheckResourceAttr(resName, "name", ""),
				),
			},
		},
	})
}

func TestAccResourceDomainRecord_AAAANoName(t *testing.T) {
	t.Parallel()

	resName := "linode_domain_record.foobar"
	domainName := acctest.RandomWithPrefix("tf-test-") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkDomainRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigAAAANoName(domainName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "record_type", "AAAA"),
					resource.TestCheckResourceAttr(resName, "name", ""),
				),
			},
		},
	})
}

func TestAccResourceDomainRecord_CAANoName(t *testing.T) {
	t.Parallel()

	resName := "linode_domain_record.foobar"
	domainName := acctest.RandomWithPrefix("tf-test-") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkDomainRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigAANoName(domainName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "record_type", "CAA"),
					resource.TestCheckResourceAttr(resName, "name", ""),
				),
			},
		},
	})
}

func TestAccResourceDomainRecord_SRV(t *testing.T) {
	t.Parallel()

	resName := "linode_domain_record.foobar"
	domainName := acctest.RandomWithPrefix("tftest") + ".example"
	expectedName := "_myservice._tcp"
	expectedTarget := "mysubdomain." + domainName
	expectedTargetExternal := "subdomain.example.com"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkDomainRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigSRV(domainName, expectedTarget),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", expectedName),
					resource.TestCheckResourceAttr(resName, "target", expectedTarget),
					resource.TestCheckResourceAttr(resName, "record_type", "SRV"),
					resource.TestCheckResourceAttr(resName, "protocol", "tcp"),
				),
			},
			{
				Config: resourceConfigSRV(domainName, expectedTargetExternal),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", expectedName),
					resource.TestCheckResourceAttr(resName, "target", expectedTargetExternal),
					resource.TestCheckResourceAttr(resName, "record_type", "SRV"),
					resource.TestCheckResourceAttr(resName, "protocol", "tcp"),
				),
			},
		},
	})
}

func TestAccResourceDomainRecord_SRVNoFQDN(t *testing.T) {
	t.Parallel()

	resName := "linode_domain_record.foobar"
	domainName := acctest.RandomWithPrefix("tftest") + ".example"
	expectedName := "_myservice._tcp"
	expectedTarget := "mysubdomain." + domainName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkDomainRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigSRV(domainName, "mysubdomain"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", expectedName),
					resource.TestCheckResourceAttr(resName, "target", expectedTarget),
					resource.TestCheckResourceAttr(resName, "record_type", "SRV"),
					resource.TestCheckResourceAttr(resName, "protocol", "tcp"),
				),
			},
			{
				Config: resourceConfigSRV(domainName, "mysubdomainbutnew"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", expectedName),
					resource.TestCheckResourceAttr(resName, "target", "mysubdomainbutnew."+domainName),
					resource.TestCheckResourceAttr(resName, "record_type", "SRV"),
					resource.TestCheckResourceAttr(resName, "protocol", "tcp"),
				),
			},
		},
	})
}

func TestAccResourceDomainRecord_update(t *testing.T) {
	t.Parallel()

	domainRecordName := acctest.RandomWithPrefix("tf-test-")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkDomainRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasic(domainRecordName),
				Check: resource.ComposeTestCheckFunc(
					checkDomainRecordExists,
					resource.TestCheckResourceAttr("linode_domain_record.foobar", "name", domainRecordName),
				),
			},
			{
				Config: resourceConfigUpdates(domainRecordName),
				Check: resource.ComposeTestCheckFunc(
					checkDomainRecordExists,
					resource.TestCheckResourceAttr("linode_domain_record.foobar", "name", fmt.Sprintf("renamed-%s", domainRecordName)),
				),
			},
		},
	})
}

func importStateID(s *terraform.State) (string, error) {
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

func checkDomainRecordExists(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

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

func checkDomainRecordDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
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

func resourceConfigBasic(domainRecord string) string {
	return configBasic(domainRecord+".example") + fmt.Sprintf(`
resource "linode_domain_record" "foobar" {
	domain_id = "${linode_domain.foobar.id}"
	name = "%s"
	record_type = "CNAME"
	target = "target.%s.example"
}`, domainRecord, domainRecord)
}

func resourceConfigWithTTL(domainRecord string, ttlSec int) string {
	return configBasic(domainRecord+".example") + fmt.Sprintf(`
resource "linode_domain_record" "foobar" {
	domain_id = "${linode_domain.foobar.id}"
	name = "%s"
	record_type = "CNAME"
	target = "target.%s.example"
	ttl_sec = %d
}`, domainRecord, domainRecord, ttlSec)
}

func resourceConfigUpdates(domainRecord string) string {
	return configBasic(domainRecord+".example") + fmt.Sprintf(`
resource "linode_domain_record" "foobar" {
	domain_id = "${linode_domain.foobar.id}"
	name = "renamed-%s"
	record_type = "CNAME"
	target = "target.%s.example"
}`, domainRecord, domainRecord)
}

func resourceConfigSRVNoName(domainName string) string {
	return configBasic(domainName) + fmt.Sprintf(`
resource "linode_domain_record" "foobar" {
	domain_id = "${linode_domain.foobar.id}"
	record_type = "SRV"
	target = "target.%s"
}`, domainName)
}

func resourceConfigANoName(domainName string) string {
	return configBasic(domainName) + `
resource "linode_domain_record" "foobar" {
	domain_id = "${linode_domain.foobar.id}"
	record_type = "A"
	target = "192.168.1.1"
}`
}

func resourceConfigAAAANoName(domainName string) string {
	return configBasic(domainName) + `
resource "linode_domain_record" "foobar" {
	domain_id = "${linode_domain.foobar.id}"
	record_type = "AAAA"
	target = "2400:3f00::22"
}`
}

func resourceConfigAANoName(domainName string) string {
	return configBasic(domainName) + fmt.Sprintf(`
resource "linode_domain_record" "foobar" {
	domain_id = "${linode_domain.foobar.id}"
	record_type = "CAA"
	target = "target.%s"
	tag = "issue"
}`, domainName)
}

func resourceConfigSRV(domainName string, target string) string {
	return configBasic(domainName) + fmt.Sprintf(`
resource "linode_domain_record" "foobar" {
	domain_id = "${linode_domain.foobar.id}"
	record_type = "SRV"
	target      = "%s"
	service     = "myservice"
	protocol    = "tcp"
	port        = 1001
	priority    = 10
	weight      = 0
}`, target)
}

func configBasic(domain string) string {
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
