//go:build integration || objbucket

package objbucket_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/objbucket"
	"github.com/linode/terraform-provider-linode/v2/linode/objbucket/tmpl"
)

const (
	objAccessKeyEnvVar = "LINODE_OBJ_ACCESS_KEY"
	objSecretKeyEnvVar = "LINODE_OBJ_SECRET_KEY"
)

var testCluster string

func init() {
	cluster, err := acceptance.GetRandomOBJCluster()
	if err != nil {
		log.Fatal(err)
	}

	testCluster = cluster
}

func init() {
	resource.AddTestSweepers("linode_object_storage_bucket", &resource.Sweeper{
		Name: "linode_object_storage_bucket",
		F:    sweep,
	})
}

func generateTestCert(domain string) (certificate, privateKey string, err error) {
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate key: %s", err)
	}
	keyUsage := x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment

	validFrom := time.Now()
	validUntil := validFrom.Add(time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate serial number: %s", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Linode"},
		},
		NotBefore:             validFrom,
		NotAfter:              validUntil,
		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{domain},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, priv.Public(), priv)
	if err != nil {
		return "", "", fmt.Errorf("failed to create certificate: %s", err)
	}
	certBuffer := new(bytes.Buffer)
	if err := pem.Encode(certBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return "", "", fmt.Errorf("failed to encode certificate to PEM: %s", err)
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal private key: %s", err)
	}
	keyBuffer := new(bytes.Buffer)
	if err := pem.Encode(keyBuffer, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		return "", "", fmt.Errorf("failed to encode private key to PEM: %s", err)
	}

	return string(certBuffer.Bytes()), string(keyBuffer.Bytes()), nil
}

func sweep(prefix string) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	listOpts := acceptance.SweeperListOptions(prefix, "label")
	objectStorageBuckets, err := client.ListObjectStorageBuckets(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("Error getting object_storage_buckets: %s", err)
	}

	accessKey, accessKeyOk := os.LookupEnv(objAccessKeyEnvVar)
	secretKey, secretKeyOk := os.LookupEnv(objSecretKeyEnvVar)
	haveBucketAccess := accessKeyOk && secretKeyOk

	for _, objectStorageBucket := range objectStorageBuckets {
		if !acceptance.ShouldSweep(prefix, objectStorageBucket.Label) {
			continue
		}
		bucket := objectStorageBucket.Label
		s3client, err := helper.S3Connection(
			context.Background(),
			helper.ComputeS3EndpointFromBucket(context.Background(), objectStorageBucket),
			accessKey,
			secretKey,
		)
		if err != nil {
			tflog.Error(context.Background(), fmt.Sprintf("failed to create s3 client: %v", err))
		}

		helper.PurgeAllObjects(context.Background(), bucket, s3client, true, true)

		_, err = s3client.DeleteBucket(
			context.Background(),
			&s3.DeleteBucketInput{
				Bucket: aws.String(bucket),
			},
		)
		if err != nil {
			if apiErr, ok := err.(*linodego.Error); ok && !haveBucketAccess && strings.HasPrefix(
				apiErr.Message, fmt.Sprintf("Bucket %s is not empty", bucket)) {
				tflog.Warn(
					context.Background(),
					fmt.Sprintf(
						"will not delete Object Storage Bucket (%s) as it needs to be emptied; "+
							"specify %q and %q env variables for bucket access",
						bucket, objAccessKeyEnvVar, objSecretKeyEnvVar,
					),
				)
				continue
			}
			return fmt.Errorf("Error destroying %s during sweep: %s", bucket, err)
		}
	}

	return nil
}

