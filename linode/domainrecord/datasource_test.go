//go:build integration || domainrecord

package domainrecord_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/domainrecord/tmpl"
)

func TestAccDataSourceDomainRecord_basic(t *testing.T) {
	datasourceName := "data.linode_domain_record.record"
	domain := acctest.RandomWithPrefix("recordtest") + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
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

func TestAccDataSourceDomainRecord_idLookup(t *testing.T) {
	datasourceName := "data.linode_domain_record.record"
	domain := acctest.RandomWithPrefix("idloikup") + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataID(t, domain),
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

func TestAccDataSourceDomainRecord_srv(t *testing.T) {
	datasourceName := "data.linode_domain_record.record"
	domain := acctest.RandomWithPrefix("srv") + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataSRV(t, domain),
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

func TestAccDataSourceDomainRecord_caa(t *testing.T) {
	datasourceName := "data.linode_domain_record.record"
	domain := acctest.RandomWithPrefix("caa") + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataCAA(t, domain),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "name", "caa_test"),
					resource.TestCheckResourceAttr(datasourceName, "type", "CAA"),
					resource.TestCheckResourceAttr(datasourceName, "tag", "issue"),
					resource.TestCheckResourceAttr(datasourceName, "target", "example.com"),
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
					resource.TestCheckResourceAttrSet(datasourceName, "domain_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "type"),
				),
			},
		},
	})
}
