package obj_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/obj/tmpl"
)

func TestAccResourceObject_basic(t *testing.T) {
	t.Parallel()

	validateObject := func(resourceName, key, content string) resource.TestCheckFunc {
		var object s3.GetObjectOutput

		return resource.ComposeTestCheckFunc(
			checkObjectExists(resourceName, &object),
			checkObjectBodyContains(&object, content),
			resource.TestCheckResourceAttr(resourceName, "key", key),
		)
	}

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

	bucketName := acctest.RandomWithPrefix("tf-test")
	keyName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, bucketName, keyName, content, contentSource.Name()),
				Check: resource.ComposeTestCheckFunc(
					validateObject(getObjectResourceName("basic"), "test_basic", content),
					validateObject(getObjectResourceName("base64"), "test_base64", content),
					validateObject(getObjectResourceName("source"), "test_source", content),
				),
			},
			{
				Config: tmpl.Updates(t, bucketName, keyName, contentUpdated, contentSourceUpdated.Name()),
				Check: resource.ComposeTestCheckFunc(
					validateObjectUpdates(getObjectResourceName("basic"), "test_basic", contentUpdated),
					validateObjectUpdates(getObjectResourceName("base64"), "test_base64", contentUpdated),
					validateObjectUpdates(getObjectResourceName("source"), "test_source", contentUpdated),
				),
			},
		},
	})
}

func getObject(rs *terraform.ResourceState) (*s3.GetObjectOutput, error) {
	bucket := rs.Primary.Attributes["bucket"]
	key := rs.Primary.Attributes["key"]
	etag := rs.Primary.Attributes["etag"]
	accessKey := rs.Primary.Attributes["access_key"]
	secretKey := rs.Primary.Attributes["secret_key"]
	cluster := rs.Primary.Attributes["cluster"]

	conn := s3.New(session.New(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:    aws.String(fmt.Sprintf(helper.LinodeObjectsEndpoint, cluster)),
	}))

	return conn.GetObject(
		&s3.GetObjectInput{
			Bucket:  &bucket,
			Key:     &key,
			IfMatch: &etag,
		})
}

func checkObjectExists(resourceName string, obj *s3.GetObjectOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("could not find resource %s in root module", resourceName)
		}
		key := rs.Primary.Attributes["key"]
		bucket := rs.Primary.Attributes["bucket"]

		out, err := getObject(rs)
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

		if _, err := getObject(rs); err == nil {
			return fmt.Errorf("object with %s Key still exists", key)
		}
	}

	return nil
}

func checkObjectBodyContains(obj *s3.GetObjectOutput, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		body, err := ioutil.ReadAll(obj.Body)
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
