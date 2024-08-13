//go:build integration || lkenodepool

package lkenodepool_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	acceptanceTmpl "github.com/linode/terraform-provider-linode/v2/linode/acceptance/tmpl"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/lkenodepool/tmpl"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
)

var (
	clusterID  string
	k8sVersion string
	testRegion string
)

func init() {
	resource.AddTestSweepers("linode_lke_node_pool", &resource.Sweeper{
		Name: "linode_lke_node_pool",
		F:    sweep,
	})

	clusterID = os.Getenv("LINODE_TEST_CLUSTER_ID")

	if clusterID == "" {
		client, err := acceptance.GetTestClient()
		if err != nil {
			log.Fatalf("failed to get client: %s", err)
		}

		versions, err := client.ListLKEVersions(context.Background(), nil)
		if err != nil {
			log.Fatal(err)
		}

		k8sVersions := make([]string, len(versions))
		for i, v := range versions {
			k8sVersions[i] = v.ID
		}

		sort.Strings(k8sVersions)

		if len(k8sVersions) < 1 {
			log.Fatal("no k8s versions found")
		}

		k8sVersion = k8sVersions[len(k8sVersions)-1]

		region, err := acceptance.GetRandomRegionWithCaps([]string{"kubernetes"}, "core")
		if err != nil {
			log.Fatal(err)
		}

		testRegion = region
	}
}

func sweep(prefix string) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}
	clusterID, err := strconv.Atoi(os.Getenv("LINODE_TEST_CLUSTER_ID"))

	clusters, err := client.ListLKEClusters(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("Error getting clusters: %s", err)
	}
	for _, cluster := range clusters {
		if acceptance.ShouldSweep(prefix, cluster.Label) {
			if err := client.DeleteLKECluster(context.Background(), cluster.ID); err != nil {
				return fmt.Errorf("Error destroying LKE cluster %d during sweep: %s", cluster.ID, err)
			}
		} else {
			pools, err := client.ListLKENodePools(context.Background(), clusterID, nil)
			if err != nil {
				return fmt.Errorf("Error getting node pools: %s", err)
			}
			for _, pool := range pools {
				if containsTagWithPrefix(pool, prefix) {
					log.Printf("[DEBUG] Found a leaked node pool, clusterID: %d, poolID: %d. Deleting", clusterID, pool.ID)
					err := client.DeleteLKENodePool(context.Background(), clusterID, pool.ID)
					if err != nil {
						return fmt.Errorf("Error destroying nodepool %v during sweep: %s", pool.ID, err)
					}
				}
			}
		}
	}

	return nil
}

func TestAccResourceNodePool_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_lke_node_pool.foobar"
	clusterLabel := acctest.RandomWithPrefix("tf_test_")
	poolTag := acctest.RandomWithPrefix("tf_test_")

	templateData := createTemplateData()
	templateData.ClusterLabel = clusterLabel
	templateData.PoolTag = poolTag
	templateData.AutoscalerEnabled = true
	templateData.AutoscalerMin = 1
	templateData.AutoscalerMax = 2
	createConfig := createResourceConfig(t, &templateData)
	templateData.AutoscalerMin = 2
	templateData.AutoscalerMax = 3
	updateConfig := createResourceConfig(t, &templateData)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					checkNodePoolExists,
					resource.TestCheckResourceAttr(resName, "type", "g6-standard-1"),
					resource.TestCheckResourceAttrSet(resName, "disk_encryption"),
					resource.TestCheckResourceAttr(resName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resName, "tags.0", "external"),
					resource.TestCheckResourceAttr(resName, "tags.1", poolTag),
					resource.TestCheckResourceAttr(resName, "autoscaler.#", "1"),
					resource.TestCheckResourceAttr(resName, "autoscaler.0.min", "1"),
					resource.TestCheckResourceAttr(resName, "autoscaler.0.max", "2"),
					resource.TestCheckResourceAttr(resName, "node_count", "1"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: resourceImportStateID,
			},
			{
				Config: updateConfig,
				Check: resource.ComposeTestCheckFunc(
					checkNodePoolExists,
					resource.TestCheckResourceAttr(resName, "autoscaler.0.min", "2"),
					resource.TestCheckResourceAttr(resName, "autoscaler.0.max", "3"),
				),
			},
		},
	})
}

func TestAccResourceNodePool_disableAutoscaling(t *testing.T) {
	t.Parallel()

	resName := "linode_lke_node_pool.foobar"
	clusterLabel := acctest.RandomWithPrefix("tf_test_")
	poolTag := acctest.RandomWithPrefix("tf_test_")

	templateData := createTemplateData()
	templateData.ClusterLabel = clusterLabel
	templateData.PoolTag = poolTag
	templateData.AutoscalerEnabled = true
	templateData.AutoscalerMin = 1
	templateData.AutoscalerMax = 2
	createConfig := createResourceConfig(t, &templateData)
	templateData.AutoscalerEnabled = false
	templateData.NodeCount = 2
	updateConfig := createResourceConfig(t, &templateData)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					checkNodePoolExists,
					resource.TestCheckResourceAttr(resName, "type", "g6-standard-1"),
					resource.TestCheckResourceAttr(resName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resName, "tags.0", "external"),
					resource.TestCheckResourceAttr(resName, "tags.1", poolTag),
					resource.TestCheckResourceAttr(resName, "autoscaler.0.min", "1"),
					resource.TestCheckResourceAttr(resName, "autoscaler.0.max", "2"),
					resource.TestCheckResourceAttr(resName, "node_count", "1"),
				),
			},
			{
				Config: updateConfig,
				Check: resource.ComposeTestCheckFunc(
					checkNodePoolExists,
					resource.TestCheckResourceAttr(resName, "autoscaler.#", "0"),
					resource.TestCheckResourceAttr(resName, "node_count", "2"),
				),
			},
		},
	})
}

