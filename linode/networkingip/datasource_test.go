package networkingip_test

import (
	"context"
	"log"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/networkingip/tmpl"
)

var testRegion string

func init() {
	region, err := acceptance.GetRandomRegionWithCaps([]string{"linodes"})
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

	provider, providerMap := acceptance.CreateTestProvider()
	acceptance.ModifyProviderMeta(provider, func(ctx context.Context, config *helper.ProviderMeta) error {
		config.Config.SkipInstanceReadyPoll = true
		config.Config.SkipInstanceDeletePoll = true
		return nil
	})

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: providerMap,
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
