//go:build integration || networkingip

package networkingip_test

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/networkingip/tmpl"
)

func init() {
	region, err := acceptance.GetRandomRegionWithCaps([]string{"linodes"}, "core")
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccResourceNetworkingIP_reserved(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")

	resourceName := "linode_networking_ip.reserved_ip"
	instanceResourceName := "linode_instance.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.NetworkingIPReservedAssigned(t, label, testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "reserved", "true"),
					resource.TestCheckResourceAttr(resourceName, "public", "true"),
					resource.TestCheckResourceAttr(resourceName, "type", "ipv4"),
					resource.TestCheckResourceAttrSet(resourceName, "address"),
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttrPair(resourceName, "linode_id", instanceResourceName, "id"),
				),
			},
		},
	})
}
