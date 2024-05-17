//go:build integration || placementgroupassignment

package placementgroupassignment_test

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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

//func resourceImportStateID(s *terraform.State) (string, error) {
//	for _, rs := range s.RootModule().Resources {
//		if rs.Type != "linode_firewall_device" {
//			continue
//		}
//
//		id, err := strconv.Atoi(rs.Primary.ID)
//		if err != nil {
//			return "", fmt.Errorf("Error parsing ID %v to int", rs.Primary.ID)
//		}
//
//		firewallID, err := strconv.Atoi(rs.Primary.Attributes["firewall_id"])
//		if err != nil {
//			return "", fmt.Errorf("Error parsing firewall_id %v to int", rs.Primary.Attributes["firewall_id"])
//		}
//		return fmt.Sprintf("%d,%d", firewallID, id), nil
//	}
//
//	return "", fmt.Errorf("Error finding firewall_device")
//}
