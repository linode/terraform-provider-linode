package lke_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/lke/tmpl"
)

const resourceClusterName = "linode_lke_cluster.test"

func init() {
	resource.AddTestSweepers("linode_lke_cluster", &resource.Sweeper{
		Name: "linode_lke_cluster",
		F:    sweep,
	})
}

func sweep(prefix string) error {
	client, err := acceptance.GetClientForSweepers()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	clusters, err := client.ListLKEClusters(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("Error getting clusters: %s", err)
	}
	for _, cluster := range clusters {
		if !acceptance.ShouldSweep(prefix, cluster.Label) {
			continue
		}
		if err := client.DeleteLKECluster(context.Background(), cluster.ID); err != nil {
			return fmt.Errorf("Error destroying LKE cluster %d during sweep: %s", cluster.ID, err)
		}
	}

	return nil
}

func checkLKEExists(cluster *linodego.LKECluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

		rs, ok := s.RootModule().Resources[resourceClusterName]
		if !ok {
			return fmt.Errorf("could not find resource %s", resourceClusterName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		found, err := client.GetLKECluster(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of Instance %s: %s", rs.Primary.Attributes["label"], err)
		}

		*cluster = *found
		return nil
	}
}

func TestAccResourceLKECluster_basic(t *testing.T) {
	t.Parallel()

	clusterName := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CheckLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceClusterName, "label", clusterName),
					resource.TestCheckResourceAttr(resourceClusterName, "region", "us-central"),
					resource.TestCheckResourceAttr(resourceClusterName, "k8s_version", "1.20"),
					resource.TestCheckResourceAttr(resourceClusterName, "status", "ready"),
					resource.TestCheckResourceAttr(resourceClusterName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.#", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.type", "g6-standard-2"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.count", "3"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.nodes.#", "3"),
					resource.TestCheckResourceAttr(resourceClusterName, "control_plane.#", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "control_plane.0.high_availability", "false"),
					resource.TestCheckResourceAttrSet(resourceClusterName, "id"),
					resource.TestCheckResourceAttrSet(resourceClusterName, "pool.0.id"),
					resource.TestCheckResourceAttrSet(resourceClusterName, "kubeconfig"),
				),
			},
		},
	})
}

func TestAccResourceLKECluster_k8sUpgrade(t *testing.T) {
	t.Parallel()

	clusterName := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CheckLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.ManyPools(t, clusterName, "1.20"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceClusterName, "label", clusterName),
					resource.TestCheckResourceAttr(resourceClusterName, "region", "us-central"),
					resource.TestCheckResourceAttr(resourceClusterName, "k8s_version", "1.20"),
				),
			},
			{
				Config: tmpl.ManyPools(t, clusterName, "1.21"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceClusterName, "label", clusterName),
					resource.TestCheckResourceAttr(resourceClusterName, "region", "us-central"),
					resource.TestCheckResourceAttr(resourceClusterName, "k8s_version", "1.21"),
				),
			},
		},
	})
}

func TestAccResourceLKECluster_basicUpdates(t *testing.T) {
	t.Parallel()

	clusterName := acctest.RandomWithPrefix("tf_test")
	newClusterName := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CheckLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceClusterName, "label", clusterName),
					resource.TestCheckResourceAttr(resourceClusterName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.#", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.count", "3"),
				),
			},
			{
				Config: tmpl.Updates(t, newClusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceClusterName, "label", newClusterName),
					resource.TestCheckResourceAttr(resourceClusterName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.count", "4"),
				),
			},
		},
	})
}

func TestAccResourceLKECluster_poolUpdates(t *testing.T) {
	t.Parallel()

	clusterName := acctest.RandomWithPrefix("tf_test")
	newClusterName := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CheckLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceClusterName, "label", clusterName),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.#", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.count", "3"),
				),
			},
			{
				Config: tmpl.ComplexPools(t, newClusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceClusterName, "label", newClusterName),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.count", "2"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.1.count", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.2.count", "2"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.3.count", "2"),
				),
			},
			{
				Config: tmpl.Basic(t, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceClusterName, "label", clusterName),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.#", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.count", "3"),
				),
			},
		},
	})
}

