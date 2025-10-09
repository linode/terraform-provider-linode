package objbucket

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"cluster": schema.StringAttribute{
			Description: "The ID of the Object Storage Cluster this bucket is in.",
			DeprecationMessage: "The cluster attribute has been deprecated, please consider " +
				"switching to the region attribute. For example, a cluster value of `us-mia-1` " +
				"can be translated to a region value of `us-mia`.",
			Optional: true,
			Computed: true,
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(
					path.MatchRelative().AtParent().AtName("region"),
				),
			},
		},
		"region": schema.StringAttribute{
			Description: "The ID of the region this bucket is in.",
			Optional:    true,
			Computed:    true,
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(
					path.MatchRelative().AtParent().AtName("cluster"),
				),
			},
		},
		"endpoint_type": schema.StringAttribute{
			Description: "The type of the S3 endpoint of the bucket.",
			Computed:    true,
		},
		"s3_endpoint": schema.StringAttribute{
			Description: "The S3 endpoint URL of the bucket, based on the `endpoint_type` and `region`.",
			Computed:    true,
		},
		"created": schema.StringAttribute{
			Description: "When this bucket was created.",
			CustomType:  timetypes.RFC3339Type{},
			Computed:    true,
		},
		"hostname": schema.StringAttribute{
			Description: "The hostname where this bucket can be accessed." +
				"This hostname can be accessed through a browser if the bucket is made public.",
			Computed: true,
		},
		"id": schema.StringAttribute{
			Description: "The id of this bucket.",
			Computed:    true,
		},
		"label": schema.StringAttribute{
			Description: "The name of this bucket.",
			Required:    true,
		},
		"objects": schema.Int64Attribute{
			Description: "The number of objects stored in this bucket.",
			Computed:    true,
		},
		"size": schema.Int64Attribute{
			Description: "The size of the bucket in bytes.",
			Computed:    true,
		},
	},
}
