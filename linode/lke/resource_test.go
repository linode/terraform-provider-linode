package lke_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/helper"
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
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CheckLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasic(clusterName),
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
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CheckLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigManyPools(clusterName, "1.20"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceClusterName, "label", clusterName),
					resource.TestCheckResourceAttr(resourceClusterName, "region", "us-central"),
					resource.TestCheckResourceAttr(resourceClusterName, "k8s_version", "1.20"),
				),
			},
			{
				Config: resourceConfigManyPools(clusterName, "1.21"),
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
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CheckLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceClusterName, "label", clusterName),
					resource.TestCheckResourceAttr(resourceClusterName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.#", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.count", "3"),
				),
			},
			{
				Config: resourceConfigBasicUpdates(newClusterName),
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
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CheckLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceClusterName, "label", clusterName),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.#", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.count", "3"),
				),
			},
			{
				Config: resourceConfigComplexPools(newClusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceClusterName, "label", newClusterName),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.0.count", "2"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.1.count", "1"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.2.count", "2"),
					resource.TestCheckResourceAttr(resourceClusterName, "pool.3.count", "2"),
				),
			},
			{
				Config: resourceConfigBasic(clusterName),
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
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CheckLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasic(clusterName),
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
				Config: resourceConfigBasic(clusterName),
				Check:  resource.TestCheckResourceAttr(resourceClusterName, "pool.#", "1"),
			},
		},
	})
}

func resourceConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "linode_lke_cluster" "test" {
	label       = "%s"
	region      = "us-central"
	k8s_version = "1.20"
	tags        = ["test"]

	pool {
		type  = "g6-standard-2"
		count = 3
	}
}`, name)
}

func resourceConfigManyPools(name, k8sVersion string) string {
	return fmt.Sprintf(`
resource "linode_lke_cluster" "test" {
	label       = "%s"
	region      = "us-central"
	k8s_version = "%s"
	tags        = ["test"]

	pool {
		type  = "g6-standard-2"
		count = 3
	}

	pool {
		type = "g6-standard-2"
		count = 1
	}

	pool {
		type = "g6-standard-2"
		count = 1
	}
}`, name, k8sVersion)
}

func resourceConfigBasicUpdates(name string) string {
	return fmt.Sprintf(`
resource "linode_lke_cluster" "test" {
	label       = "%s"
	region      = "us-central"
	k8s_version = "1.20"
	tags        = ["test", "new_tag"]

	pool {
		type  = "g6-standard-2"
		count = 4
	}
}`, name)
}

func resourceConfigComplexPools(name string) string {
	return fmt.Sprintf(`
resource "linode_lke_cluster" "test" {
	label       = "%s"
	region      = "us-central"
	k8s_version = "1.20"
	tags        = ["test"]

	pool {
		type  = "g6-standard-2"
		count = 2
	}

	pool {
		type = "g6-standard-1"
		count = 1
	}

	pool {
		type = "g6-standard-1"
		count = 2
	}

	pool {
		type = "g6-standard-4"
		count = 2
	}
}`, name)
}
