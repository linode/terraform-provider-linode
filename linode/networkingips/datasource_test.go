//go:build integration || networkingip

package networkingips_test

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/networkingips/tmpl"
)

var testRegion string

func init() {
	region, err := acceptance.GetRandomRegionWithCaps([]string{"linodes"}, "core")
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccDataSourceNetworkingIP_list(t *testing.T) {
	t.Parallel()

	dataResourceName := "data.linode_networking_ips.list"

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
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[dataResourceName]
						if !ok {
							return fmt.Errorf("resource not found: %s", dataResourceName)
						}

						gateway := rs.Primary.Attributes["ip_addresses.0.gateway"]

						// Validate gateway: allow null (empty string) or a value ending in '.1'
						if gateway != "" && !regexp.MustCompile(`\.1$`).MatchString(gateway) {
							return fmt.Errorf("attribute ip_addresses.0.gateway has invalid value: %s", gateway)
						}
						return nil
					},
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

	dataResourceName := "data.linode_networking_ips.filtered"

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
