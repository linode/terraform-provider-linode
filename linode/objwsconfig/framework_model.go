package objwsconfig

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

// ResourceModel describes the Terraform resource rm model to match the
// resource schema.
type ResourceModel struct {
	Bucket          types.String `tfsdk:"bucket"`
	Cluster         types.String `tfsdk:"cluster"`
	AccessKey       types.String `tfsdk:"access_key"`
	SecretKey       types.String `tfsdk:"secret_key"`
	IndexDocument   types.String `tfsdk:"index_document"`
	ErrorDocument   types.String `tfsdk:"error_document"`
	WebsiteEndpoint types.String `tfsdk:"website_endpoint"`
}

func (rm *ResourceModel) CopyFrom(other ResourceModel, preserveKnown bool) {
	rm.Bucket = helper.KeepOrUpdateValue(rm.Bucket, other.Bucket, preserveKnown)
	rm.Cluster = helper.KeepOrUpdateValue(rm.Cluster, other.Cluster, preserveKnown)
	rm.AccessKey = helper.KeepOrUpdateValue(rm.AccessKey, other.AccessKey, preserveKnown)
	rm.SecretKey = helper.KeepOrUpdateValue(rm.SecretKey, other.SecretKey, preserveKnown)
	rm.IndexDocument = helper.KeepOrUpdateValue(rm.IndexDocument, other.IndexDocument, preserveKnown)
	rm.ErrorDocument = helper.KeepOrUpdateValue(rm.ErrorDocument, other.ErrorDocument, preserveKnown)
	rm.WebsiteEndpoint = helper.KeepOrUpdateValue(rm.WebsiteEndpoint, other.WebsiteEndpoint, preserveKnown)
}

func (rm *ResourceModel) ComputeWebsiteEndpoint(websiteDomain string) {
	rm.WebsiteEndpoint = types.StringValue(fmt.Sprintf("%s.%s", rm.Bucket.ValueString(), websiteDomain))
}

func (rm *ResourceModel) FlattenBucketWebsite(ws *s3.GetBucketWebsiteOutput, preserveKnown bool) {
	var indexDocument *string
	var errorDocument *string

	if ws.IndexDocument != nil && ws.IndexDocument.Suffix != nil {
		indexDocument = ws.IndexDocument.Suffix
	}
	if ws.ErrorDocument != nil && ws.ErrorDocument.Key != nil {
		errorDocument = ws.ErrorDocument.Key
	}

	rm.IndexDocument = helper.KeepOrUpdateStringPointer(rm.IndexDocument, indexDocument, preserveKnown)
	rm.ErrorDocument = helper.KeepOrUpdateStringPointer(rm.ErrorDocument, errorDocument, preserveKnown)
}
