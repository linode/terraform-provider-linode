package nfsfilesystem

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var resourceStatsAttribute = schema.SingleNestedAttribute{
	Description: "Live usage statistics for this filesystem.",
	Computed:    true,
	PlanModifiers: []planmodifier.Object{
		objectplanmodifier.UseStateForUnknown(),
	},
	Attributes: map[string]schema.Attribute{
		"used_capacity_bytes": schema.Int64Attribute{
			Description: "Capacity currently in use, in bytes.",
			Computed:    true,
		},
		"used_file_count": schema.Int64Attribute{
			Description: "Number of files and directories currently in use.",
			Computed:    true,
		},
		"collected_at": schema.StringAttribute{
			Description: "When these statistics were collected.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
	},
}

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique numeric ID of the NFS Filesystem.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"space_id": schema.Int64Attribute{
			Description: "The ID of the NFS Space this filesystem belongs to. Immutable after creation.",
			Required:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
		},
		"label": schema.StringAttribute{
			Description: "The customer-provided label of the NFS Filesystem. Unique within the NFS " +
				"Space. 3-63 characters, lowercase alphanumeric and hyphens only.",
			Required: true,
		},
		"region": schema.StringAttribute{
			Description: "The region where this filesystem is provisioned. Immutable after creation.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"max_capacity_bytes": schema.Int64Attribute{
			Description: "The provisioned capacity of the filesystem, in bytes.",
			Required:    true,
		},
		"protocol_versions": schema.ListAttribute{
			Description: "The NFS protocol versions supported by this filesystem. Currently only " +
				"'nfsv4' is supported.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
				listplanmodifier.RequiresReplace(),
			},
		},
		"status": schema.StringAttribute{
			Description: "The current lifecycle status of the filesystem.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"mount_target_ips": schema.ListAttribute{
			Description: "A list of NFS mount strings using NFS endpoint IPv6 addresses.",
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
		},
		"mount_target_fqdn": schema.StringAttribute{
			Description: "The full NFS mount string using the DNS-resolvable FQDN.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"snapshot_usage_bytes": schema.Int64Attribute{
			Description: "The total unique snapshot capacity in bytes for this filesystem.",
			Computed:    true,
		},
		"ldap_config_id": schema.StringAttribute{
			Description: "The ID of the LDAP Config associated with this filesystem's parent NFS " +
				"Space, if any.",
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"source_snapshot_id": schema.Int64Attribute{
			Description: "The source snapshot ID if this filesystem was created via snapshot clone.",
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"created": schema.StringAttribute{
			Description: "When this filesystem was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated": schema.StringAttribute{
			Description: "When this filesystem was last updated.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"tags": schema.SetAttribute{
			Description: "Customer-defined metadata tags.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		},
		"stats": resourceStatsAttribute,
	},
}
