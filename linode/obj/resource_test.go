package obj_test

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
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

const objectResourceName = "linode_object_storage_object.object"

func TestAccResourceObject_basic(t *testing.T) {
	t.Parallel()

	content := "testing123"
	bucketName := acctest.RandomWithPrefix("tf-test")
	keyName := acctest.RandomWithPrefix("tf_test")

	var object s3.GetObjectOutput

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, bucketName, keyName, content),
				Check: resource.ComposeTestCheckFunc(
					checkObjectExists(&object),
					checkObjectBodyContains(&object, content),
					resource.TestCheckResourceAttr(objectResourceName, "key", "test"),
				),
			},
		},
	})
}

func TestAccResourceObject_base64(t *testing.T) {
	t.Parallel()

	content := "testing123"
	base64EncodedContent := base64.StdEncoding.EncodeToString([]byte(content))
	bucketName := acctest.RandomWithPrefix("tf-test")
	keyName := acctest.RandomWithPrefix("tf_test")

	var object s3.GetObjectOutput

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Base64(t, bucketName, keyName, base64EncodedContent),
				Check: resource.ComposeTestCheckFunc(
					checkObjectExists(&object),
					checkObjectBodyContains(&object, content),
					resource.TestCheckResourceAttr(objectResourceName, "key", "test"),
				),
			},
		},
	})
}

func TestAccResourceObject_source(t *testing.T) {
	t.Parallel()

	content := "testing123"
	bucketName := acctest.RandomWithPrefix("tf-test")
	keyName := acctest.RandomWithPrefix("tf_test")

	file, err := ioutil.TempFile(os.TempDir(), "tf-test-obj-source")
	if err != nil {
		t.Fatalf("failed to create temp file: %s", err)
	}
	defer os.Remove(file.Name())

	if _, err := file.Write([]byte(content)); err != nil {
		t.Fatalf("failed to write to temp file: %s", err)
	}

	var object s3.GetObjectOutput

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Source(t, bucketName, keyName, file.Name()),
				Check: resource.ComposeTestCheckFunc(
					checkObjectExists(&object),
					checkObjectBodyContains(&object, content),
					resource.TestCheckResourceAttr(objectResourceName, "key", "test"),
				),
			},
		},
	})
}

func TestAccResourceObject_contentUpdate(t *testing.T) {
	t.Parallel()

	content := "testing123"
	bucketName := acctest.RandomWithPrefix("tf-test")
	keyName := acctest.RandomWithPrefix("tf_test")

	var object s3.GetObjectOutput

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, bucketName, keyName, content),
				Check: resource.ComposeTestCheckFunc(
					checkObjectExists(&object),
					checkObjectBodyContains(&object, content),
				),
			},
			{
				PreConfig: func() {
					content = "updated456"
				},
				Config: tmpl.Basic(t, bucketName, keyName, content),
				Check: resource.ComposeTestCheckFunc(
					checkObjectExists(&object),
					checkObjectBodyContains(&object, content),
				),
			},
		},
	})
}

func TestAccResourceObject_updates(t *testing.T) {
	t.Parallel()

	content := "testing123"
	bucketName := acctest.RandomWithPrefix("tf-test")
	keyName := acctest.RandomWithPrefix("tf_test")

	var object s3.GetObjectOutput

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, bucketName, keyName, content),
				Check: resource.ComposeTestCheckFunc(
					checkObjectExists(&object),
					checkObjectBodyContains(&object, content),
					resource.TestCheckResourceAttr(objectResourceName, "key", "test"),
					resource.TestCheckResourceAttr(objectResourceName, "acl", "private"),
					resource.TestCheckResourceAttr(objectResourceName, "force_destroy", "false"),
					resource.TestCheckNoResourceAttr(objectResourceName, "metadata"),
				),
			},
			{
				Config: tmpl.Updates(t, bucketName, keyName, content),
				Check: resource.ComposeTestCheckFunc(
					checkObjectExists(&object),
					checkObjectBodyContains(&object, content),
					resource.TestCheckResourceAttr(objectResourceName, "key", "test"),
					resource.TestCheckResourceAttr(objectResourceName, "acl", "public-read"),
					resource.TestCheckResourceAttr(objectResourceName, "force_destroy", "true"),
					resource.TestCheckResourceAttr(objectResourceName, "cache_control", "max-age=2592000"),
					resource.TestCheckResourceAttr(objectResourceName, "content_disposition", "attachment"),
					resource.TestCheckResourceAttr(objectResourceName, "content_type", "text/plain"),
					resource.TestCheckResourceAttr(objectResourceName, "content_encoding", "utf8"),
					resource.TestCheckResourceAttr(objectResourceName, "content_language", "en"),
					resource.TestCheckResourceAttr(objectResourceName, "website_redirect", "test.com"),
					resource.TestCheckResourceAttr(objectResourceName, "metadata.%", "2"),
					resource.TestCheckResourceAttr(objectResourceName, "metadata.foo", "bar"),
					resource.TestCheckResourceAttr(objectResourceName, "metadata.bar", "foo"),
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

func checkObjectExists(obj *s3.GetObjectOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[objectResourceName]
		if !ok {
			return fmt.Errorf("could not find resource %s in root module", objectResourceName)
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
