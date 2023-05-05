package backup

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var diskObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"label":      types.StringType,
		"size":       types.Int64Type,
		"filesystem": types.StringType,
	},
}

var backupObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":        types.Int64Type,
		"label":     types.StringType,
		"status":    types.StringType,
		"type":      types.StringType,
		"created":   types.StringType,
		"updated":   types.StringType,
		"finished":  types.StringType,
		"configs":   types.ListType{ElemType: types.StringType},
		"disks":     types.ListType{ElemType: diskObjectType},
		"available": types.BoolType,
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"linode_id": schema.Int64Attribute{
			Description: "The ID of the Linode to get backups for.",
			Required:    true,
		},
		"automatic": schema.ListAttribute{
			Description: "A list of backups or snapshots for a Linode.",
			Computed:    true,
			ElementType: backupObjectType,
		},
		"current": schema.ListAttribute{
			Description: "The current Backup for a Linode.",
			Computed:    true,
			ElementType: backupObjectType,
		},
		"in_progress": schema.ListAttribute{
			Description: "The in-progress Backup for a Linode",
			Computed:    true,
			ElementType: backupObjectType,
		},
		"id": schema.Int64Attribute{
			Description: "The ID of the Backup",
			Computed:    true,
		},
	},
}
