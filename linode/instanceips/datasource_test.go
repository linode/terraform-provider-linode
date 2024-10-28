//go:build integration || instanceips

package instanceips_test

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/instanceips/tmpl"
)

const testInstanceIPsDataName = "data.linode_instance_ips.test"

var testRegion string

func init() {
	region, err := acceptance.GetRandomRegionWithCaps(nil, "core")
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccDataSourceInstanceIPs_basic(t *testing.T) {
	t.Parallel()

	var instance linodego.Instance

	name := acctest.RandomWithPrefix("tf_test")
	resourceName := "linode_instance.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, name, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttrSet(testInstanceIPsDataName, "id"),
					resource.TestCheckResourceAttrSet(testInstanceIPsDataName, "ipv4.public.0"),
					resource.TestCheckResourceAttrSet(testInstanceIPsDataName, "ipv4.private.#"),
					resource.TestCheckResourceAttrSet(testInstanceIPsDataName, "ipv6.slaac"),
					resource.TestCheckResourceAttrSet(testInstanceIPsDataName, "ipv6.link_local"),
					resource.TestCheckResourceAttr(testInstanceIPsDataName, "ipv6.global.#", "0"),
					resource.TestCheckResourceAttrPair(testInstanceIPsDataName, "id", resourceName, "id"),
				),
			},
		},
	})
}
