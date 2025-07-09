//go:build integration || networkingip

package networkingips_test

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/networkingips/tmpl"
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
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataList(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataResourceName, "ip_addresses.#"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[dataResourceName]
						if !ok {
							return fmt.Errorf("resource not found: %s", dataResourceName)
						}

						numAddresses, err := strconv.Atoi(rs.Primary.Attributes["ip_addresses.#"])
						if err != nil {
							return fmt.Errorf("failed to parse ip_addresses.#: %v", err)
						}

						for i := 0; i < numAddresses; i++ {
							prefix := fmt.Sprintf("ip_addresses.%d.", i)

							// Check if all required fields are set
							if rs.Primary.Attributes[prefix+"gateway"] != "" &&
								rs.Primary.Attributes[prefix+"rdns"] != "" &&
								rs.Primary.Attributes[prefix+"address"] != "" &&
								rs.Primary.Attributes[prefix+"linode_id"] != "" &&
								rs.Primary.Attributes[prefix+"region"] != "" &&
								rs.Primary.Attributes[prefix+"type"] != "" &&
								rs.Primary.Attributes[prefix+"public"] != "" &&
								rs.Primary.Attributes[prefix+"prefix"] != "" &&
								rs.Primary.Attributes[prefix+"subnet_mask"] != "" &&
								rs.Primary.Attributes[prefix+"reserved"] != "" {

								// Perform assertions for the selected IP address
								if !regexp.MustCompile(`\.1$`).MatchString(rs.Primary.Attributes[prefix+"gateway"]) {
									return fmt.Errorf("attribute %sgateway has invalid value: %s", prefix, rs.Primary.Attributes[prefix+"gateway"])
								}

								if !regexp.MustCompile(`.ip.linodeusercontent.com$`).MatchString(rs.Primary.Attributes[prefix+"rdns"]) {
									return fmt.Errorf("attribute %srdns has invalid value: %s", prefix, rs.Primary.Attributes[prefix+"rdns"])
								}

								return nil
							}
						}

						return fmt.Errorf("no IP address found with all attributes set")
					},
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
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
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
