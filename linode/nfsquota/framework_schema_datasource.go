package nfsquota

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func dsCapacityLimitAttributes(description string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: description,
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"max_capacity_bytes": schema.Int64Attribute{
				Description: "The maximum capacity in bytes allowed by this limit.",
				Computed:    true,
			},
			"max_file_count": schema.Int64Attribute{
				Description: "The maximum number of files allowed by this limit.",
				Computed:    true,
			},
		},
	}
}

func dsIdentifiedLimitAttributes(description string) schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Description: description,
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"identifier_type": schema.StringAttribute{
					Description: "The identifier type. One of: uid, username, gid, groupname.",
					Computed:    true,
				},
				"identifier": schema.StringAttribute{
					Description: "The user or group identifier value.",
					Computed:    true,
				},
				"max_capacity_bytes": schema.Int64Attribute{
					Description: "The maximum capacity in bytes allowed for this identifier.",
					Computed:    true,
				},
				"max_file_count": schema.Int64Attribute{
					Description: "The maximum number of files allowed for this identifier.",
					Computed:    true,
				},
			},
		},
	}
}

func dsIdentityAttribute(description string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: description,
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"identifier_type": schema.StringAttribute{
				Description: "The identifier type. One of: uid, username, gid, groupname.",
				Computed:    true,
			},
			"identifier": schema.StringAttribute{
				Description: "The user or group identifier value.",
				Computed:    true,
			},
		},
	}
}

var DataSourceSchema = schema.Schema{
	Description: "Provides details about a Linode NFS Quota.",
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The unique ID of this NFS Quota.",
			Required:    true,
		},
		"space_id": schema.Int64Attribute{
			Description: "The ID of the NFS Space this quota belongs to.",
			Required:    true,
		},
		"filesystem_id": schema.Int64Attribute{
			Description: "The ID of the NFS Filesystem this quota belongs to.",
			Required:    true,
		},
		"path": schema.StringAttribute{
			Description: "The path within the filesystem that this quota applies to.",
			Computed:    true,
		},
		"max_capacity_bytes": schema.Int64Attribute{
			Description: "The maximum capacity in bytes allowed by this quota.",
			Computed:    true,
		},
		"max_file_count": schema.Int64Attribute{
			Description: "The maximum number of files allowed by this quota.",
			Computed:    true,
		},
		"is_root": schema.BoolAttribute{
			Description: "Whether this quota applies to the root of the filesystem.",
			Computed:    true,
		},
		"used_capacity_bytes": schema.Int64Attribute{
			Description: "The capacity in bytes currently used under this quota.",
			Computed:    true,
		},
		"quota_exceeded": schema.BoolAttribute{
			Description: "Whether the quota has been exceeded.",
			Computed:    true,
		},
		"used_file_count": schema.Int64Attribute{
			Description: "The number of files currently used under this quota.",
			Computed:    true,
		},
		"status": schema.StringAttribute{
			Description: "The current status of this NFS Quota.",
			Computed:    true,
		},
		"created": schema.StringAttribute{
			Description: "When this NFS Quota was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"updated": schema.StringAttribute{
			Description: "When this NFS Quota was last updated.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"collected_at": schema.StringAttribute{
			Description: "When this NFS Quota's usage data was last collected.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"user_group_config": schema.SingleNestedAttribute{
			Description: "User and group level quota configuration.",
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"default_user_limit":  dsCapacityLimitAttributes("The default limit applied to each user."),
				"default_group_limit": dsCapacityLimitAttributes("The default limit applied to each group."),
				"user_limits":         dsIdentifiedLimitAttributes("Per-user quota limits."),
				"group_limits":        dsIdentifiedLimitAttributes("Per-group quota limits."),
				"user_group_limits": schema.ListNestedAttribute{
					Description: "Combined user and group quota limits.",
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"user":  dsIdentityAttribute("The user portion of the combined limit."),
							"group": dsIdentityAttribute("The group portion of the combined limit."),
							"max_capacity_bytes": schema.Int64Attribute{
								Description: "The maximum capacity in bytes allowed for this combined limit.",
								Computed:    true,
							},
							"max_file_count": schema.Int64Attribute{
								Description: "The maximum number of files allowed for this combined limit.",
								Computed:    true,
							},
						},
					},
				},
			},
		},
	},
}
