package databasemysqlv2

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The id of the VPC.",
			Required:    true,
		},

		"engine_id": schema.StringAttribute{
			Description: "The unique ID of the database engine and version to use. (e.g. mysql/8)",
			Computed:    true,
		},
		"label": schema.StringAttribute{
			Description: "A unique, user-defined string referring to the Managed Database.",
			Computed:    true,
		},
		"region": schema.StringAttribute{
			Description: "The Region ID for the Managed Database.",
			Computed:    true,
		},
		"type": schema.StringAttribute{
			Description: "The Linode Instance type used by the Managed Database for its nodes.",
			Computed:    true,
		},

		"allow_list": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "A list of IP addresses that can access the Managed Database. " +
				"Each item can be a single IP address or a range in CIDR format.",
		},
		"ca_cert": schema.StringAttribute{
			Description: "The base64-encoded SSL CA certificate for the Managed Database.",
			Computed:    true,
			Sensitive:   true,
		},
		"cluster_size": schema.Int64Attribute{
			Computed:    true,
			Description: "The number of Linode instance nodes deployed to the Managed Database.",
		},
		"fork_restore_time": schema.StringAttribute{
			Description: "The database timestamp from which it was restored.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"fork_source": schema.Int64Attribute{
			Description: "The ID of the database that was forked from.",
			Computed:    true,
		},
		"updates": schema.ObjectAttribute{
			Description:    "Configuration settings for automated patch update maintenance for the Managed Database.",
			AttributeTypes: updatesAttributes,
			Computed:       true,
		},

		"created": schema.StringAttribute{
			Description: "When this Managed Database was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"encrypted": schema.BoolAttribute{
			Description: "Whether the Managed Databases is encrypted.",
			Computed:    true,
		},
		"engine": schema.StringAttribute{
			Description: "The Managed Database engine in engine/version format.",
			Computed:    true,
		},
		"host_primary": schema.StringAttribute{
			Description: "The primary host for the Managed Database.",
			Computed:    true,
		},
		"host_secondary": schema.StringAttribute{
			Description: "The secondary/private host for the Managed Database.",
			Computed:    true,
		},
		"members": schema.MapAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "A mapping between IP addresses and strings designating them as primary or failover.",
		},
		"oldest_restore_time": schema.StringAttribute{
			Description: "The oldest time to which a database can be restored.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"pending_updates": schema.SetAttribute{
			Description: "A set of pending updates.",
			Computed:    true,
			ElementType: types.ObjectType{AttrTypes: pendingUpdateAttributes},
		},
		"platform": schema.StringAttribute{
			Computed:    true,
			Description: "The back-end platform for relational databases used by the service.",
		},
		"port": schema.Int64Attribute{
			Description: "The access port for this Managed Database.",
			Computed:    true,
		},
		"root_password": schema.StringAttribute{
			Description: "The randomly generated root password for the Managed Database instance.",
			Computed:    true,
			Sensitive:   true,
		},
		"root_username": schema.StringAttribute{
			Description: "The root username for the Managed Database instance.",
			Computed:    true,
			Sensitive:   true,
		},
		"ssl_connection": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether to require SSL credentials to establish a connection to the Managed Database.",
		},
		"status": schema.StringAttribute{
			Computed:    true,
			Description: "The operating status of the Managed Database.",
		},
		"suspended": schema.BoolAttribute{
			Description: "Whether this database is suspended.",
			Computed:    true,
		},
		"updated": schema.StringAttribute{
			Description: "When this Managed Database was last updated.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"version": schema.StringAttribute{
			Description: "The Managed Database engine version.",
			Computed:    true,
		},
	},
}
