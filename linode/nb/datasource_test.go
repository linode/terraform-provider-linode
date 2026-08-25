//go:build integration || nb

package nb_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/nb/tmpl"
)

func TestAccDataSourceNodeBalancer_basic(t *testing.T) {
	t.Parallel()

	resName := "data.linode_nodebalancer.foobar"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkNodeBalancerDestroy,

		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, nodebalancerName, testRegion),
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
					resource.TestCheckResourceAttr(resName, "transfer.#", "1"),
					resource.TestCheckResourceAttrSet(resName, "transfer.0.in"),
					resource.TestCheckResourceAttrSet(resName, "transfer.0.out"),
					resource.TestCheckResourceAttrSet(resName, "transfer.0.total"),
					resource.TestCheckResourceAttr(resName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resName, "tags.0", "tf_test"),
					resource.TestCheckResourceAttr(resName, "lke_cluster.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceNodeBalancer_firewalls(t *testing.T) {
	t.Parallel()

	resName := "data.linode_nodebalancer.foobar"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkNodeBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataFirewalls(t, nodebalancerName, testRegion),
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
					resource.TestCheckResourceAttr(resName, "transfer.#", "1"),
					resource.TestCheckResourceAttrSet(resName, "transfer.0.in"),
					resource.TestCheckResourceAttrSet(resName, "transfer.0.out"),
					resource.TestCheckResourceAttrSet(resName, "transfer.0.total"),
					resource.TestCheckResourceAttr(resName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resName, "tags.0", "tf_test"),
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
		},
	})
}

func TestAccDataSourceNodeBalancer_vpc(t *testing.T) {
	t.Parallel()

	dsName := "data.linode_nodebalancer.test"
	nodebalancerName := acctest.RandomWithPrefix("tf-test")

	targetRegion, err := acceptance.GetRandomRegionWithCaps([]linodego.RegionCapability{linodego.CapabilityNodeBalancers, linodego.CapabilityVPCs, linodego.CapabilityVPCDualStack}, "core")
	if err != nil {
		log.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkNodeBalancerDestroy,

		Steps: []resource.TestStep{
			{
				Config: tmpl.DataVPC(t, nodebalancerName, targetRegion),
				Check:  checkNodeBalancerExists,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						dsName,
						tfjsonpath.New("vpcs").AtSliceIndex(0).AtMapKey("subnet_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dsName,
						tfjsonpath.New("vpcs").AtSliceIndex(0).AtMapKey("ipv4_range"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dsName,
						tfjsonpath.New("vpcs").AtSliceIndex(0).AtMapKey("ipv6_range"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

// TestAccDataSourceNodeBalancer_lkeClusterAPL verifies that a NodeBalancer managed
// by an APL-enabled LKE cluster has the lke_cluster field populated in the data source.
//
// The LKE cluster is created via the API (not Terraform) so that we know the
// NodeBalancer ID before building the Terraform config, enabling a proper
// data source test with ConfigStateChecks.
func TestAccDataSourceNodeBalancer_lkeClusterAPL(t *testing.T) {
	acceptance.LongRunningTest(t)
	t.Parallel()

	acceptance.PreCheck(t)

	if lkeRegion == "" || k8sVersion == "" {
		t.Skip("no LKE-capable region or K8s version available")
	}

	client, err := acceptance.GetTestClient()
	if err != nil {
		t.Fatalf("failed to get test client: %s", err)
	}

	clusterName := acctest.RandomWithPrefix("tf-test")

	// Create APL LKE cluster via API.
	cluster, err := client.CreateLKECluster(context.Background(), linodego.LKEClusterCreateOptions{
		Label:      clusterName,
		Region:     lkeRegion,
		K8sVersion: k8sVersion,
		APLEnabled: true,
		NodePools: []linodego.LKENodePoolCreateOptions{
			{Type: "g6-dedicated-4", Count: 3},
		},
	})
	if err != nil {
		t.Fatalf("failed to create LKE cluster: %s", err)
	}

	t.Cleanup(func() {
		// Delete NB first (not in Terraform state, managed by APL).
		nbs, err := client.ListNodeBalancers(context.Background(), nil)
		if err != nil {
			t.Logf("cleanup: failed to list nodebalancers: %s", err)
		} else {
			for _, nb := range nbs {
				if nb.LKECluster != nil && nb.LKECluster.ID == cluster.ID {
					if err := client.DeleteNodeBalancer(context.Background(), nb.ID); err != nil {
						t.Logf("cleanup: failed to delete nodebalancer %d: %s", nb.ID, err)
					}
				}
			}
		}
		if err := client.DeleteLKECluster(context.Background(), cluster.ID); err != nil {
			t.Logf("cleanup: failed to delete LKE cluster %d: %s", cluster.ID, err)
		}
	})

	// Poll until the APL NodeBalancer appears (may take several minutes).
	const (
		pollInterval = 30 * time.Second
		pollTimeout  = 10 * time.Minute
	)

	var nbID int
	deadline := time.Now().Add(pollTimeout)
	for time.Now().Before(deadline) {
		nbs, err := client.ListNodeBalancers(context.Background(), nil)
		if err != nil {
			t.Fatalf("failed to list nodebalancers: %s", err)
		}
		for _, nb := range nbs {
			if nb.LKECluster != nil && nb.LKECluster.ID == cluster.ID {
				nbID = nb.ID
				break
			}
		}
		if nbID != 0 {
			break
		}
		log.Printf("[DEBUG] NodeBalancer for LKE cluster %d not yet available, retrying in %s...", cluster.ID, pollInterval)
		time.Sleep(pollInterval)
	}

	if nbID == 0 {
		t.Fatalf("timed out waiting for nodebalancer for APL LKE cluster %d after %s", cluster.ID, pollTimeout)
	}

	const dsName = "data.linode_nodebalancer.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: tmpl.LKEClusterData(t, nbID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						dsName,
						tfjsonpath.New("lke_cluster").AtSliceIndex(0).AtMapKey("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dsName,
						tfjsonpath.New("lke_cluster").AtSliceIndex(0).AtMapKey("label"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dsName,
						tfjsonpath.New("lke_cluster").AtSliceIndex(0).AtMapKey("type"),
						knownvalue.StringExact("lkecluster"),
					),
					statecheck.ExpectKnownValue(
						dsName,
						tfjsonpath.New("lke_cluster").AtSliceIndex(0).AtMapKey("url"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}
