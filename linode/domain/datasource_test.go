//go:build integration || domain

package domain_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/domain/tmpl"
)

func TestAccDataSourceDomain_basic_smoke(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_domain.foobar"
	domainName := acctest.RandomWithPrefix("tf-test") + ".example"

	// TODO(ellisbenjamin) -- This test passes only because of the Destroy: true statement and needs attention.

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, domainName),
			},
			{
				Config: tmpl.DataBasic(t, domainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "domain", domainName),
					resource.TestCheckResourceAttr(resourceName, "type", "master"),
					resource.TestCheckResourceAttr(resourceName, "description", "tf-testing"),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tf_test"),
					resource.TestCheckResourceAttr(resourceName, "soa_email", "example@"+domainName),
					resource.TestCheckResourceAttrSet(resourceName, "retry_sec"),
					resource.TestCheckResourceAttrSet(resourceName, "expire_sec"),
				),
			},
			{
				Config: tmpl.DataByID(t, domainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "domain", domainName),
				),
				Destroy: true,
			},
			{
				Config:      tmpl.DataBasic(t, domainName),
				ExpectError: regexp.MustCompile(domainName + " was not found"),
			},
		},
	})
}
