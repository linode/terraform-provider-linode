//go:build integration || reservedip

package reservedip_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/reservedip/tmpl"
)

var testRegion string

func init() {
	region, err := acceptance.GetRandomRegionWithCaps([]string{"linodes"}, "core")
	if err != nil {
		panic(fmt.Sprintf("Error getting random region: %s", err))
	}
	testRegion = region
	fmt.Println(testRegion)
}

func TestAccResource_reserveIP(t *testing.T) {
	t.Parallel()

	resName := "linode_reserved_ip.test"
	instanceName := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.PreCheck(t)
		},
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.ReserveIP(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						return nil
					},
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttrSet(resName, "address"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					resource.TestCheckResourceAttr(resName, "reserved", "true"),
				),
			},
		},
	})

	t.Log("Finished TestAccResource_reserveIP")
}
