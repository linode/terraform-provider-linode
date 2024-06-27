//go:build integration || placementgroupassignment

package placementgroupassignment_test

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/placementgroupassignment/tmpl"
)

var testRegion string

func init() {
	region, err := acceptance.GetRandomRegionWithCaps([]string{"Placement Group"})
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccResourcePlacementGroupAssignment_basic(t *testing.T) {
	t.Parallel()

	pgName := "linode_placement_group.test"
	instanceName := "linode_instance.test"
	assignmentName := "linode_placement_group_assignment.test"

	var instance linodego.Instance

	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreventPostDestroyRefresh: true,
		PreCheck:                  func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories:  acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, testRegion, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					acceptance.CheckInstanceExists(instanceName, &instance),
					resource.TestCheckResourceAttrSet(assignmentName, "id"),
					resource.TestCheckResourceAttrSet(assignmentName, "placement_group_id"),
					resource.TestCheckResourceAttrSet(assignmentName, "linode_id"),
				),
			},
			// Refresh the plan and make sure the assignment exists under the PG
			{
				Config: tmpl.Basic(t, label, testRegion, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(pgName, "members.#", "1"),
				),
			},
			// Attempt to import the assignment resource
			{
				ResourceName:      assignmentName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: resourceImportStateID,
			},
			// Drop the assignment resource
			{
				Config: tmpl.Basic(t, label, testRegion, false),
			},
			// Refresh the plan and make sure the assignment does not exist under the PG
			{
				Config: tmpl.Basic(t, label, testRegion, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(pgName, "members.#", "0"),
				),
			},
		},
	})
}

func resourceImportStateID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_placement_group_assignment" {
			continue
		}

		pgID, err := strconv.Atoi(rs.Primary.Attributes["placement_group_id"])
		if err != nil {
			return "", fmt.Errorf("Error parsing ID %v to int", rs.Primary.ID)
		}

		linodeID, err := strconv.Atoi(rs.Primary.Attributes["linode_id"])
		if err != nil {
			return "", fmt.Errorf("Error parsing ID %v to int", rs.Primary.ID)
		}

		return fmt.Sprintf("%d,%d", pgID, linodeID), nil
	}

	return "", fmt.Errorf("Error finding linode_placement_group_assignment")
}
