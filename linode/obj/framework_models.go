package obj

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type BaseModel struct {
	Bucket             types.String `tfsdk:"bucket"`
	Cluster            types.String `tfsdk:"cluster"`
	Region             types.String `tfsdk:"region"`
	Key                types.String `tfsdk:"key"`
	SecreteKey         types.String `tfsdk:"secret_key"`
	AccessKey          types.String `tfsdk:"access_key"`
	Content            types.String `tfsdk:"content"`
	ContentBase64      types.String `tfsdk:"content_base64"`
	Source             types.String `tfsdk:"source"`
	ACL                types.String `tfsdk:"acl"`
	CacheControl       types.String `tfsdk:"cache_control"`
	ContentDisposition types.String `tfsdk:"content_disposition"`
	ContentEncoding    types.String `tfsdk:"content_encoding"`
	ContentLanguage    types.String `tfsdk:"content_language"`
	ContentType        types.String `tfsdk:"content_type"`
	Endpoint           types.String `tfsdk:"endpoint"`
	ETag               types.String `tfsdk:"etag"`
	ForceDestroy       types.Bool   `tfsdk:"force_destroy"`
	Metadata           types.Map    `tfsdk:"metadata"`
	VersionID          types.String `tfsdk:"version_id"`
	WebsiteRedirect    types.String `tfsdk:"website_redirect"`
}

// TODO: consider merging two models when resource's ID change to int type
type ResourceModel struct {
	ID types.String `tfsdk:"id"`
	BaseModel
}

func (data ResourceModel) getObjectBody(diags *diag.Diagnostics) (body *s3manager.ReaderSeekerCloser) {
	if !data.Source.IsNull() && !data.Source.IsUnknown() {
		sourcePath := data.Source.ValueString()

		file, err := os.Open(filepath.Clean(sourcePath))
		if err != nil {
			diags.AddError(fmt.Sprintf("Failed to Open the File at %q", sourcePath), err.Error())
			return
		}

		return s3manager.ReadSeekCloser(file)
	}

	var contentBytes []byte
	var err error

	if !data.ContentBase64.IsNull() && !data.ContentBase64.IsUnknown() {
		contentBytes, err = base64.StdEncoding.DecodeString(data.ContentBase64.ValueString())
		if err != nil {
			diags.AddError("Failed to Decode the base64 Content", err.Error())
		}
	} else if !data.Content.IsNull() && !data.Content.IsUnknown() {
		contentBytes = []byte(data.Content.ValueString())
	}

	return s3manager.ReadSeekCloser(bytes.NewReader(contentBytes))
}

func (data ResourceModel) GetObjectStorageKeys(
	ctx context.Context,
	client *linodego.Client,
	config *helper.FrameworkProviderModel,
	permissions string,
	diags *diag.Diagnostics,
) (*ObjectKeys, func()) {
	result := &ObjectKeys{}

	result.AccessKey = data.AccessKey.ValueString()
	result.SecretKey = data.SecreteKey.ValueString()

	if result.Ok() {
		return result, nil
	}

	result.AccessKey = config.ObjAccessKey.ValueString()
	result.SecretKey = config.ObjSecretKey.ValueString()

	if result.Ok() {
		return result, nil
	}

	if config.ObjUseTempKeys.ValueBool() {
		objKey := fwCreateTempKeys(ctx, client, data.Bucket.ValueString(), data.GetRegionOrCluster(ctx), permissions, diags)
		if diags.HasError() {
			return nil, nil
		}

		result.AccessKey = objKey.AccessKey
		result.SecretKey = objKey.SecretKey

		teardownTempKeysCleanUp := func() { cleanUpTempKeys(ctx, client, objKey.ID) }

		return result, teardownTempKeysCleanUp
	}

	diags.AddError(
		"Keys Not Found",
		"`access_key` and `secret_key` are Required but not Configured",
	)

	return nil, nil
}

func (plan *ResourceModel) ComputeEndpointIfUnknown(ctx context.Context, client *linodego.Client, diags *diag.Diagnostics) {
	if !plan.Endpoint.IsUnknown() {
		return
	}

	bucketName := plan.Bucket.ValueString()
	regionOrCluster := plan.GetRegionOrCluster(ctx)

	bucket, err := client.GetObjectStorageBucket(ctx, regionOrCluster, bucketName)
	if err != nil {
		diags.AddError(
			"Failed to Find the Specified Linode ObjectStorageBucket",
			err.Error(),
		)
		return
	}

	plan.Endpoint = types.StringValue(
		strings.TrimPrefix(bucket.Hostname, fmt.Sprintf("%s.", bucket.Label)),
	)
}

func (data *ResourceModel) GenerateObjectStorageObjectID(apply bool, preserveKnown bool) string {
	id := fmt.Sprintf("%s/%s", data.Bucket.ValueString(), data.Key.ValueString())

	if apply {
		data.ID = types.StringValue(id)
	}

	return id
}

