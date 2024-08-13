//go:build integration || placementgroups

package placementgroups_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/placementgroups/tmpl"
)

func TestAccDataSourcePlacementGroups_basic(t *testing.T) {
	t.Parallel()

	const dsAllName = "data.linode_placement_groups.all"
	const dsByLabelName = "data.linode_placement_groups.by-label"
	const dsByATName = "data.linode_placement_groups.by-placement-group-type"

	baseLabel := acctest.RandomWithPrefix("tf-test")

	testRegion, err := acceptance.GetRandomRegionWithCaps([]string{"Placement Group"}, "core")
	if err != nil {
		t.Error(fmt.Errorf("failed to get region with PG capability: %w", err))
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, baseLabel, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckResourceAttrGreaterThan(dsAllName, "placement_groups.#", 2),
					resource.TestCheckResourceAttrSet(dsAllName, "placement_groups.0.id"),
					resource.TestCheckResourceAttrSet(dsAllName, "placement_groups.0.label"),
					resource.TestCheckResourceAttrSet(dsAllName, "placement_groups.0.placement_group_type"),
					resource.TestCheckResourceAttrSet(dsAllName, "placement_groups.0.region"),
					resource.TestCheckResourceAttrSet(dsAllName, "placement_groups.0.is_compliant"),
					resource.TestCheckResourceAttrSet(dsAllName, "placement_groups.0.placement_group_policy"),
					resource.TestCheckResourceAttrSet(dsAllName, "placement_groups.0.members.#"),

					resource.TestCheckResourceAttr(dsByLabelName, "placement_groups.#", "1"),
					resource.TestCheckResourceAttrSet(dsByLabelName, "placement_groups.0.id"),
					resource.TestCheckResourceAttr(dsByLabelName, "placement_groups.0.label", baseLabel+"-1"),
					resource.TestCheckResourceAttr(dsByLabelName, "placement_groups.0.placement_group_type", "anti_affinity:local"),
					resource.TestCheckResourceAttr(dsByLabelName, "placement_groups.0.region", testRegion),
					resource.TestCheckResourceAttrSet(dsByLabelName, "placement_groups.0.is_compliant"),
					resource.TestCheckResourceAttr(dsByLabelName, "placement_groups.0.placement_group_policy", "strict"),
					resource.TestCheckResourceAttr(dsByLabelName, "placement_groups.0.members.#", "0"),

					acceptance.CheckResourceAttrGreaterThan(dsByATName, "placement_groups.#", 2),
					resource.TestCheckResourceAttrSet(dsByATName, "placement_groups.0.id"),
					resource.TestCheckResourceAttrSet(dsByATName, "placement_groups.0.label"),
					resource.TestCheckResourceAttrSet(dsByATName, "placement_groups.0.placement_group_type"),
					resource.TestCheckResourceAttrSet(dsByATName, "placement_groups.0.region"),
					resource.TestCheckResourceAttrSet(dsByATName, "placement_groups.0.is_compliant"),
					resource.TestCheckResourceAttrSet(dsByATName, "placement_groups.0.placement_group_policy"),
					resource.TestCheckResourceAttrSet(dsByATName, "placement_groups.0.members.#"),
				),
			},
		},
	})
}
