//go:build integration || region

package region_test

import (
	"context"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/region/tmpl"
)

var (
	testRegion string
	testLabel  string
)

func init() {
	region, err := acceptance.GetRandomRegionWithCaps([]string{"linodes"}, "core")
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region

	client, err := acceptance.GetTestClient()
	if err != nil {
		log.Fatal(err)
	}

	r, err := client.GetRegion(context.Background(), testRegion)
	if err != nil {
		log.Fatal(err)
	}

	testLabel = r.Label
}

func TestAccDataSourceRegion_basic(t *testing.T) {
	t.Parallel()

	regionID := testRegion
	label := testLabel
	resourceName := "data.linode_region.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, regionID, label),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "country"),
					resource.TestCheckResourceAttr(resourceName, "id", regionID),
					resource.TestCheckResourceAttr(resourceName, "label", label),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "site_type"),
					resource.TestCheckResourceAttrSet(resourceName, "resolvers.0.ipv4"),
					resource.TestCheckResourceAttrSet(resourceName, "resolvers.0.ipv6"),
					resource.TestCheckResourceAttrSet(resourceName, "placement_group_limits.0.maximum_pgs_per_customer"),
					resource.TestCheckResourceAttrSet(resourceName, "placement_group_limits.0.maximum_linodes_per_pg"),
					acceptance.CheckResourceAttrGreaterThan(resourceName, "capabilities.#", 0),
				),
			},
		},
	})
}
