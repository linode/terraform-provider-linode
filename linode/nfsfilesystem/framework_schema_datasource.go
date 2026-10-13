package nfsfilesystem

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var FilesystemAttributes = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description: "The unique numeric ID of the NFS Filesystem.",
		Required:    true,
	},
	"space_id": schema.Int64Attribute{
		Description: "The ID of the NFS Space this filesystem belongs to.",
		Computed:    true,
	},
	"label": schema.StringAttribute{
		Description: "The customer-provided label of the NFS Filesystem.",
		Computed:    true,
	},
	"region": schema.StringAttribute{
		Description: "The region where this filesystem is provisioned.",
		Computed:    true,
	},
	"max_capacity_bytes": schema.Int64Attribute{
		Description: "The provisioned capacity of the filesystem, in bytes.",
		Computed:    true,
	},
	"protocol_versions": schema.ListAttribute{
		Description: "The NFS protocol versions supported by this filesystem.",
		Computed:    true,
		ElementType: types.StringType,
	},
	"status": schema.StringAttribute{
		Description: "The current lifecycle status of the filesystem.",
		Computed:    true,
	},
	"mount_target_ips": schema.ListAttribute{
		Description: "A list of NFS mount strings using NFS endpoint IPv6 addresses.",
		Computed:    true,
		ElementType: types.StringType,
	},
	"mount_target_fqdn": schema.StringAttribute{
		Description: "The full NFS mount string using the DNS-resolvable FQDN.",
		Computed:    true,
	},
	"snapshot_usage_bytes": schema.Int64Attribute{
		Description: "The total unique snapshot capacity in bytes for this filesystem.",
		Computed:    true,
	},
	"ldap_config_id": schema.StringAttribute{
		Description: "The ID of the LDAP Config associated with this filesystem's parent NFS Space, if any.",
		Computed:    true,
	},
	"source_snapshot_id": schema.Int64Attribute{
		Description: "The source snapshot ID if this filesystem was created via snapshot clone.",
		Computed:    true,
	},
	"created": schema.StringAttribute{
		Description: "When this filesystem was created.",
		Computed:    true,
		CustomType:  timetypes.RFC3339Type{},
	},
	"updated": schema.StringAttribute{
		Description: "When this filesystem was last updated.",
		Computed:    true,
		CustomType:  timetypes.RFC3339Type{},
	},
	"tags": schema.SetAttribute{
		Description: "Customer-defined metadata tags.",
		Computed:    true,
		ElementType: types.StringType,
	},
	"stats": schema.SingleNestedAttribute{
		Description: "Live usage statistics for this filesystem.",
		Computed:    true,
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
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: FilesystemAttributes,
}