func TestAccResourceBucket_basic_smoke(t *testing.T) {
	t.Parallel()

	acceptance.RunTestRetry(t, 5, func(retryT *acceptance.TRetry) {
		resName := "linode_object_storage_bucket.foobar"
		objectStorageBucketName := acctest.RandomWithPrefix("tf-test")

		resource.Test(retryT, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
			CheckDestroy:             checkBucketDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.Basic(t, objectStorageBucketName, testCluster),
					Check: resource.ComposeTestCheckFunc(
						checkBucketExists,
						resource.TestCheckResourceAttr(resName, "label", objectStorageBucketName),
						resource.TestCheckResourceAttrSet(resName, "hostname"),
					),
				},
				{
					ResourceName:      resName,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})
}

func TestAccResourceBucket_access(t *testing.T) {
	t.Parallel()

	acceptance.RunTestRetry(t, 5, func(retryT *acceptance.TRetry) {
		resName := "linode_object_storage_bucket.foobar"
		objectStorageBucketName := acctest.RandomWithPrefix("tf-test")

		resource.Test(retryT, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
			CheckDestroy:             checkBucketDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.Access(t, objectStorageBucketName, testCluster, "public-read", true),
					Check: resource.ComposeTestCheckFunc(
						checkBucketExists,
						resource.TestCheckResourceAttr(resName, "label", objectStorageBucketName),
						resource.TestCheckResourceAttr(resName, "acl", "public-read"),
						resource.TestCheckResourceAttr(resName, "cors_enabled", "true"),
					),
				},
				{
					Config: tmpl.Access(t, objectStorageBucketName, testCluster, "private", false),
					Check: resource.ComposeTestCheckFunc(
						checkBucketExists,
						resource.TestCheckResourceAttr(resName, "label", objectStorageBucketName),
						resource.TestCheckResourceAttr(resName, "acl", "private"),
						resource.TestCheckResourceAttr(resName, "cors_enabled", "false"),
					),
				},
			},
		})
	})
}

func TestAccResourceBucket_versioning(t *testing.T) {
	t.Parallel()

	acceptance.RunTestRetry(t, 5, func(retryT *acceptance.TRetry) {
		resName := "linode_object_storage_bucket.foobar"
		objectStorageBucketName := acctest.RandomWithPrefix("tf-test")
		objectStorageKeyName := acctest.RandomWithPrefix("tf-test")

		resource.Test(retryT, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
			CheckDestroy:             checkBucketDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.Versioning(t, objectStorageBucketName, testCluster, objectStorageKeyName, true),
					Check: resource.ComposeTestCheckFunc(
						checkBucketExists,
						resource.TestCheckResourceAttr(resName, "label", objectStorageBucketName),
						resource.TestCheckResourceAttr(resName, "versioning", "true"),
					),
				},
				{
					Config: tmpl.Versioning(t, objectStorageBucketName, testCluster, objectStorageKeyName, false),
					Check: resource.ComposeTestCheckFunc(
						checkBucketExists,
						resource.TestCheckResourceAttr(resName, "label", objectStorageBucketName),
						resource.TestCheckResourceAttr(resName, "versioning", "false"),
					),
				},
			},
		})
	})
}

func TestAccResourceBucket_lifecycle(t *testing.T) {
	t.Parallel()

	acceptance.RunTestRetry(t, 5, func(retryT *acceptance.TRetry) {
		resName := "linode_object_storage_bucket.foobar"
		objectStorageBucketName := acctest.RandomWithPrefix("tf-test")
		objectStorageKeyName := acctest.RandomWithPrefix("tf-test")

		resource.Test(retryT, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
			CheckDestroy:             checkBucketDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.LifeCycle(t, objectStorageBucketName, testCluster, objectStorageKeyName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "label", objectStorageBucketName),
						resource.TestCheckResourceAttr(resName, "cluster", testCluster),
						resource.TestCheckResourceAttr(resName, "lifecycle_rule.#", "1"),
						resource.TestCheckResourceAttr(resName, "lifecycle_rule.0.id", "test-rule"),
						resource.TestCheckResourceAttr(resName, "lifecycle_rule.0.prefix", "tf"),
						resource.TestCheckResourceAttr(resName, "lifecycle_rule.0.enabled", "true"),
						resource.TestCheckResourceAttr(resName, "lifecycle_rule.0.expiration.#", "1"),
						resource.TestCheckResourceAttr(resName, "lifecycle_rule.0.abort_incomplete_multipart_upload_days", "5"),
						resource.TestCheckResourceAttrSet(resName, "lifecycle_rule.0.expiration.0.date"),
					),
				},
				{
					Config: tmpl.LifeCycleUpdates(t, objectStorageBucketName, testCluster, objectStorageKeyName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "label", objectStorageBucketName),
						resource.TestCheckResourceAttr(resName, "cluster", testCluster),
						resource.TestCheckResourceAttr(resName, "lifecycle_rule.#", "1"),
						resource.TestCheckResourceAttr(resName, "lifecycle_rule.0.id", "test-rule-update"),
						resource.TestCheckResourceAttr(resName, "lifecycle_rule.0.prefix", "tf-update"),
						resource.TestCheckResourceAttr(resName, "lifecycle_rule.0.enabled", "false"),
						resource.TestCheckResourceAttr(resName, "lifecycle_rule.0.abort_incomplete_multipart_upload_days", "42"),
						resource.TestCheckResourceAttr(resName, "lifecycle_rule.0.expiration.#", "1"),
						resource.TestCheckResourceAttr(resName, "lifecycle_rule.0.expiration.0.days", "37"),
					),
				},
				{
					Config: tmpl.LifeCycleRemoved(t, objectStorageBucketName, testCluster, objectStorageKeyName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "label", objectStorageBucketName),
						resource.TestCheckResourceAttr(resName, "cluster", testCluster),
						resource.TestCheckResourceAttr(resName, "lifecycle_rule.#", "0"),
					),
				},
			},
		})
	})
}

