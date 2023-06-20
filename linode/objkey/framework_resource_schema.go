package objkey

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"label": schema.StringAttribute{
			Description: "The label given to this key. For display purposes only.",
			Required:    true,
		},

		"id": schema.Int64Attribute{
			Description: "The unique ID of this Object Storage key.",
			Computed:    true,
		},
		"access_key": schema.StringAttribute{
			Description: "This keypair's access key. This is not secret.",
			Computed:    true,
		},
		"secret_key": schema.StringAttribute{
			Description: "This keypair's secret key.",
			Sensitive:   true,
			Computed:    true,
		},
		"limited": schema.BoolAttribute{
			Description: "Whether or not this key is a limited access key.",
			Computed:    true,
		},
	},
	Blocks: map[string]schema.Block{
		"bucket_access": schema.ListNestedBlock{
			// TODO: force new
			Description: "A list of permissions to grant this limited access key.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"bucket_name": schema.StringAttribute{
						Description: "The unique label of the bucket to which the key will grant limited access.",
						Required:    true,
					},
					"cluster": schema.StringAttribute{
						Description: "The Object Storage cluster where a bucket to which the key is granting access is hosted.",
						Required:    true,
					},
					"permissions": schema.StringAttribute{
						Description: "This Limited Access Key's permissions for the selected bucket.",
						Required:    true,
						// TODO: validate perms
					},
				},
			},
		},
	},
}