func (data ResourceModel) GetRegionOrCluster(ctx context.Context) string {
	if !data.Region.IsNull() && !data.Region.IsUnknown() {
		return data.Region.ValueString()
	}

	return data.Cluster.ValueString()
}

func (data *ResourceModel) FlattenObject(
	obj s3.HeadObjectOutput, preserveKnown bool,
) {
	data.CacheControl = helper.KeepOrUpdateStringPointer(data.CacheControl, obj.CacheControl, preserveKnown)
	data.ContentDisposition = helper.KeepOrUpdateStringPointer(data.ContentDisposition, obj.ContentDisposition, preserveKnown)
	data.ContentEncoding = helper.KeepOrUpdateStringPointer(data.ContentEncoding, obj.ContentEncoding, preserveKnown)
	data.ContentLanguage = helper.KeepOrUpdateStringPointer(data.ContentLanguage, obj.ContentLanguage, preserveKnown)
	data.ContentType = helper.KeepOrUpdateStringPointer(data.ContentType, obj.ContentType, preserveKnown)
	data.ETag = helper.KeepOrUpdateStringPointer(data.ETag, getQuotesTrimmedETag(obj), preserveKnown)
	data.WebsiteRedirect = helper.KeepOrUpdateStringPointer(data.WebsiteRedirect, obj.WebsiteRedirectLocation, preserveKnown)
	data.VersionID = helper.KeepOrUpdateStringPointer(data.VersionID, obj.VersionId, preserveKnown)
	data.Metadata = helper.KeepOrUpdateValue(data.Metadata, types.MapValueMust(types.StringType, flattenObjectMetadata(obj.Metadata)), preserveKnown)
	data.ContentDisposition = helper.KeepOrUpdateStringPointer(data.ContentDisposition, obj.ContentDisposition, preserveKnown)

	data.GenerateObjectStorageObjectID(true, preserveKnown)
}

func (data ResourceModel) ETagChanged(
	obj s3.HeadObjectOutput,
) bool {
	return !data.ETag.Equal(types.StringPointerValue(getQuotesTrimmedETag(obj)))
}

func (plan *ResourceModel) CopyFrom(state ResourceModel, preserveKnown bool) {
	plan.ID = helper.KeepOrUpdateValue(plan.ID, state.ID, preserveKnown)
	plan.Bucket = helper.KeepOrUpdateValue(plan.Bucket, state.Bucket, preserveKnown)
	plan.Cluster = helper.KeepOrUpdateValue(plan.Cluster, state.Cluster, preserveKnown)
	plan.Region = helper.KeepOrUpdateValue(plan.Region, state.Region, preserveKnown)
	plan.Key = helper.KeepOrUpdateValue(plan.Key, state.Key, preserveKnown)
	plan.SecreteKey = helper.KeepOrUpdateValue(plan.SecreteKey, state.SecreteKey, preserveKnown)
	plan.AccessKey = helper.KeepOrUpdateValue(plan.AccessKey, state.AccessKey, preserveKnown)
	plan.Content = helper.KeepOrUpdateValue(plan.Content, state.Content, preserveKnown)
	plan.ContentBase64 = helper.KeepOrUpdateValue(plan.ContentBase64, state.ContentBase64, preserveKnown)
	plan.Source = helper.KeepOrUpdateValue(plan.Source, state.Source, preserveKnown)
	plan.ACL = helper.KeepOrUpdateValue(plan.ACL, state.ACL, preserveKnown)
	plan.CacheControl = helper.KeepOrUpdateValue(plan.CacheControl, state.CacheControl, preserveKnown)
	plan.ContentDisposition = helper.KeepOrUpdateValue(plan.ContentDisposition, state.ContentDisposition, preserveKnown)
	plan.ContentEncoding = helper.KeepOrUpdateValue(plan.ContentEncoding, state.ContentEncoding, preserveKnown)
	plan.ContentLanguage = helper.KeepOrUpdateValue(plan.ContentLanguage, state.ContentLanguage, preserveKnown)
	plan.ContentType = helper.KeepOrUpdateValue(plan.ContentType, state.ContentType, preserveKnown)
	plan.Endpoint = helper.KeepOrUpdateValue(plan.Endpoint, state.Endpoint, preserveKnown)
	plan.ETag = helper.KeepOrUpdateValue(plan.ETag, state.ETag, preserveKnown)
	plan.ForceDestroy = helper.KeepOrUpdateValue(plan.ForceDestroy, state.ForceDestroy, preserveKnown)
	plan.Metadata = helper.KeepOrUpdateValue(plan.Metadata, state.Metadata, preserveKnown)
	plan.VersionID = helper.KeepOrUpdateValue(plan.VersionID, state.VersionID, preserveKnown)
	plan.WebsiteRedirect = helper.KeepOrUpdateValue(plan.WebsiteRedirect, state.WebsiteRedirect, preserveKnown)
}