func TestAccResourceBucket_lifecycleNoID(t *testing.T) {
	t.Parallel()

	resName := "linode_object_storage_bucket.foobar"
	objectStorageBucketName := acctest.RandomWithPrefix("tf-test")
	objectStorageKeyName := acctest.RandomWithPrefix("tf-test")

	acceptance.RunTestRetry(t, 5, func(retryT *acceptance.TRetry) {
		resource.Test(retryT, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
			CheckDestroy:             checkBucketDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.LifeCycleNoID(t, objectStorageBucketName, testCluster, objectStorageKeyName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "label", objectStorageBucketName),
						resource.TestCheckResourceAttr(resName, "cluster", testCluster),
						resource.TestCheckResourceAttr(resName, "lifecycle_rule.#", "1"),
						resource.TestCheckResourceAttrSet(resName, "lifecycle_rule.0.id"),
						resource.TestCheckResourceAttr(resName, "lifecycle_rule.0.prefix", "tf"),
						resource.TestCheckResourceAttr(resName, "lifecycle_rule.0.enabled", "true"),
						resource.TestCheckResourceAttr(resName, "lifecycle_rule.0.expiration.#", "1"),
						resource.TestCheckResourceAttr(resName, "lifecycle_rule.0.abort_incomplete_multipart_upload_days", "5"),
						resource.TestCheckResourceAttrSet(resName, "lifecycle_rule.0.expiration.0.date"),
					),
				},
			},
		})
	})
}

func TestAccResourceBucket_cert(t *testing.T) {
	t.Parallel()

	resName := "linode_object_storage_bucket.foobar"
	objectStorageBucketName := acctest.RandomWithPrefix("tf-test") + ".io"
	cert, key, err := generateTestCert(objectStorageBucketName)
	if err != nil {
		t.Fatal(err)
	}

	invalidCert, invalidKey, err := generateTestCert("bogusdomain.com")
	if err != nil {
		t.Fatal(err)
	}

	otherCert, otherKey, err := generateTestCert(objectStorageBucketName)
	if err != nil {
		t.Fatal(err)
	}

	acceptance.RunTestRetry(t, 5, func(retryT *acceptance.TRetry) {
		resource.Test(retryT, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
			CheckDestroy:             checkBucketDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.Cert(t, objectStorageBucketName, testCluster, cert, key),
					Check: resource.ComposeTestCheckFunc(
						checkBucketExists,
						checkBucketHasSSL(true),
						resource.TestCheckResourceAttr(resName, "label", objectStorageBucketName),
					),
				},
				{
					Config:      tmpl.Cert(t, objectStorageBucketName, testCluster, invalidCert, invalidKey),
					ExpectError: regexp.MustCompile("failed to upload new bucket cert"),
				},
				{
					Config: tmpl.Cert(t, objectStorageBucketName, testCluster, otherCert, otherKey),
					Check: resource.ComposeTestCheckFunc(
						checkBucketExists,
						checkBucketHasSSL(true),
						resource.TestCheckResourceAttr(resName, "label", objectStorageBucketName),
					),
				},
				{
					Config: tmpl.Basic(t, objectStorageBucketName, testCluster),
					Check: resource.ComposeTestCheckFunc(
						checkBucketExists,
						checkBucketHasSSL(false),
						resource.TestCheckResourceAttr(resName, "label", objectStorageBucketName),
					),
				},
			},
		})
	})
}

