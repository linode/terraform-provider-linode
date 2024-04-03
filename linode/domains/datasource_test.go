//go:build integration || domains

package domains_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/domains/tmpl"
)

func TestAccDataSourceDomains_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_domains.foo"

	domainName := acctest.RandomWithPrefix("tf-test") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, domainName),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckResourceAttrGreaterThan(resourceName, "domains.#", 1),
					resource.TestCheckResourceAttrSet(resourceName, "domains.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "domains.0.domain"),
					resource.TestCheckResourceAttrSet(resourceName, "domains.0.type"),
					resource.TestCheckResourceAttrSet(resourceName, "domains.0.description"),
					resource.TestCheckResourceAttrSet(resourceName, "domains.0.status"),
					resource.TestCheckResourceAttrSet(resourceName, "domains.0.tags.0"),
					resource.TestCheckResourceAttrSet(resourceName, "domains.0.soa_email"),
					resource.TestCheckResourceAttrSet(resourceName, "domains.0.retry_sec"),
					resource.TestCheckResourceAttrSet(resourceName, "domains.0.expire_sec"),
				),
			},
			{
				Config: tmpl.DataFilter(t, domainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "domains.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "domains.0.type", "master"),
					resource.TestCheckResourceAttr(resourceName, "domains.0.description", "tf-testing-master"),
				),
			},
			{
				Config: tmpl.DataAPIFilter(t, domainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "domains.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "domains.0.domain", domainName),
				),
			},
		},
	})
}
