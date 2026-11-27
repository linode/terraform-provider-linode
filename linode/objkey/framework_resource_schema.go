package objkey

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
)

var RegionDetailType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":            types.StringType,
		"s3_endpoint":   types.StringType,
		"endpoint_type": types.StringType,
	},
}

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
		"regions": schema.SetAttribute{
			Description: "A set of regions where the key will grant access to create buckets.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			Validators:  []validator.Set{setvalidator.SizeAtLeast(1)},
		},
		"regions_details": schema.SetAttribute{
			Description: "A set of objects containing the detailed info of the regions where " +
				"the key will grant access.",
			Computed:    true,
			ElementType: RegionDetailType,
		},
		"key_class": schema.StringAttribute{
			Description: "The storage class this key works with: E or IA. One key holds one class, " +
				"so its regions must all support it. When omitted, the API infers the class from the key's regions " +
				"and buckets. Without Infrequent Access API access, the API ignores this field.",
			Optional: true,
			Computed: true,
			Validators: []validator.String{
				stringvalidator.OneOf(string(linodego.ObjectStorageKeyClassE), string(linodego.ObjectStorageKeyClassIA)),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
				stringplanmodifier.UseStateForUnknown(),
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
					"region": schema.StringAttribute{
						Description: "The region where the bucket resides.",
						Required:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
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
