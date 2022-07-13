package domainzonefile_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/domainrecord/tmpl"
)

func TestAccDataSourceDomainRecord_basic(t *testing.T) {
	datasourceName := "data.linode_domain_zonefile"
	domain := acctest.RandomWithPrefix("recordtest") + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, domain),
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
