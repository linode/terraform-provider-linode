package objbucket

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"cluster": schema.StringAttribute{
			Description: "The ID of the Object Storage Cluster this bucket is in.",
			Required:    true,
		},
		"created": schema.StringAttribute{
			Description: "When this bucket was created.",
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
