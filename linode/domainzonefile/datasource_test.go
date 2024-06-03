//go:build integration || domainzonefile

package domainzonefile_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/domainzonefile/tmpl"
)

func TestAccDataSourceDomainZonefile_basic(t *testing.T) {
	datasourceName := "data.linode_domain_zonefile.foobar"
	domain := acctest.RandomWithPrefix("tf-test") + ".example"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, domain),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "domain_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "zone_file.0"),
				),
			},
		},
	})
}
