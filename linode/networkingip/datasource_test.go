package networkingip_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

func TestAccDataSourceNetworkingIP_basic(t *testing.T) {
	t.Parallel()

	resourceName := "linode_instance.foobar"
	dataResourceName := "data.linode_networking_ip.foobar"

	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: acceptance.AccTestWithProvider(dataSourceConfigBasic(label), map[string]interface{}{
					acceptance.SkipInstanceReadyPollKey: true,
				}),
			},
			{
				Config: acceptance.AccTestWithProvider(dataSourceConfigBasic(label), map[string]interface{}{
					acceptance.SkipInstanceReadyPollKey: true,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataResourceName, "address", resourceName, "ip_address"),
					resource.TestCheckResourceAttrPair(dataResourceName, "linode_id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataResourceName, "region", resourceName, "region"),
					resource.TestMatchResourceAttr(dataResourceName, "gateway", regexp.MustCompile(`\.1$`)),
					resource.TestCheckResourceAttr(dataResourceName, "type", "ipv4"),
					resource.TestCheckResourceAttr(dataResourceName, "public", "true"),
					resource.TestCheckResourceAttr(dataResourceName, "prefix", "24"),
					resource.TestMatchResourceAttr(dataResourceName, "rdns", regexp.MustCompile(`\.members\.linode\.com$`)),
				),
			},
		},
	})
}

func dataSourceConfigBasic(label string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	group = "tf_test"
	image = "linode/alpine3.12"
	type = "g6-standard-1"
	region = "us-east"
}

data "linode_networking_ip" "foobar" {
	address = "${linode_instance.foobar.ip_address}"
}`, label)
}
