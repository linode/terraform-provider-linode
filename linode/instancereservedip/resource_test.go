//go:build integration || instancereservedip

package instancereservedip_test

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/instancereservedip/tmpl"
)

const testInstanceIPResName = "linode_instance_reserved_ip.test"

var testRegion string

func init() {
	region, err := acceptance.GetRandomRegionWithCaps(nil, "core")
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccInstanceIP_addReservedIP(t *testing.T) {
	t.Parallel()

	var instance linodego.Instance
	name := acctest.RandomWithPrefix("tf_test")
	reservedIP := "50.116.51.242" // Replace with your actual reserved IP address
	testRegion = "us-east"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.AddReservedIP(t, name, testRegion, reservedIP),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists("linode_instance.foobar", &instance),
					resource.TestCheckResourceAttr(testInstanceIPResName, "address", reservedIP),
					resource.TestCheckResourceAttr(testInstanceIPResName, "public", "true"),
					resource.TestCheckResourceAttrSet(testInstanceIPResName, "linode_id"),
					resource.TestCheckResourceAttrSet(testInstanceIPResName, "gateway"),
					resource.TestCheckResourceAttrSet(testInstanceIPResName, "subnet_mask"),
					resource.TestCheckResourceAttrSet(testInstanceIPResName, "prefix"),
					resource.TestCheckResourceAttrSet(testInstanceIPResName, "rdns"),
					resource.TestCheckResourceAttr(testInstanceIPResName, "region", testRegion),
					resource.TestCheckResourceAttr(testInstanceIPResName, "type", "ipv4"),
				),
			},
		},
	})
}
