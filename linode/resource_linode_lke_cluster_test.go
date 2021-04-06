package linode

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
)

const testLKEClusterResName = "linode_lke_cluster.test"

func init() {
	resource.AddTestSweepers("linode_lke_cluster", &resource.Sweeper{
		Name: "linode_lke_cluster",
		F:    testSweepLinodeLKECluster,
	})
}

func TestReconcileLKEClusterPoolSpecs(t *testing.T) {
	for _, tc := range []struct {
		name             string
		specs            []linodeLKEClusterPoolSpec
		provisionedPools []linodego.LKEClusterPool

		expectedToDelete []int
		expectedToCreate []linodego.LKEClusterPoolCreateOptions
		expectedToUpdate map[int]linodego.LKEClusterPoolUpdateOptions
	}{
		{
			name: "no change",
			provisionedPools: []linodego.LKEClusterPool{
				{ID: 123, Type: "g6-standard-1", Count: 2},
			},
			specs: []linodeLKEClusterPoolSpec{
				{Type: "g6-standard-1", Count: 2},
			},
			expectedToUpdate: map[int]linodego.LKEClusterPoolUpdateOptions{},
		},
		{
			name: "upsize a single pool",
			provisionedPools: []linodego.LKEClusterPool{
				{ID: 123, Type: "g6-standard-1", Count: 2},
			},
			specs: []linodeLKEClusterPoolSpec{
				{Type: "g6-standard-1", Count: 3},
			},
			expectedToUpdate: map[int]linodego.LKEClusterPoolUpdateOptions{
				123: {Count: 3},
			},
		},
		{
			name: "change single pool type",
			provisionedPools: []linodego.LKEClusterPool{
				{ID: 123, Type: "g6-standard-1", Count: 2},
			},
			specs: []linodeLKEClusterPoolSpec{
				{Type: "g6-standard-2", Count: 2},
			},
			expectedToCreate: []linodego.LKEClusterPoolCreateOptions{
				{Type: "g6-standard-2", Count: 2},
			},
			expectedToDelete: []int{123},
			expectedToUpdate: map[int]linodego.LKEClusterPoolUpdateOptions{},
		},
		{
			name: "reuse cluster for resize",
			provisionedPools: []linodego.LKEClusterPool{
				{ID: 123, Type: "g6-standard-1", Count: 1},
				{ID: 124, Type: "g6-standard-1", Count: 10},
			},
			specs: []linodeLKEClusterPoolSpec{
				{Type: "g6-standard-1", Count: 9},  // bumped from 1 to 9
				{Type: "g6-standard-2", Count: 10}, // type changed
			},
			expectedToDelete: []int{123},
			expectedToUpdate: map[int]linodego.LKEClusterPoolUpdateOptions{
				124: {Count: 9},
			},
			expectedToCreate: []linodego.LKEClusterPoolCreateOptions{
				{Type: "g6-standard-2", Count: 10},
			},
		},
		{
			name: "competing resizes",
			provisionedPools: []linodego.LKEClusterPool{
				{ID: 123, Type: "g6-standard-3", Count: 3},
				{ID: 124, Type: "g6-standard-3", Count: 7},
				{ID: 126, Type: "g6-standard-3", Count: 4},
				{ID: 127, Type: "g6-standard-3", Count: 2},
			},
			specs: []linodeLKEClusterPoolSpec{
				{Type: "g6-standard-3", Count: 2},
				{Type: "g6-standard-3", Count: 9},
				{Type: "g6-standard-3", Count: 8},
				{Type: "g6-standard-3", Count: 2},
			},
			expectedToUpdate: map[int]linodego.LKEClusterPoolUpdateOptions{
				123: {Count: 2}, // -1
				124: {Count: 8}, // +1
				126: {Count: 9}, // +5
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			updates := reconcileLKEClusterPoolSpecs(tc.specs, tc.provisionedPools)
			if !reflect.DeepEqual(tc.expectedToCreate, updates.ToCreate) {
				t.Errorf("expected to create:\n%#v\ngot:\n%#v", tc.expectedToCreate, updates.ToCreate)
			}
			if !reflect.DeepEqual(tc.expectedToUpdate, updates.ToUpdate) {
				t.Errorf("expected to update:\n%#v\ngot:\n%#v", tc.expectedToUpdate, updates.ToUpdate)
			}
			if !reflect.DeepEqual(tc.expectedToDelete, updates.ToDelete) {
				t.Errorf("expected to delete:\n%#v\ngot:\n%#v", tc.expectedToDelete, updates.ToDelete)
			}
		})
	}
}