func TestAccResourceBucket_dataSource(t *testing.T) {
	t.Parallel()

	resName := "linode_object_storage_bucket.foobar"
	objectStorageBucketName := acctest.RandomWithPrefix("tf-test")

	acceptance.RunTestRetry(t, 5, func(retryT *acceptance.TRetry) {
		resource.Test(retryT, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
			CheckDestroy:             checkBucketDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.ClusterDataBasic(t, objectStorageBucketName, testCluster),
					Check: resource.ComposeTestCheckFunc(
						checkBucketExists,
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
	})
}

func TestAccResourceBucket_update(t *testing.T) {
	t.Parallel()

	objectStorageBucketName := acctest.RandomWithPrefix("tf-test")
	resName := "linode_object_storage_bucket.foobar"

	acceptance.RunTestRetry(t, 5, func(retryT *acceptance.TRetry) {
		resource.Test(retryT, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
			CheckDestroy:             checkBucketDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.Basic(t, objectStorageBucketName, testCluster),
					Check: resource.ComposeTestCheckFunc(
						checkBucketExists,
						resource.TestCheckResourceAttr(resName, "label", objectStorageBucketName),
					),
				},
				{
					Config: tmpl.Updates(t, objectStorageBucketName, testCluster),
					Check: resource.ComposeTestCheckFunc(
						checkBucketExists,
						resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s-renamed", objectStorageBucketName)),
					),
				},
			},
		})
	})
}

func TestAccResourceBucket_credsConfiged(t *testing.T) {
	t.Parallel()

	resName := "linode_object_storage_bucket.foobar"
	objectStorageBucketName := acctest.RandomWithPrefix("tf-test")
	objectStorageKeyName := acctest.RandomWithPrefix("tf-test")

	acceptance.RunTestRetry(t, 5, func(retryT *acceptance.TRetry) {
		resource.Test(retryT, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
			CheckDestroy:             checkBucketDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.CredsConfiged(t, objectStorageBucketName, testCluster, objectStorageKeyName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "label", objectStorageBucketName),
						resource.TestCheckResourceAttr(resName, "cluster", testCluster),
						resource.TestCheckResourceAttr(resName, "lifecycle_rule.#", "1"),
					),
				},
			},
		})
	})
}

func TestAccResourceBucket_tempKeys(t *testing.T) {
	t.Parallel()

	resName := "linode_object_storage_bucket.foobar"
	objectStorageBucketName := acctest.RandomWithPrefix("tf-test")
	objectStorageKeyName := acctest.RandomWithPrefix("tf-test")

	acceptance.RunTestRetry(t, 5, func(retryT *acceptance.TRetry) {
		resource.Test(retryT, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
			CheckDestroy:             checkBucketDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.TempKeys(t, objectStorageBucketName, testCluster, objectStorageKeyName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "label", objectStorageBucketName),
						resource.TestCheckResourceAttr(resName, "cluster", testCluster),
						resource.TestCheckResourceAttr(resName, "lifecycle_rule.#", "1"),
					),
				},
			},
		})
	})
}

func checkBucketExists(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_object_storage_objbucket" {
			continue
		}

		cluster, label, err := objbucket.DecodeBucketID(context.Background(), rs.Primary.ID)
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

func checkBucketHasSSL(expected bool) func(*terraform.State) error {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "linode_object_storage_bucket" {
				continue
			}

			cluster, label, err := objbucket.DecodeBucketID(context.Background(), rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("could not parse bucket ID %s: %s", rs.Primary.ID, err)
			}

			cert, err := client.GetObjectStorageBucketCert(context.TODO(), cluster, label)
			if err != nil {
				return fmt.Errorf("failed to get bucket cert: %s", err)
			}

			if cert.SSL != expected {
				return fmt.Errorf("expected cert.SSL to be %v; got %v", expected, cert.SSL)
			}
		}
		return nil
	}
}

func checkBucketDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_object_storage_bucket" {
			continue
		}

		id := rs.Primary.ID
		cluster, label, err := objbucket.DecodeBucketID(context.Background(), id)
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

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 && apiErr.Code != 500 {
			return fmt.Errorf("Error requesting Linode ObjectStorageBucket with id %s: %s", id, err)
		}
	}

	return nil
}
