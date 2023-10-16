package objkey

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"label": schema.StringAttribute{
			Description: "The label given to this key. For display purposes only.",
			Required:    true,
		},

		"id": schema.StringAttribute{
			Description: "The unique ID of this Object Storage key.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"access_key": schema.StringAttribute{
			Description: "This keypair's access key. This is not secret.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"secret_key": schema.StringAttribute{
			Description: "This keypair's secret key.",
			Sensitive:   true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"limited": schema.BoolAttribute{
			Description: "Whether or not this key is a limited access key.",
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
	},
	Blocks: map[string]schema.Block{
		"bucket_access": schema.SetNestedBlock{
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
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.RequiresReplace(),
				setplanmodifier.UseStateForUnknown(),
			},
		},
	},
}
