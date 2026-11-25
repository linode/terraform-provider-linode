package nfsfilesystemaccesspolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
)

var frameworkResourceSchema = schema.Schema{
	Description: "Manages the access policy of a Linode NFS Filesystem. The access policy controls " +
		"which Linodes can mount the filesystem.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique ID of this resource, in the format `space_id,filesystem_id`.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"space_id": schema.Int64Attribute{
			Description: "The ID of the NFS Space the filesystem belongs to.",
			Required:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
		},
		"filesystem_id": schema.Int64Attribute{
			Description: "The ID of the NFS Filesystem this access policy belongs to.",
			Required:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
		},
		"label": schema.StringAttribute{
			Description: "The descriptive label of this access policy. Mirrors the label of the NFS Filesystem.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"enabled": schema.BoolAttribute{
			Description: "Whether this access policy is active. When false, no Linodes can mount " +
				"this filesystem.",
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"inherit": schema.BoolAttribute{
			Description: "Whether this filesystem falls back to the NFS Space access policy when no " +
				"`linode_ids` are configured. When false, access is restricted to `linode_ids` only.",
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"linode_ids": schema.SetAttribute{
			Description: "The IDs of the Linodes permitted to mount this filesystem. Each Linode " +
				"must belong to a VPC permitted by the NFS Space access policy.",
			Optional:    true,
			Computed:    true,
			ElementType: types.Int64Type,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		},
		"squash_policy": schema.StringAttribute{
			Description: "The user-ID squash policy. (`none`, `root_squash`, `all_squash`)",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf(
					string(linodego.NFSSquashPolicyNone),
					string(linodego.NFSSquashPolicyRootSquash),
					string(linodego.NFSSquashPolicyAllSquash),
				),
			},
		},
		"protocols": schema.SetAttribute{
			Description: "The NFS protocol versions permitted. Defaults to the filesystem's " +
				"`protocol_versions`.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		},
		"linode_acl": schema.ListNestedAttribute{
			Description: "The resolved Linode ACL entries permitted to mount this filesystem.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Description: "The ID of the Linode.",
						Computed:    true,
					},
					"label": schema.StringAttribute{
						Description: "The label of the Linode.",
						Computed:    true,
					},
					"url": schema.StringAttribute{
						Description: "The relative URL of the Linode.",
						Computed:    true,
					},
					"ipv6": schema.StringAttribute{
						Description: "The resolved IPv6 address of the Linode.",
						Computed:    true,
					},
				},
			},
		},
		"status": schema.StringAttribute{
			Description: "The current lifecycle status of the access policy.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"created": schema.StringAttribute{
			Description: "When this access policy was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated": schema.StringAttribute{
			Description: "When this access policy was last updated.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
	},
}
