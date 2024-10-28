//go:build integration || networkingip

package networkingips_test

import (
	"log"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/networkingip/tmpl"
)

var testRegion string

func init() {
	region, err := acceptance.GetRandomRegionWithCaps([]string{"linodes"}, "core")
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccDataSourceNetworkingIP_basic(t *testing.T) {
	t.Parallel()

	resourceName := "linode_instance.foobar"
	dataResourceName := "data.linode_networking_ip.foobar"

	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, label, testRegion),
			},
			{
				Config: tmpl.DataBasic(t, label, testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataResourceName, "address", resourceName, "ip_address"),
					resource.TestCheckResourceAttrPair(dataResourceName, "linode_id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataResourceName, "region", resourceName, "region"),
					resource.TestMatchResourceAttr(dataResourceName, "gateway", regexp.MustCompile(`\.1$`)),
					resource.TestCheckResourceAttr(dataResourceName, "type", "ipv4"),
					resource.TestCheckResourceAttr(dataResourceName, "public", "true"),
					resource.TestCheckResourceAttr(dataResourceName, "prefix", "24"),
					resource.TestMatchResourceAttr(dataResourceName, "rdns", regexp.MustCompile(`.ip.linodeusercontent.com$`)),
				),
			},
		},
	})
}

func TestAccDataSourceNetworkingIP_list(t *testing.T) {
	t.Parallel()

	dataResourceName := "data.linode_networking_ip.list"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataList(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataResourceName, "ip_addresses.#"),
					resource.TestCheckResourceAttrSet(dataResourceName, "ip_addresses.0.address"),
					resource.TestCheckResourceAttrSet(dataResourceName, "ip_addresses.0.linode_id"),
					resource.TestCheckResourceAttrSet(dataResourceName, "ip_addresses.0.region"),
					resource.TestMatchResourceAttr(dataResourceName, "ip_addresses.0.gateway", regexp.MustCompile(`\.1$`)),
					resource.TestCheckResourceAttr(dataResourceName, "ip_addresses.0.type", "ipv4"),
					resource.TestCheckResourceAttr(dataResourceName, "ip_addresses.0.public", "true"),
					resource.TestCheckResourceAttr(dataResourceName, "ip_addresses.0.prefix", "24"),
					resource.TestMatchResourceAttr(dataResourceName, "ip_addresses.0.rdns", regexp.MustCompile(`.ip.linodeusercontent.com$`)),
					resource.TestCheckResourceAttrSet(dataResourceName, "ip_addresses.0.subnet_mask"),
					resource.TestCheckResourceAttrSet(dataResourceName, "ip_addresses.0.reserved"),
				),
			},
		},
	})
}

func TestAccDataSourceNetworkingIP_filterReserved(t *testing.T) {
	t.Parallel()

	dataResourceName := "data.linode_networking_ip.filtered"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataFilterReserved(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataResourceName, "ip_addresses.#"),
					resource.TestCheckResourceAttr(dataResourceName, "ip_addresses.0.reserved", "true"),
					resource.TestCheckResourceAttrSet(dataResourceName, "ip_addresses.0.address"),
					resource.TestCheckResourceAttrSet(dataResourceName, "ip_addresses.0.linode_id"),
					resource.TestCheckResourceAttrSet(dataResourceName, "ip_addresses.0.region"),
					resource.TestMatchResourceAttr(dataResourceName, "ip_addresses.0.gateway", regexp.MustCompile(`\.1$`)),
					resource.TestCheckResourceAttr(dataResourceName, "ip_addresses.0.type", "ipv4"),
					resource.TestCheckResourceAttr(dataResourceName, "ip_addresses.0.public", "true"),
					resource.TestCheckResourceAttr(dataResourceName, "ip_addresses.0.prefix", "24"),
					resource.TestMatchResourceAttr(dataResourceName, "ip_addresses.0.rdns", regexp.MustCompile(`.ip.linodeusercontent.com$`)),
					resource.TestCheckResourceAttrSet(dataResourceName, "ip_addresses.0.subnet_mask"),
				),
			},
		},
	})
}
