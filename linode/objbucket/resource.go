package objbucket

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func resourceLifecycleExpiration() *schema.Resource {
	return &schema.Resource{
		Schema: resourceSchemaExpiration,
	}
}

func resourceLifecycleNoncurrentExp() *schema.Resource {
	return &schema.Resource{
		Schema: resourceSchemaNonCurrentExp,
	}
}

func resourceLifeCycle() *schema.Resource {
	return &schema.Resource{
		Schema: resourceSchemaLifeCycle,
	}
}

func Resource() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceSchema,
		ReadContext:   readResource,
		CreateContext: createResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func readResource(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	cluster, label, err := DecodeBucketID(d.Id())
	if err != nil {
		return diag.Errorf("failed to parse Linode ObjectStorageBucket id %s", d.Id())
	}

	bucket, err := client.GetObjectStorageBucket(ctx, cluster, label)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing Object Storage Bucket %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to find the specified Linode ObjectStorageBucket: %s", err)
	}

	access, err := client.GetObjectStorageBucketAccess(ctx, cluster, label)
	if err != nil {
		return diag.Errorf("failed to find the access config for the specified Linode ObjectStorageBucket: %s", err)
	}

	// Functionality requiring direct S3 API access
	accessKey := d.Get("access_key").(string)
	secretKey := d.Get("secret_key").(string)

	_, versioningPresent := d.GetOk("versioning")
	_, lifecyclePresent := d.GetOk("lifecycle_rule")

	if versioningPresent || lifecyclePresent {
		if accessKey == "" || secretKey == "" {
			return diag.Errorf("access_key and secret_key are required to get versioning and lifecycle info")
		}

		conn := helper.S3ConnFromResourceData(d)

		if err := readBucketLifecycle(d, conn); err != nil {
			return diag.Errorf("failed to find get object storage bucket lifecycle: %s", err)
		}

		if err := readBucketVersioning(d, conn); err != nil {
			return diag.Errorf("failed to find get object storage bucket versioning: %s", err)
		}
	}

	d.SetId(fmt.Sprintf("%s:%s", bucket.Cluster, bucket.Label))
	d.Set("cluster", bucket.Cluster)
	d.Set("label", bucket.Label)
	d.Set("acl", access.ACL)
	d.Set("cors_enabled", access.CorsEnabled)

	return nil
}

func createResource(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	cluster := d.Get("cluster").(string)
	label := d.Get("label").(string)
	acl := d.Get("acl").(string)
	corsEnabled := d.Get("cors_enabled").(bool)

	createOpts := linodego.ObjectStorageBucketCreateOptions{
		Cluster:     cluster,
		Label:       label,
		ACL:         linodego.ObjectStorageACL(acl),
		CorsEnabled: &corsEnabled,
	}

	bucket, err := client.CreateObjectStorageBucket(ctx, createOpts)
	if err != nil {
		return diag.Errorf("failed to create a Linode ObjectStorageBucket: %s", err)
	}

	d.SetId(fmt.Sprintf("%s:%s", bucket.Cluster, bucket.Label))

	return updateResource(ctx, d, meta)
}

