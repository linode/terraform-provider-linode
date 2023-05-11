package instancenetworking_test

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/instancenetworking/tmpl"
)

const testInstanceNetworkResName = "data.linode_instance_networking.test"

var testRegion string

func init() {
	region, err := acceptance.GetRandomRegionWithCaps(nil)
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccDataSourceInstanceNetworking_basic(t *testing.T) {
	t.Parallel()

	var instance linodego.Instance

	name := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, name, testRegion),
			},
			{
				Config: tmpl.DataBasic(t, name, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists("linode_instance.foobar", &instance),
					resource.TestCheckResourceAttrSet(testInstanceNetworkResName, "ipv4.private.#"),
					resource.TestCheckResourceAttrSet(testInstanceNetworkResName, "ipv4.public.#"),
					resource.TestCheckResourceAttrSet(testInstanceNetworkResName, "ipv4.reserved.#"),
					resource.TestCheckResourceAttrSet(testInstanceNetworkResName, "ipv4.shared.#"),
					resource.TestCheckResourceAttrSet(testInstanceNetworkResName, "ipv6.global.#"),
					resource.TestCheckResourceAttrSet(testInstanceNetworkResName, "ipv6.link_local.%"),
					resource.TestCheckResourceAttrSet(testInstanceNetworkResName, "ipv6.slaac.%"),
				),
			},
		},
	})
}
