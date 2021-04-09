package linode

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	linodeObjectsEndpoint = "https://%s.linodeobjects.com"
)

// s3ConnFromResourceData builds an S3 client from the linode_object_storage_object
// resource's access_key, secret_key, and cluster fields.
func s3ConnFromResourceData(d *schema.ResourceData) *s3.S3 {
	accessKey := d.Get("access_key").(string)
	secretKey := d.Get("secret_key").(string)
	cluster := d.Get("cluster").(string)

	sess := session.New(&aws.Config{
		// This region is hardcoded strictly for preflight validation purposes.
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:    aws.String(fmt.Sprintf(linodeObjectsEndpoint, cluster)),
	})
	return s3.New(sess)
}

func buildObjectStorageObjectID(d *schema.ResourceData) string {
	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	return fmt.Sprintf("%s/%s", bucket, key)
}