func testSweepLinodeLKECluster(prefix string) error {
	client, err := getClientForSweepers()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	clusters, err := client.ListLKEClusters(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("Error getting clusters: %s", err)
	}
	for _, cluster := range clusters {
		if !shouldSweepAcceptanceTestResource(prefix, cluster.Label) {
			continue
		}
		if err := client.DeleteLKECluster(context.Background(), cluster.ID); err != nil {
			return fmt.Errorf("Error destroying LKE cluster %d during sweep: %s", cluster.ID, err)
		}
	}

	return nil
}

func testAccCheckLinodeLKEClusterExists(cluster *linodego.LKECluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ProviderMeta).Client

		rs, ok := s.RootModule().Resources[testLKEClusterResName]
		if !ok {
			return fmt.Errorf("could not find resource %s", testLKEClusterResName)
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

func testAccCheckLinodeLKEClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderMeta).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_lke_cluster" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to parse LKE Cluster ID: %s", err)
		}

		if id == 0 {
			return fmt.Errorf("should not have LKE Cluster ID of 0")
		}

		if _, err = client.GetLKECluster(context.Background(), id); err == nil {
			return fmt.Errorf("should not find Linode ID %d existing after delete", id)
		} else if apiErr, ok := err.(*linodego.Error); !ok {
			return fmt.Errorf("expected API Error but got %#v", err)
		} else if apiErr.Code != 404 {
			return fmt.Errorf("expected an error 404 but got %#v", apiErr)
		}
	}

	return nil
}

func TestAccLinodeLKECluster_basic(t *testing.T) {
	t.Parallel()

	clusterName := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeLKEClusterBasic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testLKEClusterResName, "label", clusterName),
					resource.TestCheckResourceAttr(testLKEClusterResName, "region", "us-central"),
					resource.TestCheckResourceAttr(testLKEClusterResName, "k8s_version", "1.20"),
					resource.TestCheckResourceAttr(testLKEClusterResName, "status", "ready"),
					resource.TestCheckResourceAttr(testLKEClusterResName, "tags.#", "1"),
					resource.TestCheckResourceAttr(testLKEClusterResName, "pool.#", "1"),
					resource.TestCheckResourceAttr(testLKEClusterResName, "pool.0.type", "g6-standard-2"),
					resource.TestCheckResourceAttr(testLKEClusterResName, "pool.0.count", "3"),
					resource.TestCheckResourceAttr(testLKEClusterResName, "pool.0.nodes.#", "3"),
					resource.TestCheckResourceAttrSet(testLKEClusterResName, "id"),
					resource.TestCheckResourceAttrSet(testLKEClusterResName, "pool.0.id"),
					resource.TestCheckResourceAttrSet(testLKEClusterResName, "kubeconfig"),
				),
			},
		},
	})
}

