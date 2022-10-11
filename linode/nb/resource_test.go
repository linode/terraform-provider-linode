package nb_test

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"testing"

	"github.com/linode/terraform-provider-linode/linode/nb"
	"github.com/linode/terraform-provider-linode/linode/nb/tmpl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

var testRegion string

func init() {
	resource.AddTestSweepers("linode_nodebalancer", &resource.Sweeper{
		Name: "linode_nodebalancer",
		F:    sweep,
	})

	region, err := acceptance.GetRandomRegionWithCaps([]string{"nodebalancers"})
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func sweep(prefix string) error {
	client, err := acceptance.GetClientForSweepers()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	listOpts := acceptance.SweeperListOptions(prefix, "label")
	nodebalancers, err := client.ListNodeBalancers(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("Error getting instances: %s", err)
	}
	for _, nodebalancer := range nodebalancers {
		if nodebalancer.Label == nil || !acceptance.ShouldSweep(prefix, *nodebalancer.Label) {
			continue
		}
		err := client.DeleteNodeBalancer(context.Background(), nodebalancer.ID)
		if err != nil {
			return fmt.Errorf("Error destroying %v during sweep: %s", nodebalancer.Label, err)
		}
	}

	return nil
}

func TestAccResourceNodeBalancer_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_nodebalancer.foobar"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkNodeBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, nodebalancerName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkNodeBalancerExists,
					resource.TestCheckResourceAttr(resName, "label", nodebalancerName),
					resource.TestCheckResourceAttr(resName, "client_conn_throttle", "20"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),

					resource.TestCheckResourceAttrSet(resName, "hostname"),
					resource.TestCheckResourceAttrSet(resName, "ipv4"),
					resource.TestCheckResourceAttrSet(resName, "ipv6"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttrSet(resName, "updated"),
					resource.TestCheckResourceAttr(resName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resName, "tags.0", "tf_test"),
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

func TestAccResourceNodeBalancer_update(t *testing.T) {
	t.Parallel()

	resName := "linode_nodebalancer.foobar"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkNodeBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, nodebalancerName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkNodeBalancerExists,
					resource.TestCheckResourceAttr(resName, "label", nodebalancerName),
					resource.TestCheckResourceAttr(resName, "client_conn_throttle", "20"),
				),
			},
			{
				Config: tmpl.Updates(t, nodebalancerName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkNodeBalancerExists,
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s_r", nodebalancerName)),
					resource.TestCheckResourceAttr(resName, "client_conn_throttle", "0"),
					resource.TestCheckResourceAttr(resName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resName, "tags.0", "tf_test"),
					resource.TestCheckResourceAttr(resName, "tags.1", "tf_test_2"),
				),
			},
		},
	})
}

func TestLinodeNodeBalancer_UpgradeV0(t *testing.T) {
	t.Parallel()

	oldState := map[string]interface{}{
		"transfer": map[string]interface{}{
			"in":    "1337",
			"out":   "1338",
			"total": "1339",
		},
	}

	desiredState := map[string]interface{}{
		"transfer": []map[string]interface{}{
			{
				"in":    1337.0,
				"out":   1338.0,
				"total": 1339.0,
			},
		},
	}

	newState, err := nb.ResourceNodeBalancerV0Upgrade(context.Background(), oldState, nil)
	if err != nil {
		t.Fatalf("error migrating state: %v", err)
	}

	if !reflect.DeepEqual(desiredState, newState) {
		t.Fatalf("expected %v, got %v", desiredState, newState)
	}
}

func TestLinodeNodeBalancer_UpgradeV0Empty(t *testing.T) {
	t.Parallel()

	oldState := map[string]interface{}{
		"transfer": map[string]interface{}{
			"in":    "",
			"out":   "",
			"total": "",
		},
	}

	desiredState := map[string]interface{}{
		"transfer": []map[string]interface{}{
			{
				"in":    0.0,
				"out":   0.0,
				"total": 0.0,
			},
		},
	}

	newState, err := nb.ResourceNodeBalancerV0Upgrade(context.Background(), oldState, nil)
	if err != nil {
		t.Fatalf("error migrating state: %v", err)
	}

	if !reflect.DeepEqual(desiredState, newState) {
		t.Fatalf("expected %v, got %v", desiredState, newState)
	}
}

func checkNodeBalancerExists(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_nodebalancer" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		_, err = client.GetNodeBalancer(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of NodeBalancer %s: %s", rs.Primary.Attributes["label"], err)
		}
	}

	return nil
}

func checkNodeBalancerDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_nodebalancer" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}
		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)
		}

		_, err = client.GetNodeBalancer(context.Background(), id)

		if err == nil {
			return fmt.Errorf("NodeBalancer with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error requesting NodeBalancer with id %d", id)
		}
	}

	return nil
}
