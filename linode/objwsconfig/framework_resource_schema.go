package objwsconfig

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"bucket": schema.StringAttribute{
			Description: "The target bucket to create website.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"cluster": schema.StringAttribute{
			Description: "The target cluster that the bucket is in.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"access_key": schema.StringAttribute{
			Description: "The S3 access key with access to the target bucket.",
			Required:    true,
		},
		"secret_key": schema.StringAttribute{
			Description: "The S3 secret key with access to the target bucket.",
			Sensitive:   true,
			Required:    true,
		},
		"index_document": schema.StringAttribute{
			Description: "Object suffix that is appended to a request that is for a directory on the website endpoint. The index document must not be empty and must not include a slash character.",
			Required:    true,
		},
		"error_document": schema.StringAttribute{
			Description: "Object key name to use when error occurs.",
			Optional:    true,
		},
		"website_endpoint": schema.StringAttribute{
			Description: "The endpoint with bucket to access.",
			Computed:    true,
		},
	},
}
