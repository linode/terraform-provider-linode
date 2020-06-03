package linode

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLinodeDomainRecord_basic(t *testing.T) {
	datasourceName := "data.linode_domain_record.record"
	domain := acctest.RandomWithPrefix("recordtest") + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLinodeDomainRecordConfigBasic(domain),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", "www"),
					resource.TestCheckResourceAttr(datasourceName, "type", "CNAME"),
					resource.TestCheckResourceAttr(datasourceName, "ttl_sec", "7200"),
					resource.TestCheckResourceAttr(datasourceName, "target", domain),
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
					resource.TestCheckResourceAttrSet(datasourceName, "domain_id"),
				),
			},
		},
	})
}

func TestAccDataSourceLinodeDomainRecord_idLookup(t *testing.T) {
	datasourceName := "data.linode_domain_record.record"
	domain := acctest.RandomWithPrefix("idloikup") + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLinodeDomainRecordConfigIDLookup(domain),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", "www"),
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
					resource.TestCheckResourceAttrSet(datasourceName, "domain_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "type"),
				),
			},
		},
	})
}

func TestAccDataSourceLinodeDomainRecord_srv(t *testing.T) {
	datasourceName := "data.linode_domain_record.record"
	domain := acctest.RandomWithPrefix("srv") + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLinodeDomainRecordConfigSRV(domain),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "type", "SRV"),
					resource.TestCheckResourceAttr(datasourceName, "port", "80"),
					resource.TestCheckResourceAttr(datasourceName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(datasourceName, "service", "sip"),
					resource.TestCheckResourceAttr(datasourceName, "weight", "5"),
					resource.TestCheckResourceAttr(datasourceName, "priority", "10"),
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
					resource.TestCheckResourceAttrSet(datasourceName, "domain_id"),
				),
			},
		},
	})
}

func TestAccDataSourceLinodeDomainRecord_caa(t *testing.T) {
	datasourceName := "data.linode_domain_record.record"
	domain := acctest.RandomWithPrefix("caa") + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLinodeDomainRecordConfigCAA(domain),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", "caa_test"),
					resource.TestCheckResourceAttr(datasourceName, "type", "CAA"),
					resource.TestCheckResourceAttr(datasourceName, "tag", "issue"),
					resource.TestCheckResourceAttr(datasourceName, "target", "test"),
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
					resource.TestCheckResourceAttrSet(datasourceName, "domain_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "type"),
				),
			},
		},
	})
}

func testAccDataSourceLinodeDomainRecordConfigBasic(domain string) string {
	return fmt.Sprintf(`
resource "linode_domain" "domain" {
	type = "master"
	domain = "%[1]s"
	soa_email = "example@%[1]s"
}

resource "linode_domain_record" "record" {
	domain_id = linode_domain.domain.id
	name = "www"
	record_type = "CNAME"
	target = "%[1]s"
	ttl_sec = 7200
}

data "linode_domain_record" "record" {
	domain_id = linode_domain.domain.id
	id = linode_domain_record.record.id
}
`, domain)
}

func testAccDataSourceLinodeDomainRecordConfigIDLookup(domain string) string {
	return fmt.Sprintf(`
resource "linode_domain" "domain" {
	type = "master"
	domain = "%[1]s"
	soa_email = "example@%[1]s"
}

resource "linode_domain_record" "record" {
	domain_id = linode_domain.domain.id
	name = "www"
	record_type = "CNAME"
	target = "%[1]s"
}

data "linode_domain_record" "record" {
	domain_id = linode_domain.domain.id
	name = linode_domain_record.record.name
}
`, domain)
}

func testAccDataSourceLinodeDomainRecordConfigSRV(domain string) string {
	return fmt.Sprintf(`
resource "linode_domain" "domain" {
	type = "master"
	domain = "%[1]s"
	soa_email = "example@%[1]s"
}

resource "linode_domain_record" "record" {
	domain_id = linode_domain.domain.id
	record_type = "SRV"
	target = "%[1]s"
	port = 80
	protocol = "tcp"
	service = "sip"
	weight = 5
	priority = 10
}

data "linode_domain_record" "record" {
	domain_id = linode_domain.domain.id
	name = linode_domain_record.record.name
}
`, domain)
}

func testAccDataSourceLinodeDomainRecordConfigCAA(domain string) string {
	return fmt.Sprintf(`
resource "linode_domain" "domain" {
	type = "master"
	domain = "%[1]s"
	soa_email = "example@%[1]s"
}

resource "linode_domain_record" "record" {
	name = "caa_test"
	domain_id = linode_domain.domain.id
	record_type = "CAA"
	tag = "issue"
	target = "test"
}

data "linode_domain_record" "record" {
	domain_id = linode_domain.domain.id
	id = linode_domain_record.record.id
}
`, domain)
}
