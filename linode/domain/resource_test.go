//go:build integration || domain

package domain_test

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
	"github.com/linode/terraform-provider-linode/v2/linode/domain/tmpl"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func init() {
	resource.AddTestSweepers("linode_domain", &resource.Sweeper{
		Name: "linode_domain",
		F:    sweep,
	})
}

func sweep(prefix string) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	listOpts := acceptance.SweeperListOptions(prefix, "domain")
	domains, err := client.ListDomains(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("Error getting domains: %s", err)
	}
	for _, domain := range domains {
		if !acceptance.ShouldSweep(prefix, domain.Domain) {
			continue
		}
		err := client.DeleteDomain(context.Background(), domain.ID)
		if err != nil {
			return fmt.Errorf("Error destroying %s during sweep: %s", domain.Domain, err)
		}
	}

	return nil
}

func TestSmokeTests_domain_resource(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{"TestAccResourceDomain_basic_smoke", TestAccResourceDomain_basic_smoke},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestAccResourceDomain_basic_smoke(t *testing.T) {
	t.Parallel()

	resName := "linode_domain.foobar"
	domainName := acctest.RandomWithPrefix("tf-test") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, domainName),
				Check: resource.ComposeTestCheckFunc(
					checkDomainExists,
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
					resource.TestCheckResourceAttr(resName, "tags.0", "tf_test"),
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

func TestAccResourceDomain_update(t *testing.T) {
	t.Parallel()

	domainName := acctest.RandomWithPrefix("tf-test") + ".example"
	resName := "linode_domain.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, domainName),
				Check: resource.ComposeTestCheckFunc(
					checkDomainExists,
					resource.TestCheckResourceAttr(resName, "domain", domainName),
				),
			},
			{
				Config: tmpl.Updates(t, domainName),
				Check: resource.ComposeTestCheckFunc(
					checkDomainExists,
					resource.TestCheckResourceAttr(resName, "domain", fmt.Sprintf("renamed-%s", domainName)),
					resource.TestCheckResourceAttr(resName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resName, "tags.0", "tf_test"),
					resource.TestCheckResourceAttr(resName, "tags.1", "tf_test_2"),
				),
			},
		},
	})
}

func TestAccResourceDomain_roundedDomainSecs(t *testing.T) {
	t.Parallel()

	domainName := acctest.RandomWithPrefix("tf-test") + ".example"
	resName := "linode_domain.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.RoundedSec(t, domainName),
				Check: resource.ComposeTestCheckFunc(
					checkDomainExists,
					resource.TestCheckResourceAttr(resName, "domain", domainName),
					resource.TestCheckResourceAttr(resName, "refresh_sec", "3600"),
					resource.TestCheckResourceAttr(resName, "retry_sec", "7200"),
					resource.TestCheckResourceAttr(resName, "ttl_sec", "300"),
					resource.TestCheckResourceAttr(resName, "expire_sec", "2419200"),
				),
			},
			{
				Config:            tmpl.RoundedSec(t, domainName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceDomain_zeroSecs(t *testing.T) {
	t.Parallel()

	domainName := acctest.RandomWithPrefix("tf-test") + ".example"
	resName := "linode_domain.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.ZeroSec(t, domainName),
				Check: resource.ComposeTestCheckFunc(
					checkDomainExists,
					resource.TestCheckResourceAttr(resName, "domain", domainName),
					resource.TestCheckResourceAttr(resName, "refresh_sec", "0"),
					resource.TestCheckResourceAttr(resName, "retry_sec", "0"),
					resource.TestCheckResourceAttr(resName, "ttl_sec", "0"),
					resource.TestCheckResourceAttr(resName, "expire_sec", "0"),
				),
			},
			{
				Config:            tmpl.ZeroSec(t, domainName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceDomain_updateIPs(t *testing.T) {
	t.Parallel()

	domainName := acctest.RandomWithPrefix("tf-test") + ".example"
	resName := "linode_domain.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.IPS(t, domainName),
				Check: resource.ComposeTestCheckFunc(
					checkDomainExists,
					resource.TestCheckResourceAttr(resName, "domain", domainName),
					resource.TestCheckResourceAttr(resName, "master_ips.0", "12.34.56.78"),
					resource.TestCheckResourceAttr(resName, "axfr_ips.0", "87.65.43.21"),
				),
			},
			{
				Config: tmpl.IPSUpdates(t, domainName),
				Check: resource.ComposeTestCheckFunc(
					checkDomainExists,
					resource.TestCheckResourceAttr(resName, "master_ips.#", "0"),
					resource.TestCheckResourceAttr(resName, "axfr_ips.#", "0"),
				),
			},
		},
	})
}

func checkDomainExists(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

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

func checkDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
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

func configRoundedSec(domain string) string {
	return fmt.Sprintf(`
resource "linode_domain" "foobar" {
	domain = "%s"
	type = "master"
	status = "active"
	soa_email = "example@%[1]s"
	description = "tf-testing"
	ttl_sec = 299
	refresh_sec = 600
	retry_sec = 3601
	expire_sec = 2419201
	tags = ["tf_test"]
}`, domain)
}
