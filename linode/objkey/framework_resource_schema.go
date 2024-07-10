package objkey

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var RegionDetailType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":          types.StringType,
		"s3_endpoint": types.StringType,
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
		},
		"regions_details": schema.SetAttribute{
			Description: "A set of objects containing the detailed info of the regions where " +
				"the key will grant access.",
			Computed:    true,
			ElementType: RegionDetailType,
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
						Description: "The Object Storage cluster where the bucket resides. " +
							"Deprecated in favor of `region`",
						Optional: true,
						Computed: true,
						DeprecationMessage: "The `cluster` attribute in a `bucket_access` block has " +
							"been deprecated in favor of `region` attribute. A cluster value can be " +
							"converted to a region value by removing -x at the end, for example, a " +
							"cluster value `us-mia-1` can be converted to region value `us-mia`",
						Validators: []validator.String{
							stringvalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("region"),
							),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"region": schema.StringAttribute{
						Description: "The region where the bucket resides.",
						Optional:    true,
						Computed:    true,
						Validators: []validator.String{
							stringvalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("cluster"),
							),
						},
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