func TestAccResourceNodePool_enableAutoscaling(t *testing.T) {
	t.Parallel()

	resName := "linode_lke_node_pool.foobar"
	clusterLabel := acctest.RandomWithPrefix("tf_test_")
	poolTag := acctest.RandomWithPrefix("tf_test_")

	templateData := createTemplateData()
	templateData.ClusterLabel = clusterLabel
	templateData.PoolTag = poolTag
	templateData.AutoscalerEnabled = false
	templateData.NodeCount = 2
	createConfig := createResourceConfig(t, &templateData)
	templateData.AutoscalerEnabled = true
	templateData.AutoscalerMin = 1
	templateData.AutoscalerMax = 2
	updateConfig := createResourceConfig(t, &templateData)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					checkNodePoolExists,
					resource.TestCheckResourceAttr(resName, "type", "g6-standard-1"),
					resource.TestCheckResourceAttr(resName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resName, "tags.0", "external"),
					resource.TestCheckResourceAttr(resName, "tags.1", poolTag),
					resource.TestCheckResourceAttr(resName, "autoscaler.#", "0"),
					resource.TestCheckResourceAttr(resName, "node_count", "2"),
				),
			},
			{
				Config: updateConfig,
				Check: resource.ComposeTestCheckFunc(
					checkNodePoolExists,
					resource.TestCheckResourceAttr(resName, "autoscaler.#", "1"),
					resource.TestCheckResourceAttr(resName, "autoscaler.0.min", "1"),
					resource.TestCheckResourceAttr(resName, "autoscaler.0.max", "2"),
				),
			},
		},
	})
}

func TestAccResourceNodePool_update_type(t *testing.T) {
	t.Parallel()

	resName := "linode_lke_node_pool.foobar"
	clusterLabel := acctest.RandomWithPrefix("tf_test_")
	poolTag := acctest.RandomWithPrefix("tf_test_")

	templateData := createTemplateData()
	templateData.ClusterLabel = clusterLabel
	templateData.PoolTag = poolTag
	templateData.AutoscalerEnabled = false
	templateData.NodeCount = 1
	createConfig := createResourceConfig(t, &templateData)

	templateData.PoolNodeType = "g6-standard-2"
	updateConfig := createResourceConfig(t, &templateData)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					checkNodePoolExists,
					resource.TestCheckResourceAttr(resName, "type", "g6-standard-1"),
					resource.TestCheckResourceAttr(resName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resName, "tags.0", "external"),
					resource.TestCheckResourceAttr(resName, "tags.1", poolTag),
					resource.TestCheckResourceAttr(resName, "node_count", "1"),
					resource.TestCheckResourceAttr(resName, "type", "g6-standard-1"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: resourceImportStateID,
			},
			{
				Config: updateConfig,
				Check: resource.ComposeTestCheckFunc(
					checkNodePoolExists,
					resource.TestCheckResourceAttr(resName, "type", "g6-standard-2"),
				),
			},
		},
	})
}

func checkNodePoolExists(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
	clusterID, poolID, err := extractIDs(s)
	if err != nil {
		return err
	}
	_, err = client.GetLKENodePool(context.Background(), clusterID, poolID)
	if err != nil {
		return fmt.Errorf("Error retrieving state of node pool %d: %v", poolID, err)
	}
	return nil
}

func checkNodePoolDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
	clusterID, poolID, err := extractIDs(s)
	if err != nil {
		return err
	}
	_, err = client.GetLKENodePool(context.Background(), clusterID, poolID)

	if err == nil {
		return fmt.Errorf("Node Pool with id %d still exists in cluster %d", poolID, clusterID)
	}

	if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
		return fmt.Errorf("Error requesting Node Pool with id %d in cluster %d", poolID, clusterID)
	}
	return nil
}

func containsTagWithPrefix(pool linodego.LKENodePool, prefix string) bool {
	for _, tag := range pool.Tags {
		if strings.HasPrefix(tag, prefix) {
			return true
		}
	}
	return false
}

func createResourceConfig(t *testing.T, data *tmpl.TemplateData) string {
	return acceptanceTmpl.ProviderNoPoll(t) + tmpl.Generate(t, data)
}

func createTemplateData() tmpl.TemplateData {
	var data tmpl.TemplateData
	data.ClusterID = clusterID
	data.K8sVersion = k8sVersion
	data.Region = testRegion
	data.PoolNodeType = "g6-standard-1"
	return data
}

func extractIDs(s *terraform.State) (int, int, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_lke_node_pool" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return 0, 0, fmt.Errorf("Error parsing ID %v to int", rs.Primary.ID)
		}

		clusterID, err := strconv.Atoi(rs.Primary.Attributes["cluster_id"])
		if err != nil {
			return 0, 0, fmt.Errorf("Error parsing cluster_id %v to int", rs.Primary.Attributes["cluster_id"])
		}
		return clusterID, id, nil
	}

	return 0, 0, fmt.Errorf("Error finding lke_node_pool")
}

func resourceImportStateID(s *terraform.State) (string, error) {
	clusterID, id, err := extractIDs(s)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d,%d", clusterID, id), nil
}
