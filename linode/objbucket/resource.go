package objbucket

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	ctx context.Context, d *schema.ResourceData, meta interface{},
) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	cluster, label, err := DecodeBucketID(d.Id())
	if err != nil {
		return diag.Errorf("failed to parse Linode ObjectStorageBucket id %s", d.Id())
	}

	bucket, err := client.GetObjectStorageBucket(ctx, cluster, label)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			tflog.Warn(
				ctx,
				fmt.Sprintf(
					"[WARN] removing Object Storage Bucket %q from state because it no longer exists",
					d.Id(),
				),
			)
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
	endpoint := helper.ComputeS3EndpointFromBucket(*bucket)

	_, versioningPresent := d.GetOk("versioning")
	_, lifecyclePresent := d.GetOk("lifecycle_rule")

	if versioningPresent || lifecyclePresent {
		if accessKey == "" || secretKey == "" {
			return diag.Errorf("access_key and secret_key are required to get versioning and lifecycle info")
		}

		s3Client, err := helper.S3ConnectionV2(endpoint, accessKey, secretKey)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := readBucketLifecycle(ctx, d, s3Client); err != nil {
			return diag.Errorf("failed to find get object storage bucket lifecycle: %s", err)
		}

		if err := readBucketVersioning(ctx, d, s3Client); err != nil {
			return diag.Errorf("failed to find get object storage bucket versioning: %s", err)
		}
	}

	d.SetId(fmt.Sprintf("%s:%s", bucket.Cluster, bucket.Label))
	d.Set("cluster", bucket.Cluster)
	d.Set("label", bucket.Label)
	d.Set("hostname", bucket.Hostname)
	d.Set("acl", access.ACL)
	d.Set("cors_enabled", access.CorsEnabled)
	d.Set("endpoint", endpoint)

	return nil
}

func createResource(
	ctx context.Context, d *schema.ResourceData, meta interface{},
) diag.Diagnostics {
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

	d.Set("endpoint", helper.ComputeS3EndpointFromBucket(*bucket))
	d.SetId(fmt.Sprintf("%s:%s", bucket.Cluster, bucket.Label))

	return updateResource(ctx, d, meta)
}

func updateResource(
	ctx context.Context, d *schema.ResourceData, meta interface{},
) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

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
		s3client, err := helper.S3ConnectionFromDataV2(ctx, d, meta)
		if err != nil {
			return diag.FromErr(err)
		}

		// Ensure we only update what is changed
		if versioningChanged {
			if err := updateBucketVersioning(ctx, d, s3client); err != nil {
				return diag.FromErr(err)
			}
		}

		if lifecycleChanged {
			if err := updateBucketLifecycle(ctx, d, s3client); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return readResource(ctx, d, meta)
}

