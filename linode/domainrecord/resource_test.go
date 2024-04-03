//go:build integration || domainrecord

package domainrecord_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/domainrecord/tmpl"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func TestAccResourceDomainRecord_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_domain_record.foobar"
	domainRecordName := acctest.RandomWithPrefix("tf-test-")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDomainRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, domainRecordName),
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
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDomainRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.TTL(t, domainRecordName, 299),
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

func TestAccResourceDomainRecord_TTLZero(t *testing.T) {
	t.Parallel()

	resName := "linode_domain_record.foobar"
	domainRecordName := acctest.RandomWithPrefix("tf-test-")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDomainRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.TTL(t, domainRecordName, 0),
				Check: resource.ComposeTestCheckFunc(
					checkDomainRecordExists,
					resource.TestCheckResourceAttr(resName, "name", domainRecordName),
					resource.TestCheckResourceAttr(resName, "ttl_sec", "0"),
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
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDomainRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.ANoName(t, domainName),
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
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDomainRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.AAAANoName(t, domainName),
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
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDomainRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.CAANoName(t, domainName),
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
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDomainRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.SRV(t, domainName, expectedTarget),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", expectedName),
					resource.TestCheckResourceAttr(resName, "target", expectedTarget),
					resource.TestCheckResourceAttr(resName, "record_type", "SRV"),
					resource.TestCheckResourceAttr(resName, "protocol", "tcp"),
				),
			},
			{
				Config: tmpl.SRV(t, domainName, expectedTargetExternal),
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
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDomainRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.SRV(t, domainName, "mysubdomain"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", expectedName),
					resource.TestCheckResourceAttr(resName, "target", expectedTarget),
					resource.TestCheckResourceAttr(resName, "record_type", "SRV"),
					resource.TestCheckResourceAttr(resName, "protocol", "tcp"),
				),
			},
			{
				Config: tmpl.SRV(t, domainName, "mysubdomainbutnew"),
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
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDomainRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, domainRecordName),
				Check: resource.ComposeTestCheckFunc(
					checkDomainRecordExists,
					resource.TestCheckResourceAttr("linode_domain_record.foobar", "name", domainRecordName),
				),
			},
			{
				Config: tmpl.Updates(t, domainRecordName),
				Check: resource.ComposeTestCheckFunc(
					checkDomainRecordExists,
					resource.TestCheckResourceAttr("linode_domain_record.foobar", "name", fmt.Sprintf("renamed-%s", domainRecordName)),
				),
			},
		},
	})
}

func TestAccResourceDomainRecord_reconcileName(t *testing.T) {
	t.Parallel()

	domainName := acctest.RandomWithPrefix("tf-test") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDomainRecordDestroy,
		Steps: []resource.TestStep{
			// Ensure there is no diff when using the domain name as the record name
			{
				Config: tmpl.WithDomain(t, domainName, domainName),
				Check: resource.ComposeTestCheckFunc(
					checkDomainRecordExists,
					resource.TestCheckResourceAttr("linode_domain_record.foobar", "name", domainName),
				),
			},
			{
				RefreshState: true,
				PlanOnly:     true,
			},

			// Ensure there is no diff when using an empty string as the record name
			{
				Config: tmpl.WithDomain(t, domainName, ""),
				Check: resource.ComposeTestCheckFunc(
					checkDomainRecordExists,
					resource.TestCheckResourceAttr("linode_domain_record.foobar", "name", domainName),
				),
			},
			{
				RefreshState: true,
				PlanOnly:     true,
			},

			// Ensure there is no diff when specifying a prefix with the domain
			{
				Config: tmpl.WithDomain(t, domainName, "test."+domainName),
				Check: resource.ComposeTestCheckFunc(
					checkDomainRecordExists,
					resource.TestCheckResourceAttr("linode_domain_record.foobar", "name", "test."+domainName),
				),
			},
			{
				RefreshState: true,
				PlanOnly:     true,
			},

			// Ensure there is no diff when specifying a prefix without the domain
			{
				Config: tmpl.WithDomain(t, domainName, "test"),
				Check: resource.ComposeTestCheckFunc(
					checkDomainRecordExists,
					resource.TestCheckResourceAttr("linode_domain_record.foobar", "name", "test"),
				),
			},
			{
				RefreshState: true,
				PlanOnly:     true,
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
