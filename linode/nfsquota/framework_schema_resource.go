package nfsquota

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var identifierTypeValidator = stringvalidator.OneOf("uid", "username", "gid", "groupname")

func capacityLimitAttributes(description string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: description,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			"max_capacity_bytes": schema.Int64Attribute{
				Description: "The maximum capacity in bytes allowed by this limit.",
				Optional:    true,
			},
			"max_file_count": schema.Int64Attribute{
				Description: "The maximum number of files allowed by this limit.",
				Optional:    true,
			},
		},
	}
}

func identifiedLimitAttributes(description string) schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Description: description,
		Optional:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"identifier_type": schema.StringAttribute{
					Description: "The identifier type. One of: uid, username, gid, groupname.",
					Required:    true,
					Validators:  []validator.String{identifierTypeValidator},
				},
				"identifier": schema.StringAttribute{
					Description: "The user or group identifier value.",
					Required:    true,
				},
				"max_capacity_bytes": schema.Int64Attribute{
					Description: "The maximum capacity in bytes allowed for this identifier.",
					Optional:    true,
				},
				"max_file_count": schema.Int64Attribute{
					Description: "The maximum number of files allowed for this identifier.",
					Optional:    true,
				},
			},
		},
	}
}

func identityAttribute(description string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: description,
		Required:    true,
		Attributes: map[string]schema.Attribute{
			"identifier_type": schema.StringAttribute{
				Description: "The identifier type. One of: uid, username, gid, groupname.",
				Required:    true,
				Validators:  []validator.String{identifierTypeValidator},
			},
			"identifier": schema.StringAttribute{
				Description: "The user or group identifier value.",
				Required:    true,
			},
		},
	}
}

var frameworkResourceSchema = schema.Schema{
	Description: "Manages a Linode NFS Quota.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique ID of this NFS Quota.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"space_id": schema.Int64Attribute{
			Description: "The ID of the NFS Space this quota belongs to.",
			Required:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
		},
		"filesystem_id": schema.Int64Attribute{
			Description: "The ID of the NFS Filesystem this quota belongs to.",
			Required:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
		},
		"path": schema.StringAttribute{
			Description: "The path within the filesystem that this quota applies to.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"max_capacity_bytes": schema.Int64Attribute{
			Description: "The maximum capacity in bytes allowed by this quota.",
			Required:    true,
		},
		"max_file_count": schema.Int64Attribute{
			Description: "The maximum number of files allowed by this quota.",
			Required:    true,
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
			Optional:    true,
			Attributes: map[string]schema.Attribute{
				"default_user_limit":  capacityLimitAttributes("The default limit applied to each user."),
				"default_group_limit": capacityLimitAttributes("The default limit applied to each group."),
				"user_limits":         identifiedLimitAttributes("Per-user quota limits."),
				"group_limits":        identifiedLimitAttributes("Per-group quota limits."),
				"user_group_limits": schema.ListNestedAttribute{
					Description: "Combined user and group quota limits.",
					Optional:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"user":  identityAttribute("The user portion of the combined limit."),
							"group": identityAttribute("The group portion of the combined limit."),
							"max_capacity_bytes": schema.Int64Attribute{
								Description: "The maximum capacity in bytes allowed for this combined limit.",
								Optional:    true,
							},
							"max_file_count": schema.Int64Attribute{
								Description: "The maximum number of files allowed for this combined limit.",
								Optional:    true,
							},
						},
					},
				},
			},
		},
	},
}