func deleteResource(
	ctx context.Context, d *schema.ResourceData, meta interface{},
) diag.Diagnostics {
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

func readBucketVersioning(ctx context.Context, d *schema.ResourceData, client *s3.Client) error {
	label := d.Get("label").(string)

	versioningOutput, err := client.GetBucketVersioning(
		ctx,
		&s3.GetBucketVersioningInput{Bucket: &label},
	)
	if err != nil {
		return fmt.Errorf("failed to get versioning for bucket id %s: %s", d.Id(), err)
	}

	d.Set("versioning", versioningOutput.Status == s3types.BucketVersioningStatusEnabled)

	return nil
}

func readBucketLifecycle(ctx context.Context, d *schema.ResourceData, client *s3.Client) error {
	label := d.Get("label").(string)

	lifecycleConfigOutput, err := client.GetBucketLifecycleConfiguration(
		ctx,
		&s3.GetBucketLifecycleConfigurationInput{Bucket: &label},
	)
	// A "NoSuchLifecycleConfiguration" error should be ignored in this context
	if err != nil {
		var ae smithy.APIError
		if ok := errors.As(err, &ae); !ok || ae.ErrorCode() != "NoSuchLifecycleConfiguration" {
			return fmt.Errorf("failed to get lifecycle for bucket id %s: %w", d.Id(), err)
		}
	}

	if lifecycleConfigOutput == nil {
		tflog.Debug(ctx, "'lifecycleConfigOutput' is nil, skipping further processing")
		return nil
	}

	rulesMatched := lifecycleConfigOutput.Rules
	declaredRules, ok := d.Get("lifecycle_rule").([]interface{})

	// We should match the existing lifecycle rules to the schema if they're defined
	if ok {
		rulesMatched = matchRulesWithSchema(rulesMatched, declaredRules)
	}

	d.Set("lifecycle_rule", flattenLifecycleRules(rulesMatched))

	return nil
}

func updateBucketVersioning(
	ctx context.Context,
	d *schema.ResourceData,
	client *s3.Client,
) error {
	bucket := d.Get("label").(string)
	n := d.Get("versioning").(bool)

	status := s3types.BucketVersioningStatusSuspended
	if n {
		status = s3types.BucketVersioningStatusEnabled
	}

	inputVersioningConfig := &s3.PutBucketVersioningInput{
		Bucket: &bucket,
		VersioningConfiguration: &s3types.VersioningConfiguration{
			Status: status,
		},
	}

	if _, err := client.PutBucketVersioning(ctx, inputVersioningConfig); err != nil {
		return err
	}

	return nil
}

func updateBucketLifecycle(
	ctx context.Context,
	d *schema.ResourceData,
	client *s3.Client,
) error {
	bucket := d.Get("label").(string)

	rules, err := expandLifecycleRules(d.Get("lifecycle_rule").([]interface{}))
	if err != nil {
		return err
	}

	_, err = client.PutBucketLifecycleConfiguration(
		ctx,
		&s3.PutBucketLifecycleConfigurationInput{
			Bucket: &bucket,
			LifecycleConfiguration: &s3types.BucketLifecycleConfiguration{
				Rules: rules,
			},
		},
	)

	return err
}

func updateBucketAccess(
	ctx context.Context, d *schema.ResourceData, client linodego.Client,
) error {
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
	ctx context.Context, d *schema.ResourceData, client linodego.Client,
) error {
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

func flattenLifecycleRules(rules []s3types.LifecycleRule) []map[string]interface{} {
	result := make([]map[string]interface{}, len(rules))

	for i, rule := range rules {
		ruleMap := make(map[string]interface{})

		if id := rule.ID; id != nil {
			ruleMap["id"] = *id
		}

		if prefix := rule.Prefix; prefix != nil {
			ruleMap["prefix"] = *prefix
		}

		ruleMap["enabled"] = rule.Status == s3types.ExpirationStatusEnabled

		if rule.AbortIncompleteMultipartUpload != nil && rule.AbortIncompleteMultipartUpload.DaysAfterInitiation != nil {
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

func expandLifecycleRules(ruleSpecs []interface{}) ([]s3types.LifecycleRule, error) {
	rules := make([]s3types.LifecycleRule, len(ruleSpecs))
	for i, ruleSpec := range ruleSpecs {
		ruleSpec := ruleSpec.(map[string]interface{})
		rule := s3types.LifecycleRule{}

		status := s3types.ExpirationStatusDisabled
		if ruleSpec["enabled"].(bool) {
			status = s3types.ExpirationStatusEnabled
		}
		rule.Status = status

		if id, ok := ruleSpec["id"]; ok {
			id := id.(string)
			rule.ID = &id
		}

		if prefix, ok := ruleSpec["prefix"]; ok {
			prefix := prefix.(string)
			rule.Prefix = &prefix
		}

		//nolint:lll
		abortIncompleteDays, ok := ruleSpec["abort_incomplete_multipart_upload_days"].(int)
		if ok && abortIncompleteDays > 0 {
			int32Days, err := helper.SafeIntToInt32(abortIncompleteDays)
			if err != nil {
				return nil, err
			}
			rule.AbortIncompleteMultipartUpload = &s3types.AbortIncompleteMultipartUpload{
				DaysAfterInitiation: &int32Days,
			}
		}

		if expirationList := ruleSpec["expiration"].([]interface{}); len(expirationList) > 0 {
			rule.Expiration = &s3types.LifecycleExpiration{}

			expirationMap := expirationList[0].(map[string]interface{})

			if dateStr, ok := expirationMap["date"].(string); ok && dateStr != "" {
				date, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", dateStr))
				if err != nil {
					return nil, err
				}

				rule.Expiration.Date = &date
			}

			if days, ok := expirationMap["days"].(int); ok && days > 0 {
				int32Days, err := helper.SafeIntToInt32(days)
				if err != nil {
					return nil, err
				}
				rule.Expiration.Days = &int32Days
			}

			if marker, ok := expirationMap["expired_object_delete_marker"].(bool); ok && marker {
				rule.Expiration.ExpiredObjectDeleteMarker = &marker
			}
		}

		if expirationList := ruleSpec["noncurrent_version_expiration"].([]interface{}); len(expirationList) > 0 {
			rule.NoncurrentVersionExpiration = &s3types.NoncurrentVersionExpiration{}

			expirationMap := expirationList[0].(map[string]interface{})

			if days, ok := expirationMap["days"].(int); ok && days > 0 {
				int32Days, err := helper.SafeIntToInt32(days)
				if err != nil {
					return nil, err
				}
				rule.NoncurrentVersionExpiration.NoncurrentDays = &int32Days
			}
		}

		rules[i] = rule
	}

	return rules, nil
}

// matchRulesWithSchema is for keeping the order of existing rules in the
// TF states and append any addition rules received
func matchRulesWithSchema(
	rules []s3types.LifecycleRule,
	declaredRules []interface{},
) []s3types.LifecycleRule {
	result := make([]s3types.LifecycleRule, 0)

	ruleMap := make(map[string]s3types.LifecycleRule, len(declaredRules))
	for _, rule := range rules {
		ruleMap[*rule.ID] = rule
	}

	for _, declaredRule := range declaredRules {
		declaredRule := declaredRule.(map[string]interface{})

		declaredID, ok := declaredRule["id"]

		if !ok || len(declaredID.(string)) < 1 {
			continue
		}

		if rule, ok := ruleMap[declaredID.(string)]; ok {
			result = append(result, rule)
			delete(ruleMap, declaredID.(string))
		}
	}

	// populate remaining values
	for _, rule := range ruleMap {
		result = append(result, rule)
	}

	return result
}
