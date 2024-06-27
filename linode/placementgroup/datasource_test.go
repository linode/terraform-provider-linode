//go:build integration || placementgroup

package placementgroup_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/placementgroup/tmpl"
)

func TestAccDataSourcePlacementGroup_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_placement_group.test"

	label := acctest.RandomWithPrefix("tf-test")
	affinityType := string(linodego.AffinityTypeAntiAffinityLocal)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, label, testRegion, affinityType, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "label"),
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttrSet(resourceName, "affinity_type"),
					resource.TestCheckResourceAttrSet(resourceName, "is_compliant"),
					resource.TestCheckResourceAttrSet(resourceName, "is_strict"),
				),
			},
		},
	})
}
