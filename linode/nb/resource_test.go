//go:build integration || nb

package nb_test

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/nb"
	"github.com/linode/terraform-provider-linode/v3/linode/nb/tmpl"
)

var testRegion string

func init() {
	resource.AddTestSweepers("linode_nodebalancer", &resource.Sweeper{
		Name: "linode_nodebalancer",
		F:    sweep,
	})

	region, err := acceptance.GetRandomRegionWithCaps([]string{"nodebalancers"}, "core")
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func sweep(prefix string) error {
	client, err := acceptance.GetTestClient()
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

func TestSmokeTests_nb(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{"TestAccResourceNodeBalancer_basic_smoke", TestAccResourceNodeBalancer_basic_smoke},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestAccResourceNodeBalancer_basic_smoke(t *testing.T) {
	t.Parallel()

	resName := "linode_nodebalancer.foobar"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkNodeBalancerDestroy,

		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, nodebalancerName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkNodeBalancerExists,
					resource.TestCheckResourceAttr(resName, "label", nodebalancerName),
					resource.TestCheckResourceAttr(resName, "client_conn_throttle", "20"),
					resource.TestCheckResourceAttr(resName, "client_udp_sess_throttle", "10"),
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
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"created", "updated", "firewall_id"}, // Ignore strict comparison for these attributes
			},
		},
	})
}

func TestAccResourceNodeBalancer_update(t *testing.T) {
	t.Parallel()

	resName := "linode_nodebalancer.foobar"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkNodeBalancerDestroy,

		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, nodebalancerName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkNodeBalancerExists,
					resource.TestCheckResourceAttr(resName, "label", nodebalancerName),
					resource.TestCheckResourceAttr(resName, "client_conn_throttle", "20"),
					resource.TestCheckResourceAttr(resName, "client_udp_sess_throttle", "10"),
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
				Config: tmpl.Updates(t, nodebalancerName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkNodeBalancerExists,
					resource.TestCheckResourceAttr(resName, "label", nodebalancerName+"_r"),
					resource.TestCheckResourceAttr(resName, "client_conn_throttle", "0"),
					resource.TestCheckResourceAttr(resName, "client_udp_sess_throttle", "5"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),

					resource.TestCheckResourceAttrSet(resName, "hostname"),
					resource.TestCheckResourceAttrSet(resName, "ipv4"),
					resource.TestCheckResourceAttrSet(resName, "ipv6"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttrSet(resName, "updated"),
					resource.TestCheckResourceAttr(resName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resName, "tags.0", "tf_test"),
					resource.TestCheckResourceAttr(resName, "tags.1", "tf_test_2"),
				),
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"created", "updated", "firewall_id"}, // Ignore strict comparison for these attributes
			},
		},
	})
}

func TestAccResourceNodeBalancer_firewall(t *testing.T) {
	t.Parallel()

	resName := "linode_nodebalancer.foobar"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkNodeBalancerDestroy,

		Steps: []resource.TestStep{
			{
				Config: tmpl.Firewall(t, nodebalancerName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkNodeBalancerExists,
					resource.TestCheckResourceAttr(resName, "label", nodebalancerName),
					resource.TestCheckResourceAttr(resName, "client_conn_throttle", "20"),
					acceptance.CheckResourceAttrGreaterThan(resName, "firewalls.#", 0),
					resource.TestCheckResourceAttr(resName, "firewalls.0.label", fmt.Sprintf("%v-fw", nodebalancerName)),
					resource.TestCheckResourceAttr(resName, "firewalls.0.inbound_policy", "DROP"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.inbound.#", "1"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.inbound.0.action", "ACCEPT"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.inbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.inbound.0.ports", "80"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.inbound.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.inbound.0.ipv4.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.inbound.0.ipv6.#", "1"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.inbound.0.ipv6.0", "::/0"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.outbound_policy", "DROP"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.outbound.#", "1"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.outbound.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.outbound.0.ports", "80"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.outbound.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.outbound.0.ipv4.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.outbound.0.ipv6.#", "1"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.outbound.0.ipv6.0", "2001:db8::/32"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.tags.#", "1"),
					resource.TestCheckResourceAttr(resName, "firewalls.0.tags.0", "test"),
				),
			},
			{
				Config: tmpl.FirewallUpdate(t, nodebalancerName, testRegion),
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

func TestAccResourceNodeBalancer_vpc(t *testing.T) {
	t.Parallel()

	resName := "linode_nodebalancer.test"
	nodebalancerName := acctest.RandomWithPrefix("tf-test")

	targetRegion, err := acceptance.GetRandomRegionWithCaps([]string{"NodeBalancers", "VPCs"}, "core")
	if err != nil {
		log.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkNodeBalancerDestroy,

		Steps: []resource.TestStep{
			{
				Config: tmpl.VPC(t, nodebalancerName, targetRegion),
				Check:  checkNodeBalancerExists,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("vpcs").AtSliceIndex(0).AtMapKey("subnet_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("vpcs").AtSliceIndex(0).AtMapKey("ipv4_range"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("vpcs").AtSliceIndex(0).AtMapKey("ipv4_range_auto_assign"),
						knownvalue.Bool(true),
					),
				},
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

	desiredState := map[string]attr.Value{
		"in":    types.Float64Value(1337.0),
		"out":   types.Float64Value(1338.0),
		"total": types.Float64Value(1339.0),
	}

	transferMap := oldState["transfer"].(map[string]interface{})

	newState := make(map[string]attr.Value)
	in, diag := nb.UpgradeResourceStateValue(transferMap["in"].(string))
	if diag != nil {
		t.Fatalf("error upgrading state: %v", diag.Detail())
	}
	newState["in"] = in

	out, diag := nb.UpgradeResourceStateValue(transferMap["out"].(string))
	if diag != nil {
		t.Fatalf("error upgrading state: %v", diag.Detail())
	}
	newState["out"] = out

	total, diag := nb.UpgradeResourceStateValue(transferMap["total"].(string))
	if diag != nil {
		t.Fatalf("error upgrading state: %v", diag.Detail())
	}
	newState["total"] = total

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

	desiredState := map[string]attr.Value{
		"in":    types.Float64Value(0.0),
		"out":   types.Float64Value(0.0),
		"total": types.Float64Value(0.0),
	}

	transferMap := oldState["transfer"].(map[string]interface{})

	newState := make(map[string]attr.Value)
	in, diag := nb.UpgradeResourceStateValue(transferMap["in"].(string))
	if diag != nil {
		t.Fatalf("error upgrading state: %v", diag.Detail())
	}
	newState["in"] = in

	out, diag := nb.UpgradeResourceStateValue(transferMap["out"].(string))
	if diag != nil {
		t.Fatalf("error upgrading state: %v", diag.Detail())
	}
	newState["out"] = out

	total, diag := nb.UpgradeResourceStateValue(transferMap["total"].(string))
	if diag != nil {
		t.Fatalf("error upgrading state: %v", diag.Detail())
	}
	newState["total"] = total

	if !reflect.DeepEqual(desiredState, newState) {
		t.Fatalf("expected %v, got %v", desiredState, newState)
	}
}

func checkNodeBalancerExists(s *terraform.State) error {
	client := acceptance.TestAccSDKv2Provider.Meta().(*helper.ProviderMeta).Client
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
	client := acceptance.TestAccSDKv2Provider.Meta().(*helper.ProviderMeta).Client
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
