//go:build integration || objkey

package objkey_test

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/objkey/tmpl"
)

var testCluster string

func init() {
	resource.AddTestSweepers("linode_object_storage_key", &resource.Sweeper{
		Name: "linode_object_storage_key",
		F:    sweep,
	})

	cluster, err := acceptance.GetRandomOBJCluster()
	if err != nil {
		log.Fatal(err)
	}

	testCluster = cluster
}

func sweep(prefix string) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	objectStorageKeys, err := client.ListObjectStorageKeys(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("Error getting object storage keys: %s", err)
	}
	for _, objectStorageKey := range objectStorageKeys {
		if !acceptance.ShouldSweep(prefix, objectStorageKey.Label) || !strings.HasPrefix(objectStorageKey.Label, prefix) {
			continue
		}
		err := client.DeleteObjectStorageKey(context.Background(), objectStorageKey.ID)
		if err != nil {
			return fmt.Errorf("Error destroying %s during sweep: %s", objectStorageKey.Label, err)
		}
	}

	return nil
}

func TestAccResourceObjectKey_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_object_storage_key.foobar"
	objectStorageKeyLabel := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkObjectKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, objectStorageKeyLabel),
				Check: resource.ComposeTestCheckFunc(
					checkObjectKeyExists,
					checkObjectKeySecretAccessible,
					resource.TestCheckResourceAttr(resName, "label", objectStorageKeyLabel),
					resource.TestCheckResourceAttrSet(resName, "access_key"),
					resource.TestCheckResourceAttrSet(resName, "secret_key"),
					resource.TestCheckResourceAttr(resName, "limited", "false"),
				),
			},
		},
	})
}

func TestAccResourceObjectKey_limited(t *testing.T) {
	t.Parallel()

	resName := "linode_object_storage_key.foobar"
	objectStorageKeyLabel := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkObjectKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Limited(t, objectStorageKeyLabel, testCluster),
				Check: resource.ComposeTestCheckFunc(
					checkObjectKeyExists,
					checkObjectKeySecretAccessible,
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s_key", objectStorageKeyLabel)),
					resource.TestCheckResourceAttrSet(resName, "access_key"),
					resource.TestCheckResourceAttrSet(resName, "secret_key"),
					resource.TestCheckResourceAttr(resName, "limited", "true"),
					resource.TestCheckResourceAttr(resName, "bucket_access.#", "2"),
					resource.TestCheckResourceAttrSet(resName, "bucket_access.0.bucket_name"),
					resource.TestCheckResourceAttrSet(resName, "bucket_access.1.bucket_name"),
					resource.TestCheckResourceAttr(resName, "bucket_access.0.cluster", testCluster),
					resource.TestCheckResourceAttr(resName, "bucket_access.1.cluster", testCluster),
					resource.TestCheckResourceAttr(resName, "bucket_access.0.permissions", "read_only"),
					resource.TestCheckResourceAttr(resName, "bucket_access.1.permissions", "read_write"),
				),
			},
		},
	})
}

func TestAccResourceObjectKey_update(t *testing.T) {
	t.Parallel()
	resName := "linode_object_storage_key.foobar"
	objectStorageKeyLabel := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkObjectKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, objectStorageKeyLabel),
				Check: resource.ComposeTestCheckFunc(
					checkObjectKeyExists,
					checkObjectKeySecretAccessible,
					resource.TestCheckResourceAttr(resName, "label", objectStorageKeyLabel),
					resource.TestCheckResourceAttrSet(resName, "access_key"),
				),
			},
			{
				Config: tmpl.Updates(t, objectStorageKeyLabel),
				Check: resource.ComposeTestCheckFunc(
					checkObjectKeyExists,
					checkObjectKeySecretAccessible, // should be preserved in state
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s_renamed", objectStorageKeyLabel)),
					resource.TestCheckResourceAttrSet(resName, "access_key"),
				),
			},
		},
	})
}

func findObjectKeyResource(s *terraform.State) []*terraform.ResourceState {
	keys := []*terraform.ResourceState{}
	for _, res := range s.RootModule().Resources {
		if res.Type != "linode_object_storage_key" {
			continue
		}
		keys = append(keys, res)
	}
	return keys
}

func checkObjectKeySecretAccessible(s *terraform.State) error {
	keys := findObjectKeyResource(s)
	secret := keys[0].Primary.Attributes["secret_key"]

	if secret == "[REDACTED]" {
		return fmt.Errorf("Expected secret_key to be accessible but got '%s'", secret)
	}
	return nil
}

func checkObjectKeyExists(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

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

func checkObjectKeyDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
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
