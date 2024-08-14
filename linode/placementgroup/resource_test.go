//go:build integration || placementgroup

package placementgroup_test

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/placementgroup/tmpl"
)

var testRegion string

func init() {
	resource.AddTestSweepers("linode_placement_group", &resource.Sweeper{
		Name: "linode_placement_group",
		F:    sweep,
	})

	var err error
	testRegion, err = acceptance.GetRandomRegionWithCaps([]string{"Placement Group"}, "core")
	if err != nil {
		log.Fatal(fmt.Errorf("Error getting region: %s", err))
	}
}

func TestAccResourcePG_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_placement_group.foobar"
	label := acctest.RandomWithPrefix("tf-test")
	labelUpdated := label + "-updated"
	placementGroupType := string(linodego.PlacementGroupTypeAntiAffinityLocal)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkPGDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, testRegion, placementGroupType, "flexible"),
				Check: resource.ComposeTestCheckFunc(
					checkPGExists,
					resource.TestCheckResourceAttr(resName, "label", label),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					resource.TestCheckResourceAttr(resName, "placement_group_type", placementGroupType),
					resource.TestCheckResourceAttr(resName, "placement_group_policy", "flexible"),
					resource.TestCheckResourceAttrSet(resName, "id"),
				),
			},
			{
				Config: tmpl.Basic(t, labelUpdated, testRegion, placementGroupType, "flexible"),
				Check: resource.ComposeTestCheckFunc(
					checkPGExists,
					resource.TestCheckResourceAttr(resName, "label", labelUpdated),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					resource.TestCheckResourceAttr(resName, "placement_group_type", placementGroupType),
					resource.TestCheckResourceAttr(resName, "placement_group_policy", "flexible"),
					resource.TestCheckResourceAttrSet(resName, "id"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func checkPGExists(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_placement_group" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		_, err = client.GetPlacementGroup(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of Placement Group %s: %s", rs.Primary.Attributes["label"], err)
		}
	}

	return nil
}

func checkPGDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_placement_group" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}
		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)
		}

		_, err = client.GetPlacementGroup(context.Background(), id)

		if err == nil {
			return fmt.Errorf("Placement Group with id %d still exists", id)
		}

		if !linodego.IsNotFound(err) {
			return fmt.Errorf("Error requesting Linode Placement Group with id %d", id)
		}
	}

	return nil
}

func sweep(prefix string) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		log.Fatal(fmt.Errorf("Error getting client: %s", err))
	}

	listOpts := acceptance.SweeperListOptions(prefix, "label")
	pgs, err := client.ListPlacementGroups(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("Error getting Placement Groups: %s", err)
	}

	for _, pg := range pgs {
		if !acceptance.ShouldSweep(prefix, pg.Label) {
			continue
		}
		err := client.DeletePlacementGroup(context.Background(), pg.ID)
		if err != nil {
			return fmt.Errorf("Error destroying %s during sweep: %s", pg.Label, err)
		}
	}

	return nil
}
