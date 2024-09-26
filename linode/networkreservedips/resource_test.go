//go:build integration || networkreservedips

package networkreservedips_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/networkreservedips/tmpl"
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

	t.Logf("Starting TestAccResource_reserveIP with resName: %s, instanceName: %s, region: %s", resName, instanceName, testRegion)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			t.Log("Running PreCheck")
			acceptance.PreCheck(t)
		},
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.ReserveIP(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						t.Log("Running Check function")
						return nil
					},
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttrSet(resName, "address"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
				),
			},
		},
	})

	t.Log("Finished TestAccResource_reserveIP")
}