func TestAccResourceLKECluster_removeUnmanagedPool(t *testing.T) {
	t.Parallel()

	var cluster linodego.LKECluster

	clusterName := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CheckLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, clusterName),
				Check: resource.ComposeTestCheckFunc(
					checkLKEExists(&cluster),
					resource.TestCheckResourceAttr(resourceClusterName, "label", clusterName),
					resource.TestCheckResourceAttr(resourceClusterName, "status", "ready"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.#", "1"),
				),
			},
			{
				PreConfig: func() {
					client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
					if _, err := client.CreateLKEClusterPool(context.Background(), cluster.ID, linodego.LKEClusterPoolCreateOptions{
						Count: 1,
						Type:  "g6-standard-1",
					}); err != nil {
						t.Errorf("failed to create unmanaged pool for cluster %d: %s", cluster.ID, err)
					}

					pools, err := client.ListLKEClusterPools(context.Background(), cluster.ID, nil)
					if err != nil {
						t.Errorf("failed to get pools for cluster %d: %s", cluster.ID, err)
					}

					if len(pools) != 2 {
						t.Errorf("expected cluster to have 2 pools but got %d", len(pools))
					}
				},
				Config: tmpl.Basic(t, clusterName),
				Check:  resource.TestCheckResourceAttr(resourceClusterName, "pool.#", "1"),
			},
		},
	})
}

func TestAccResourceLKECluster_autoScaler(t *testing.T) {
	t.Parallel()

	clusterName := acctest.RandomWithPrefix("tf_test")
	//newClusterName := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CheckLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceClusterName, "label", clusterName),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.#", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.count", "3"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.autoscaler.#", "0"),
				),
			},
			{
				Config: tmpl.Autoscaler(t, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceClusterName, "label", clusterName),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.#", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.count", "3"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.autoscaler.#", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.autoscaler.0.min", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.autoscaler.0.max", "5"),
				),
			},
			{
				Config: tmpl.AutoscalerUpdates(t, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceClusterName, "label", clusterName),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.#", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.count", "3"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.autoscaler.#", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.autoscaler.0.min", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.autoscaler.0.max", "8"),
				),
			},
			{
				Config: tmpl.AutoscalerManyPools(t, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceClusterName, "label", clusterName),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.#", "2"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.count", "5"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.autoscaler.#", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.autoscaler.0.min", "3"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.autoscaler.0.max", "8"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.1.count", "3"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.1.autoscaler.#", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.1.autoscaler.0.min", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.1.autoscaler.0.max", "8"),
				),
			},
			{
				Config: tmpl.Basic(t, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceClusterName, "label", clusterName),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.#", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.count", "3"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.autoscaler.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceLKECluster_controlPlane(t *testing.T) {
	t.Parallel()

	clusterName := acctest.RandomWithPrefix("tf_test")
	//newClusterName := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CheckLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.ControlPlane(t, clusterName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceClusterName, "label", clusterName),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.#", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.count", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.autoscaler.#", "0"),
					resource.TestCheckResourceAttr(resourceClusterName, "control_plane.0.high_availability", "false"),
				),
			},
			{
				Config: tmpl.ControlPlane(t, clusterName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceClusterName, "label", clusterName),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.#", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.count", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.autoscaler.#", "0"),
					resource.TestCheckResourceAttr(resourceClusterName, "control_plane.0.high_availability", "true"),
				),
			},
			{
				Config: tmpl.ControlPlane(t, clusterName, false),

				// Expect a 400 response when attempting to disable HA
				ExpectError: regexp.MustCompile("\\[400]"),
			},
		},
	})
}
