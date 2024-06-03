//go:build integration || obj

package obj_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/obj/tmpl"
)

var testCluster string

func init() {
	cluster, err := acceptance.GetRandomOBJCluster()
	if err != nil {
		log.Fatal(err)
	}

	testCluster = cluster
}

func TestAccResourceObject_basic(t *testing.T) {
	t.Parallel()

	validateObjectUpdates := func(resourceName, key, content string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			validateObject(resourceName, key, content),
			resource.TestCheckResourceAttr(resourceName, "acl", "public-read"),
			resource.TestCheckResourceAttr(resourceName, "content_type", "text/plain"),
			resource.TestCheckResourceAttr(resourceName, "content_encoding", "utf8"),
			resource.TestCheckResourceAttr(resourceName, "content_language", "en"),
			resource.TestCheckResourceAttr(resourceName, "website_redirect", "test.com"),
			resource.TestCheckResourceAttr(resourceName, "force_destroy", "true"),
			resource.TestCheckResourceAttr(resourceName, "content_disposition", "attachment"),
			resource.TestCheckResourceAttr(resourceName, "cache_control", "max-age=2592000"),
			resource.TestCheckResourceAttr(resourceName, "metadata.foo", "bar"),
			resource.TestCheckResourceAttr(resourceName, "metadata.bar", "foo"),
		)
	}

	content := "testing123"
	contentUpdated := "testing456"

	contentSource := acceptance.CreateTempFile(t, "tf-test-obj-source", content)
	contentSourceUpdated := acceptance.CreateTempFile(t, "tf-test-obj-source-updated", contentUpdated)

	acceptance.RunTestRetry(t, 6, func(tRetry *acceptance.TRetry) {
		bucketName := acctest.RandomWithPrefix("tf-test")
		keyName := acctest.RandomWithPrefix("tf_test")

		resource.Test(tRetry, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
			CheckDestroy:             checkObjectDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.Basic(t, bucketName, testCluster, keyName, content, contentSource.Name()),
					Check: resource.ComposeTestCheckFunc(
						validateObject(getObjectResourceName("basic"), "test_basic", content),
						validateObject(getObjectResourceName("base64"), "test_base64", content),
						validateObject(getObjectResourceName("source"), "test_source", content),
					),
				},
				{
					Config: tmpl.Updates(t, bucketName, testCluster, keyName, contentUpdated, contentSourceUpdated.Name()),
					Check: resource.ComposeTestCheckFunc(
						validateObjectUpdates(getObjectResourceName("basic"), "test_basic", contentUpdated),
						validateObjectUpdates(getObjectResourceName("base64"), "test_base64", contentUpdated),
						validateObjectUpdates(getObjectResourceName("source"), "test_source", contentUpdated),
					),
				},
			},
		})
	})
}

func TestAccResourceObject_credsConfiged(t *testing.T) {
	t.Parallel()

	content := "test_creds_configed"

	acceptance.RunTestRetry(t, 6, func(tRetry *acceptance.TRetry) {
		bucketName := acctest.RandomWithPrefix("tf-test")
		keyName := acctest.RandomWithPrefix("tf_test")

		resource.Test(tRetry, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
			CheckDestroy:             checkObjectDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.CredsConfiged(t, bucketName, testCluster, keyName, content),
					Check: resource.ComposeTestCheckFunc(
						validateObject(getObjectResourceName("creds_configed"), "test_creds_configed", content),
					),
				},
			},
		})
	})
}

func TestAccResourceObject_tempKeys(t *testing.T) {
	t.Parallel()

	content := "test_temp_keys"

	acceptance.RunTestRetry(t, 6, func(tRetry *acceptance.TRetry) {
		bucketName := acctest.RandomWithPrefix("tf-test")
		keyName := acctest.RandomWithPrefix("tf_test")

		resource.Test(tRetry, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
			CheckDestroy:             checkObjectDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.TempKeys(t, bucketName, testCluster, keyName, content),
					Check: resource.ComposeTestCheckFunc(
						validateObject(getObjectResourceName("temp_keys"), "test_temp_keys", content),
					),
				},
			},
		})
	})
}

func getObject(ctx context.Context, rs *terraform.ResourceState) (*s3.GetObjectOutput, error) {
	bucket := rs.Primary.Attributes["bucket"]
	key := rs.Primary.Attributes["key"]
	etag := rs.Primary.Attributes["etag"]
	accessKey := rs.Primary.Attributes["access_key"]
	secretKey := rs.Primary.Attributes["secret_key"]
	endpoint := rs.Primary.Attributes["endpoint"]

	if accessKey == "" || secretKey == "" {
		client, err := acceptance.GetTestClient()
		if err != nil {
			return nil, fmt.Errorf("Error getting client: %s", err)
		}

		createOpts := linodego.ObjectStorageKeyCreateOptions{
			Label: fmt.Sprintf("temp_%s_%v", bucket, time.Now().Unix()),
			BucketAccess: &[]linodego.ObjectStorageKeyBucketAccess{{
				BucketName:  bucket,
				Cluster:     rs.Primary.Attributes["cluster"],
				Permissions: "read_write",
			}},
		}

		key, err := client.CreateObjectStorageKey(ctx, createOpts)
		if err != nil {
			return nil, err
		}

		accessKey = key.AccessKey
		secretKey = key.SecretKey

		defer func() {
			if err := client.DeleteObjectStorageKey(ctx, key.ID); err != nil {
				log.Printf("[WARN] Failed to clean up temporary object storage keys: %s\n", err)
			}
		}()
	}

	s3client, err := helper.S3Connection(ctx, endpoint, accessKey, secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get create s3 client: %w", err)
	}

	return s3client.GetObject(
		ctx,
		&s3.GetObjectInput{
			Bucket:  &bucket,
			Key:     &key,
			IfMatch: &etag,
		},
	)
}

func checkObjectExists(resourceName string, obj *s3.GetObjectOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("could not find resource %s in root module", resourceName)
		}
		key := rs.Primary.Attributes["key"]
		bucket := rs.Primary.Attributes["bucket"]

		out, err := getObject(context.Background(), rs)
		if err != nil {
			return fmt.Errorf("failed to get Bucket (%s) Object (%s): %s", bucket, key, err)
		}

		*obj = *out
		return nil
	}
}

func checkObjectDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_object_storage_object" {
			continue
		}

		key := rs.Primary.Attributes["key"]

		if _, err := getObject(context.Background(), rs); err == nil {
			return fmt.Errorf("object with %s Key still exists", key)
		}
	}

	return nil
}

func checkObjectBodyContains(obj *s3.GetObjectOutput, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		body, err := io.ReadAll(obj.Body)
		if err != nil {
			return fmt.Errorf("failed to read body: %s", err)
		}
		obj.Body.Close()

		if got := string(body); got != expected {
			return fmt.Errorf("expected body to be %q; got %q", expected, got)
		}
		return nil
	}
}

func getObjectResourceName(name string) string {
	return fmt.Sprintf("linode_object_storage_object.%s", name)
}

func validateObject(resourceName, key, content string) resource.TestCheckFunc {
	var object s3.GetObjectOutput

	return resource.ComposeTestCheckFunc(
		checkObjectExists(resourceName, &object),
		checkObjectBodyContains(&object, content),
		resource.TestCheckResourceAttr(resourceName, "key", key),
	)
}
