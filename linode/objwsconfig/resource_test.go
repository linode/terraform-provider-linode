//go:build integration

package objwsconfig_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/objwsconfig/tmpl"
)

var testCluster string

func init() {
	cluster, err := acceptance.GetRandomOBJCluster()
	if err != nil {
		log.Fatal(err)
	}

	testCluster = cluster
}

func TestAccResourceObjectWebsiteConfig_basic(t *testing.T) {
	t.Parallel()

	testBucket := acctest.RandomWithPrefix("tf-test")
	testKeyName := acctest.RandomWithPrefix("tf_test")
	resName := "linode_object_storage_website_config.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, testCluster, testBucket, testKeyName),
				Check: resource.ComposeTestCheckFunc(
					checkBucketWebsiteConfigExists,
					checkBucketWebsiteEndpointAccessible,
					resource.TestCheckResourceAttr(resName, "cluster", testCluster),
					resource.TestCheckResourceAttr(resName, "bucket", testBucket),
					resource.TestCheckResourceAttrSet(resName, "website_endpoint"),
				),
			},
			{
				Config: tmpl.BasicDependency(t, testCluster, testBucket, testKeyName), // only destroy the website config
				Check:  checkBucketWebsiteConfigDestroy,
			},
		},
	})
}

func TestAccResourceObjectWebsiteConfig_update(t *testing.T) {
	t.Parallel()

	testBucket := acctest.RandomWithPrefix("tf-test")
	testKeyName := acctest.RandomWithPrefix("tf_test")
	resNameUpdate := "linode_object_storage_website_config.update"
	resNameReplace := "linode_object_storage_website_config.replace"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.UpdatesBefore(t, testCluster, testBucket, testKeyName),
				Check: resource.ComposeTestCheckFunc(
					checkBucketWebsiteConfigExists,
					checkBucketWebsiteEndpointAccessible,
					resource.TestCheckResourceAttr(resNameUpdate, "bucket", testBucket),
					resource.TestCheckResourceAttr(resNameUpdate, "index_document", "index.html"),
					resource.TestCheckNoResourceAttr(resNameUpdate, "error_document"),
					resource.TestCheckResourceAttrSet(resNameUpdate, "website_endpoint"),
					resource.TestCheckResourceAttr(resNameReplace, "bucket", fmt.Sprintf("%s-2", testBucket)),
					resource.TestCheckResourceAttrSet(resNameReplace, "website_endpoint"),
				),
			},
			{
				Config: tmpl.UpdatesAfter(t, testCluster, testBucket, testKeyName),
				Check: resource.ComposeTestCheckFunc(
					checkBucketWebsiteConfigExists,
					checkBucketWebsiteEndpointAccessible,
					resource.TestCheckResourceAttr(resNameUpdate, "bucket", testBucket),
					resource.TestCheckResourceAttr(resNameUpdate, "index_document", "sub/index.html"),
					resource.TestCheckResourceAttr(resNameUpdate, "error_document", "404.html"),
					resource.TestCheckResourceAttrSet(resNameUpdate, "website_endpoint"),
					resource.TestCheckResourceAttr(resNameReplace, "bucket", fmt.Sprintf("%s-3", testBucket)),
					resource.TestCheckResourceAttr(resNameReplace, "index_document", "index.html"),
					resource.TestCheckResourceAttrSet(resNameReplace, "website_endpoint"),
				),
			},
			{
				Config: tmpl.UpdateDependency(t, testCluster, testBucket, testKeyName),
				Check:  checkBucketWebsiteConfigDestroy,
			},
		},
	})
}

func findBucketWebsiteConfigResource(s *terraform.State) []*terraform.ResourceState {
	configs := []*terraform.ResourceState{}
	for _, res := range s.RootModule().Resources {
		if res.Type != "linode_object_storage_website_config" {
			continue
		}
		configs = append(configs, res)
	}
	return configs
}

func checkBucketWebsiteEndpointAccessible(s *terraform.State) error {
	configs := findBucketWebsiteConfigResource(s)
	secret := configs[0].Primary.Attributes["website_endpoint"]

	if secret == "" {
		return errors.New("Expected website_endpoint to be accessible but got empty")
	}
	return nil
}

func checkBucketWebsiteConfigExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_object_storage_website_config" {
			continue
		}

		_, err := getBucketWebsiteConfig(context.Background(), rs)
		if err != nil {
			return fmt.Errorf("Error retrieving state of bucket website config %s: %s", rs.Primary.Attributes["bucket"], err)
		}
	}

	return nil
}

func checkBucketWebsiteConfigDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_object_storage_website_config" {
			continue
		}

		_, err := getBucketWebsiteConfig(context.Background(), rs)
		if err == nil {
			return fmt.Errorf("Bucket website config still exists %s", rs.Primary.Attributes["bucket"])
		}

		var re *awshttp.ResponseError
		if !errors.As(err, &re) || re.HTTPStatusCode() != 404 {
			return fmt.Errorf("Error requesting bucket website config %s", rs.Primary.Attributes["bucket"])
		}
	}

	return nil
}

func getBucketWebsiteConfig(ctx context.Context, rs *terraform.ResourceState) (*s3.GetBucketWebsiteOutput, error) {
	cluster := rs.Primary.Attributes["cluster"]
	bucket := rs.Primary.Attributes["bucket"]
	accessKey := rs.Primary.Attributes["access_key"]
	secretKey := rs.Primary.Attributes["secret_key"]

	linodeClient := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
	b, err := linodeClient.GetObjectStorageBucket(ctx, cluster, bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to find the specified Linode ObjectStorageBucket: %s", err)
	}
	endpoint := helper.ComputeS3EndpointFromBucket(ctx, *b)

	s3client, err := helper.S3Connection(ctx, endpoint, accessKey, secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get create s3 client: %w", err)
	}

	return s3client.GetBucketWebsite(ctx, &s3.GetBucketWebsiteInput{Bucket: &bucket})
}