func updateResource(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	accessKey := d.Get("access_key")
	secretKey := d.Get("secret_key")

	conn := helper.S3ConnFromResourceData(d)

	if d.HasChanges("acl", "cors_enabled") {
		if err := updateBucketAccess(ctx, d, client); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("cert") {
		if err := updateBucketCert(ctx, d, client); err != nil {
			return diag.FromErr(err)
		}
	}

	versioningChanged := d.HasChange("versioning")
	lifecycleChanged := d.HasChange("lifecycle_rule")

	if versioningChanged || lifecycleChanged {
		if accessKey == "" || secretKey == "" {
			return diag.Errorf("access_key and secret_key are required to set versioning and lifecycle info")
		}

		// Ensure we only update what is changed
		if versioningChanged {
			if err := updateBucketVersioning(d, conn); err != nil {
				return diag.FromErr(err)
			}
		}

		if lifecycleChanged {
			if err := updateBucketLifecycle(d, conn); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return readResource(ctx, d, meta)
}

func deleteResource(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	cluster, label, err := DecodeBucketID(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing Linode ObjectStorageBucket id %s", d.Id())
	}
	err = client.DeleteObjectStorageBucket(ctx, cluster, label)
	if err != nil {
		return diag.Errorf("Error deleting Linode ObjectStorageBucket %s: %s", d.Id(), err)
	}
	return nil
}

func readBucketVersioning(d *schema.ResourceData, conn *s3.S3) error {
	label := d.Get("label").(string)

	versioningOutput, err := conn.GetBucketVersioning(&s3.GetBucketVersioningInput{
		Bucket: &label,
	})
	if err != nil {
		return fmt.Errorf("failed to get versioning for bucket id %s: %s", d.Id(), err)
	}

	d.Set("versioning", versioningOutput.Status != nil &&
		*versioningOutput.Status == s3.BucketVersioningStatusEnabled)

	return nil
}

func readBucketLifecycle(d *schema.ResourceData, conn *s3.S3) error {
	label := d.Get("label").(string)

	lifecycleConfigOutput, err := conn.GetBucketLifecycleConfiguration(
		&s3.GetBucketLifecycleConfigurationInput{Bucket: &label})

	// A "NoSuchLifecycleConfiguration" error should be ignored in this context
	if err != nil {
		if err, ok := err.(awserr.Error); !ok || (ok && err.Code() != "NoSuchLifecycleConfiguration") {
			return fmt.Errorf("failed to get lifecycle for bucket id %s: %s", d.Id(), err)
		}
	}

	d.Set("lifecycle_rule", flattenLifecycleRules(lifecycleConfigOutput.Rules))

	return nil
}

func updateBucketVersioning(d *schema.ResourceData, conn *s3.S3) error {
	bucket := d.Get("label").(string)
	n := d.Get("versioning").(bool)

	status := s3.BucketVersioningStatusSuspended
	if n {
		status = s3.BucketVersioningStatusEnabled
	}

	inputVersioningConfig := &s3.PutBucketVersioningInput{
		Bucket: &bucket,
		VersioningConfiguration: &s3.VersioningConfiguration{
			Status: &status,
		},
	}

	if _, err := conn.PutBucketVersioning(inputVersioningConfig); err != nil {
		return err
	}

	return nil
}

func updateBucketLifecycle(d *schema.ResourceData, conn *s3.S3) error {
	bucket := d.Get("label").(string)

	rules, err := expandLifecycleRules(d.Get("lifecycle_rule").([]interface{}))
	if err != nil {
		return err
	}

	_, err = conn.PutBucketLifecycleConfiguration(
		&s3.PutBucketLifecycleConfigurationInput{
			Bucket: &bucket,
			LifecycleConfiguration: &s3.BucketLifecycleConfiguration{
				Rules: rules,
			},
		})

	return err
}

func updateBucketAccess(
	ctx context.Context, d *schema.ResourceData, client linodego.Client) error {
	cluster := d.Get("cluster").(string)
	label := d.Get("label").(string)

	updateOpts := linodego.ObjectStorageBucketUpdateAccessOptions{}
	if d.HasChange("acl") {
		updateOpts.ACL = linodego.ObjectStorageACL(d.Get("acl").(string))
	}

	if d.HasChange("cors_enabled") {
		newCorsBool := d.Get("cors_enabled").(bool)
		updateOpts.CorsEnabled = &newCorsBool
	}

	if err := client.UpdateObjectStorageBucketAccess(ctx, cluster, label, updateOpts); err != nil {
		return fmt.Errorf("failed to update bucket access: %s", err)
	}

	return nil
}

func updateBucketCert(
	ctx context.Context, d *schema.ResourceData, client linodego.Client) error {
	cluster := d.Get("cluster").(string)
	label := d.Get("label").(string)
	oldCert, newCert := d.GetChange("cert")
	hasOldCert := len(oldCert.([]interface{})) != 0

	if hasOldCert {
		if err := client.DeleteObjectStorageBucketCert(ctx, cluster, label); err != nil {
			return fmt.Errorf("failed to delete old bucket cert: %s", err)
		}
	}

	certSpec := newCert.([]interface{})
	if len(certSpec) == 0 {
		return nil
	}

	uploadOptions := expandBucketCert(certSpec[0])
	if _, err := client.UploadObjectStorageBucketCert(ctx, cluster, label, uploadOptions); err != nil {
		return fmt.Errorf("failed to upload new bucket cert: %s", err)
	}
	return nil
}

func expandBucketCert(v interface{}) linodego.ObjectStorageBucketCertUploadOptions {
	certSpec := v.(map[string]interface{})
	return linodego.ObjectStorageBucketCertUploadOptions{
		Certificate: certSpec["certificate"].(string),
		PrivateKey:  certSpec["private_key"].(string),
	}
}

func DecodeBucketID(id string) (cluster, label string, err error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		err = fmt.Errorf("Linode Object Storage Bucket ID must be of the form <Cluster>:<Label>, was provided: %s", id)
		return
	}
	cluster = parts[0]
	label = parts[1]
	return
}

func flattenLifecycleRules(rules []*s3.LifecycleRule) []map[string]interface{} {
	result := make([]map[string]interface{}, len(rules))

	for i, rule := range rules {
		ruleMap := make(map[string]interface{})

		if id := rule.ID; id != nil {
			ruleMap["id"] = *id
		}

		if prefix := rule.Prefix; prefix != nil {
			ruleMap["prefix"] = *prefix
		}

		if status := rule.Status; status != nil {
			ruleMap["enabled"] = *status == "Enabled"
		}

		if rule.AbortIncompleteMultipartUpload != nil {
			ruleMap["abort_incomplete_multipart_upload_days"] = *rule.AbortIncompleteMultipartUpload.DaysAfterInitiation
		}

		if rule.Expiration != nil {
			e := make(map[string]interface{})

			if date := rule.Expiration.Date; date != nil {
				e["date"] = rule.Expiration.Date.Format("2006-01-02")
			}

			if days := rule.Expiration.Days; days != nil {
				e["days"] = *days
			}

			if marker := rule.Expiration.ExpiredObjectDeleteMarker; marker != nil && *marker {
				e["expired_object_delete_marker"] = *marker
			}

			ruleMap["expiration"] = []interface{}{e}
		}

		if rule.NoncurrentVersionExpiration != nil {
			e := make(map[string]interface{})

			if days := rule.NoncurrentVersionExpiration.NoncurrentDays; days != nil && *days > 0 {
				e["days"] = *days
			}

			ruleMap["noncurrent_version_expiration"] = []interface{}{e}
		}

		result[i] = ruleMap
	}

	return result
}

func expandLifecycleRules(ruleSpecs []interface{}) ([]*s3.LifecycleRule, error) {
	rules := make([]*s3.LifecycleRule, len(ruleSpecs))
	for i, ruleSpec := range ruleSpecs {
		ruleSpec := ruleSpec.(map[string]interface{})
		rule := &s3.LifecycleRule{}

		status := "Disabled"
		if ruleSpec["enabled"].(bool) {
			status = "Enabled"
		}
		rule.Status = &status

		if id, ok := ruleSpec["id"]; ok {
			id := id.(string)
			rule.ID = &id
		}

		if prefix, ok := ruleSpec["prefix"]; ok {
			prefix := prefix.(string)
			rule.Prefix = &prefix
		}

		//nolint:lll
		if abortIncompleteDays, ok := ruleSpec["abort_incomplete_multipart_upload_days"].(int); ok && abortIncompleteDays > 0 {
			rule.AbortIncompleteMultipartUpload = &s3.AbortIncompleteMultipartUpload{}
			abortIncompleteDays := int64(abortIncompleteDays)

			rule.AbortIncompleteMultipartUpload.DaysAfterInitiation = &abortIncompleteDays
		}

		if expirationList := ruleSpec["expiration"].([]interface{}); len(expirationList) > 0 {
			rule.Expiration = &s3.LifecycleExpiration{}

			expirationMap := expirationList[0].(map[string]interface{})

			if dateStr, ok := expirationMap["date"].(string); ok && dateStr != "" {
				date, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", dateStr))
				if err != nil {
					return nil, err
				}

				rule.Expiration.Date = &date
			}

			if days, ok := expirationMap["days"].(int); ok && days > 0 {
				days := int64(days)

				rule.Expiration.Days = &days
			}

			if marker, ok := expirationMap["expired_object_delete_marker"].(bool); ok && marker {
				rule.Expiration.ExpiredObjectDeleteMarker = &marker
			}
		}

		if expirationList := ruleSpec["noncurrent_version_expiration"].([]interface{}); len(expirationList) > 0 {
			rule.NoncurrentVersionExpiration = &s3.NoncurrentVersionExpiration{}

			expirationMap := expirationList[0].(map[string]interface{})

			if days, ok := expirationMap["days"]; ok {
				days := int64(days.(int))
				rule.NoncurrentVersionExpiration.NoncurrentDays = &days
			}
		}

		rules[i] = rule
	}

	return rules, nil
}
