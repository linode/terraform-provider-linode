package volume

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The unique id of this Volume.",
			Required:    true,
		},
		"label": schema.StringAttribute{
			Description: "The Volume's label. For display purposes only.",
			Computed:    true,
		},
		"region": schema.StringAttribute{
			Description: "The datacenter where this Volume is located.",
			Computed:    true,
		},
		"size": schema.Int64Attribute{
			Description: "The size of this Volume in GiB.",
			Computed:    true,
		},
		"linode_id": schema.Int64Attribute{
			Description: "If a Volume is attached to a specific Linode, the ID of that Linode will be displayed here.",
			Computed:    true,
		},
		"filesystem_path": schema.StringAttribute{
			Description: "The full filesystem path for the Volume based on the Volume's label. Path is " +
				"/dev/disk/by-id/scsi-0LinodeVolume + Volume label.",
			Computed: true,
		},
		"tags": schema.SetAttribute{
			ElementType: types.StringType,
			Description: "An array of tags applied to this Volume. Tags are for organizational purposes only.",
			Computed:    true,
		},
		"status": schema.StringAttribute{
			Description: "The status of the Volume. Can be one of active | creating | resizing | contact_support",
			Computed:    true,
		},
		"created": schema.StringAttribute{
			Description: "Datetime string representing when the Volume was created.",
			Computed:    true,
		},
		"updated": schema.StringAttribute{
			Description: "Datetime string representing when the Volume was last updated.",
			Computed:    true,
		},
	},
}
