package linode

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
)

func init() {
	resource.AddTestSweepers("linode_object_storage_bucket", &resource.Sweeper{
		Name: "linode_object_storage_bucket",
		F:    testSweepLinodeObjectStorageBucket,
	})
}

func testSweepLinodeObjectStorageBucket(prefix string) error {
	client, err := getClientForSweepers()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	listOpts := sweeperListOptions(prefix, "label")
	objectStorageBuckets, err := client.ListObjectStorageBuckets(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("Error getting object_storage_buckets: %s", err)
	}
	for _, objectStorageBucket := range objectStorageBuckets {
		if !shouldSweepAcceptanceTestResource(prefix, objectStorageBucket.Label) {
			continue
		}
		err := client.DeleteObjectStorageBucket(context.Background(), objectStorageBucket.Cluster, objectStorageBucket.Label)

		if err != nil {
			return fmt.Errorf("Error destroying %s during sweep: %s", objectStorageBucket.Label, err)
		}
	}

	return nil
}

func TestAccLinodeObjectStorageBucket_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_object_storage_bucket.foobar"
	var objectStorageBucketName = acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeObjectStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeObjectStorageBucketConfigBasic(objectStorageBucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeObjectStorageBucketExists,
					resource.TestCheckResourceAttr(resName, "label", objectStorageBucketName),
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

func TestAccLinodeObjectStorageBucket_dataSource(t *testing.T) {
	t.Parallel()

	resName := "linode_object_storage_bucket.foobar"
	var objectStorageBucketName = acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeObjectStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeObjectStorageBucketConfigDataSource(objectStorageBucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeObjectStorageBucketExists,
					resource.TestCheckResourceAttr(resName, "label", objectStorageBucketName),
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

func TestAccLinodeObjectStorageBucket_update(t *testing.T) {
	t.Parallel()

	var objectStorageBucketName = acctest.RandomWithPrefix("tf-test")
	resName := "linode_object_storage_bucket.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeObjectStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeObjectStorageBucketConfigBasic(objectStorageBucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeObjectStorageBucketExists,
					resource.TestCheckResourceAttr(resName, "label", objectStorageBucketName),
				),
			},
			{
				Config: testAccCheckLinodeObjectStorageBucketConfigUpdates(objectStorageBucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeObjectStorageBucketExists,
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s-renamed", objectStorageBucketName)),
				),
			},
		},
	})
}

func testAccCheckLinodeObjectStorageBucketExists(s *terraform.State) error {
	client := testAccProvider.Meta().(linodego.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_object_storage_bucket" {
			continue
		}

		cluster, label, err := decodeLinodeObjectStorageBucketID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %s, %s", rs.Primary.ID, err)
		}

		_, err = client.GetObjectStorageBucket(context.Background(), cluster, label)
		if err != nil {
			return fmt.Errorf("Error retrieving state of ObjectStorageBucket %s: %s", rs.Primary.Attributes["label"], err)
		}
	}

	return nil
}

func testAccCheckLinodeObjectStorageBucketDestroy(s *terraform.State) error {
	client, ok := testAccProvider.Meta().(linodego.Client)
	if !ok {
		return fmt.Errorf("Error getting Linode client")
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_object_storage_bucket" {
			continue
		}

		id := rs.Primary.ID
		cluster, label, err := decodeLinodeObjectStorageBucketID(id)
		if err != nil {
			return fmt.Errorf("Error parsing %s", id)
		}
		if label == "" {
			return fmt.Errorf("Would have considered %s as %s", id, label)

		}

		_, err = client.GetObjectStorageBucket(context.Background(), cluster, label)

		if err == nil {
			return fmt.Errorf("Linode ObjectStorageBucket with id %s still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error requesting Linode ObjectStorageBucket with id %s", id)
		}
	}

	return nil
}

func testAccCheckLinodeObjectStorageBucketConfigBasic(object_storage_bucket string) string {
	return fmt.Sprintf(`
resource "linode_object_storage_bucket" "foobar" {
	cluster = "us-east-1"
	label = "%s"
}`, object_storage_bucket)
}

func testAccCheckLinodeObjectStorageBucketConfigUpdates(object_storage_bucket string) string {
	return fmt.Sprintf(`
resource "linode_object_storage_bucket" "foobar" {
	cluster = "us-east-1"
	label = "%s-renamed"
}`, object_storage_bucket)
}

func testAccCheckLinodeObjectStorageBucketConfigDataSource(object_storage_bucket string) string {
	return fmt.Sprintf(`
data "linode_object_storage_cluster" "baz" {
	id = "us-east-1"
}

resource "linode_object_storage_bucket" "foobar" {
	cluster = data.linode_object_storage_cluster.baz.id
	label = "%s"
}`, object_storage_bucket)
}
