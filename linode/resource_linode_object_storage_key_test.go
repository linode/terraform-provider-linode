package linode

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
)

func init() {
	resource.AddTestSweepers("linode_object_storage_key", &resource.Sweeper{
		Name: "linode_object_storage_key",
		F:    testSweepLinodeObjectStorageKey,
	})
}

func testSweepLinodeObjectStorageKey(prefix string) error {
	client, err := getClientForSweepers()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	objectStorageKeys, err := client.ListObjectStorageKeys(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("Error getting object storage keys: %s", err)
	}
	for _, objectStorageKey := range objectStorageKeys {
		if !shouldSweepAcceptanceTestResource(prefix, objectStorageKey.Label) || !strings.HasPrefix(objectStorageKey.Label, prefix) {
			continue
		}
		err := client.DeleteObjectStorageKey(context.Background(), objectStorageKey.ID)

		if err != nil {
			return fmt.Errorf("Error destroying %s during sweep: %s", objectStorageKey.Label, err)
		}
	}

	return nil
}

func TestAccLinodeObjectStorageKey_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_object_storage_key.foobar"
	var objectStorageKeyLabel = acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeObjectStorageKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeObjectStorageKeyConfigBasic(objectStorageKeyLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeObjectStorageKeyExists,
					testAccCheckLinodeObjectStorageKeySecretKeyAccessible,
					resource.TestCheckResourceAttr(resName, "label", objectStorageKeyLabel),
					resource.TestCheckResourceAttrSet(resName, "access_key"),
					resource.TestCheckResourceAttrSet(resName, "secret_key"),
				),
			},
		},
	})
}

func TestAccLinodeObjectStorageKey_update(t *testing.T) {
	t.Parallel()
	resName := "linode_object_storage_key.foobar"
	var objectStorageKeyLabel = acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeObjectStorageKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeObjectStorageKeyConfigBasic(objectStorageKeyLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeObjectStorageKeyExists,
					testAccCheckLinodeObjectStorageKeySecretKeyAccessible,
					resource.TestCheckResourceAttr(resName, "label", objectStorageKeyLabel),
					resource.TestCheckResourceAttrSet(resName, "access_key"),
				),
			},
			{
				Config: testAccCheckLinodeObjectStorageKeyConfigUpdates(objectStorageKeyLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeObjectStorageKeyExists,
					testAccCheckLinodeObjectStorageKeySecretKeyAccessible, // should be preserved in state
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s_renamed", objectStorageKeyLabel)),
					resource.TestCheckResourceAttrSet(resName, "access_key"),
				),
			},
		},
	})
}

func findObjectStorageKeyResources(s *terraform.State) []*terraform.ResourceState {
	keys := []*terraform.ResourceState{}
	for _, res := range s.RootModule().Resources {
		if res.Type != "linode_object_storage_key" {
			continue
		}
		keys = append(keys, res)
	}
	return keys
}

func testAccCheckLinodeObjectStorageKeySecretKeyAccessible(s *terraform.State) error {
	keys := findObjectStorageKeyResources(s)
	secret := keys[0].Primary.Attributes["secret_key"]

	if secret == "[REDACTED]" {
		return fmt.Errorf("Expected secret_key to be accessible but got '%s'", secret)
	}
	return nil
}

func testAccCheckLinodeObjectStorageKeyExists(s *terraform.State) error {
	client := testAccProvider.Meta().(linodego.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_object_storage_key" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		_, err = client.GetObjectStorageKey(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of Object Storage Key %s: %s", rs.Primary.Attributes["label"], err)
		}
	}

	return nil
}

func testAccCheckLinodeObjectStorageKeyDestroy(s *terraform.State) error {
	client, ok := testAccProvider.Meta().(linodego.Client)
	if !ok {
		return fmt.Errorf("Error getting Linode client")
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_object_storage_key" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}
		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)

		}

		_, err = client.GetObjectStorageKey(context.Background(), id)

		if err == nil {
			return fmt.Errorf("Linode Object Storage Key with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error requesting Linode Object Storage Key with id %d", id)
		}
	}

	return nil
}

func testAccCheckLinodeObjectStorageKeyConfigBasic(label string) string {
	return fmt.Sprintf(`
resource "linode_object_storage_key" "foobar" {
	label = "%s"
}`, label)
}

func testAccCheckLinodeObjectStorageKeyConfigUpdates(label string) string {
	return fmt.Sprintf(`
resource "linode_object_storage_key" "foobar" {
	label = "%s_renamed"
}`, label)
}
