package nfsfilesystemaccesspolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkDatasourceSchema = schema.Schema{
	Description: "Provides information about the access policy of a Linode NFS Filesystem.",
	Attributes: map[string]schema.Attribute{
		"space_id": schema.Int64Attribute{
			Description: "The ID of the NFS Space the filesystem belongs to.",
			Required:    true,
		},
		"filesystem_id": schema.Int64Attribute{
			Description: "The ID of the NFS Filesystem this access policy belongs to.",
			Required:    true,
		},
		"id": schema.StringAttribute{
			Description: "The unique ID of this access policy, in the format `space_id,filesystem_id`.",
			Computed:    true,
		},
		"configured": schema.BoolAttribute{
			Description: "Whether the access policy has been configured. When false, the filesystem " +
				"follows the NFS Space access policy and all other policy attributes are null.",
			Computed: true,
		},
		"label": schema.StringAttribute{
			Description: "The descriptive label of this access policy. Mirrors the label of the NFS Filesystem.",
			Computed:    true,
		},
		"enabled": schema.BoolAttribute{
			Description: "Whether this access policy is active. When false, no Linodes can mount " +
				"this filesystem.",
			Computed: true,
		},
		"inherit": schema.BoolAttribute{
			Description: "Whether this filesystem falls back to the NFS Space access policy when no " +
				"`linode_ids` are configured.",
			Computed: true,
		},
		"linode_ids": schema.SetAttribute{
			Description: "The IDs of the Linodes permitted to mount this filesystem.",
			Computed:    true,
			ElementType: types.Int64Type,
		},
		"squash_policy": schema.StringAttribute{
			Description: "The user-ID squash policy.",
			Computed:    true,
		},
		"protocols": schema.SetAttribute{
			Description: "The NFS protocol versions permitted.",
			Computed:    true,
			ElementType: types.StringType,
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
		},
		"created": schema.StringAttribute{
			Description: "When this access policy was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"updated": schema.StringAttribute{
			Description: "When this access policy was last updated.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
	},
}