func TestAccLinodeLKECluster_k8sUpgrade(t *testing.T) {
	t.Parallel()

	clusterName := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeLKEClusterManyPools(clusterName, "1.19"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testLKEClusterResName, "label", clusterName),
					resource.TestCheckResourceAttr(testLKEClusterResName, "region", "us-central"),
					resource.TestCheckResourceAttr(testLKEClusterResName, "k8s_version", "1.19"),
				),
			},
			{
				Config: testAccCheckLinodeLKEClusterManyPools(clusterName, "1.20"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testLKEClusterResName, "label", clusterName),
					resource.TestCheckResourceAttr(testLKEClusterResName, "region", "us-central"),
					resource.TestCheckResourceAttr(testLKEClusterResName, "k8s_version", "1.20"),
				),
			},
		},
	})
}

func TestAccLinodeLKECluster_basicUpdates(t *testing.T) {
	t.Parallel()

	clusterName := acctest.RandomWithPrefix("tf_test")
	newClusterName := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeLKEClusterBasic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testLKEClusterResName, "label", clusterName),
					resource.TestCheckResourceAttr(testLKEClusterResName, "tags.#", "1"),
					resource.TestCheckResourceAttr(testLKEClusterResName, "pool.#", "1"),
					resource.TestCheckResourceAttr(testLKEClusterResName, "pool.0.count", "3"),
				),
			},
			{
				Config: testAccCheckLinodeLKEClusterBasicUpdates(newClusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testLKEClusterResName, "label", newClusterName),
					resource.TestCheckResourceAttr(testLKEClusterResName, "tags.#", "2"),
					resource.TestCheckResourceAttr(testLKEClusterResName, "pool.0.count", "4"),
				),
			},
		},
	})
}

func TestAccLinodeLKECluster_poolUpdates(t *testing.T) {
	t.Parallel()

	clusterName := acctest.RandomWithPrefix("tf_test")
	newClusterName := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeLKEClusterBasic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testLKEClusterResName, "label", clusterName),
					resource.TestCheckResourceAttr(testLKEClusterResName, "pool.#", "1"),
					resource.TestCheckResourceAttr(testLKEClusterResName, "pool.0.count", "3"),
				),
			},
			{
				Config: testAccCheckLinodeLKEClusterComplexPools(newClusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testLKEClusterResName, "label", newClusterName),
					resource.TestCheckResourceAttr(testLKEClusterResName, "pool.0.count", "2"),
					resource.TestCheckResourceAttr(testLKEClusterResName, "pool.1.count", "1"),
					resource.TestCheckResourceAttr(testLKEClusterResName, "pool.2.count", "2"),
					resource.TestCheckResourceAttr(testLKEClusterResName, "pool.3.count", "2"),
				),
			},
			{
				Config: testAccCheckLinodeLKEClusterBasic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testLKEClusterResName, "label", clusterName),
					resource.TestCheckResourceAttr(testLKEClusterResName, "pool.#", "1"),
					resource.TestCheckResourceAttr(testLKEClusterResName, "pool.0.count", "3"),
				),
			},
		},
	})
}

func TestAccLinodeLKECluster_removeUnmanagedPool(t *testing.T) {
	t.Parallel()

	var cluster linodego.LKECluster

	clusterName := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeLKEClusterBasic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeLKEClusterExists(&cluster),
					resource.TestCheckResourceAttr(testLKEClusterResName, "label", clusterName),
					resource.TestCheckResourceAttr(testLKEClusterResName, "status", "ready"),
					resource.TestCheckResourceAttr(testLKEClusterResName, "pool.#", "1"),
				),
			},
			{
				PreConfig: func() {
					client := testAccProvider.Meta().(*ProviderMeta).Client
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
				Config: testAccCheckLinodeLKEClusterBasic(clusterName),
				Check:  resource.TestCheckResourceAttr(testLKEClusterResName, "pool.#", "1"),
			},
		},
	})
}

func testAccCheckLinodeLKEClusterBasic(name string) string {
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

func testAccCheckLinodeLKEClusterManyPools(name, k8sVersion string) string {
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

func testAccCheckLinodeLKEClusterBasicUpdates(name string) string {
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

func testAccCheckLinodeLKEClusterComplexPools(name string) string {
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
