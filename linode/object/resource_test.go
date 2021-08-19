package object_test

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
)

const objectResourceName = "linode_object_storage_object.object"

func checkObjectExists(obj *s3.GetObjectOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[objectResourceName]
		if !ok {
			return fmt.Errorf("could not find resource %s in root module", objectResourceName)
		}

		bucket := rs.Primary.Attributes["bucket"]
		key := rs.Primary.Attributes["key"]
		etag := rs.Primary.Attributes["etag"]
		accessKey := rs.Primary.Attributes["access_key"]
		secretKey := rs.Primary.Attributes["secret_key"]
		cluster := rs.Primary.Attributes["cluster"]

		conn := s3.New(session.New(&aws.Config{
			Region:      aws.String("us-east-1"),
			Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
			Endpoint:    aws.String(fmt.Sprintf(linodeObjectsEndpoint, cluster)),
		}))

		out, err := conn.GetObject(
			&s3.GetObjectInput{
				Bucket:  &bucket,
				Key:     &key,
				IfMatch: &etag,
			})
		if err != nil {
			return fmt.Errorf("failed to get Bucket (%s) Object (%s): %s", bucket, key, err)
		}

		*obj = *out
		return nil
	}
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

func TestAccResourceObject_basic(t *testing.T) {
	t.Parallel()

	content := "testing123"
	bucketName := acctest.RandomWithPrefix("tf-test")
	keyName := acctest.RandomWithPrefix("tf_test")

	var object s3.GetObjectOutput

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckLinodeObjectStorageKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasic(bucketName, keyName, content),
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
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckLinodeObjectStorageKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBase64(bucketName, keyName, base64EncodedContent),
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
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckLinodeObjectStorageKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigSource(bucketName, keyName, file.Name()),
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
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckLinodeObjectStorageKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasic(bucketName, keyName, content),
				Check: resource.ComposeTestCheckFunc(
					checkObjectExists(&object),
					checkObjectBodyContains(&object, content),
				),
			},
			{
				PreConfig: func() {
					content = "updated456"
				},
				Config: resourceConfigBasic(bucketName, keyName, content),
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
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckLinodeObjectStorageKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasic(bucketName, keyName, content),
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
				Config: resourceConfigUpdates(bucketName, keyName, content),
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

func resourceConfigBasic(name, keyName, content string) string {
	return testAccCheckLinodeObjectStorageBucketConfigBasic(name) + testAccCheckLinodeObjectStorageKeyConfigBasic(keyName) + fmt.Sprintf(`
resource "linode_object_storage_object" "object" {
	bucket     = linode_object_storage_bucket.foobar.label
	cluster    = "us-east-1"
	access_key = linode_object_storage_key.foobar.access_key
	secret_key = linode_object_storage_key.foobar.secret_key
	key        = "test"
	content    = "%s"
}`, content)
}

func resourceConfigBase64(name, keyName, content string) string {
	return testAccCheckLinodeObjectStorageBucketConfigBasic(name) + testAccCheckLinodeObjectStorageKeyConfigBasic(keyName) + fmt.Sprintf(`
resource "linode_object_storage_object" "object" {
	bucket         = linode_object_storage_bucket.foobar.label
	cluster        = "us-east-1"
	access_key     = linode_object_storage_key.foobar.access_key
	secret_key     = linode_object_storage_key.foobar.secret_key
	key            = "test"
	content_base64 = "%s"
}`, content)
}

func resourceConfigSource(name, keyName, filePath string) string {
	return testAccCheckLinodeObjectStorageBucketConfigBasic(name) + testAccCheckLinodeObjectStorageKeyConfigBasic(keyName) + fmt.Sprintf(`
resource "linode_object_storage_object" "object" {
	bucket     = linode_object_storage_bucket.foobar.label
	cluster    = "us-east-1"
	access_key = linode_object_storage_key.foobar.access_key
	secret_key = linode_object_storage_key.foobar.secret_key
	key        = "test"
	source     = "%s"
}`, filePath)
}

func resourceConfigUpdates(name, keyName, content string) string {
	return testAccCheckLinodeObjectStorageBucketConfigBasic(name) + testAccCheckLinodeObjectStorageKeyConfigBasic(keyName) + fmt.Sprintf(`
	resource "linode_object_storage_object" "object" {
		bucket     = linode_object_storage_bucket.foobar.label
		cluster    = "us-east-1"
		access_key = linode_object_storage_key.foobar.access_key
		secret_key = linode_object_storage_key.foobar.secret_key
		key        = "test"
		content    = "%s"
		acl        = "public-read"

		content_type     = "text/plain"
		content_encoding = "utf8"
		content_language = "en"
		website_redirect = "test.com"
		force_destroy    = true

		content_disposition = "attachment"
		cache_control       = "max-age=2592000"

		metadata = {
			foo = "bar"
			bar = "foo"
		}
	}`, content)
}
