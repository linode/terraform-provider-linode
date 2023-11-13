package databasemysql

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"database_id": schema.Int64Attribute{
			DeprecationMessage: "Configure `id` instead. This attribute will be removed " +
				"in the next major version of the provider.",
			Description: "The ID of the MySQL database. DEPRECATED: Use ID instead",
			Optional:    true,
		},
		"engine_id": schema.StringAttribute{
			Description: "The Managed Database engine in engine/version format. (e.g. mysql/8.0.30)",
			Computed:    true,
		},
		"label": schema.StringAttribute{
			Description: "A unique, user-defined string referring to the Managed Database.",
			Computed:    true,
			Validators: []validator.String{
				stringvalidator.LengthBetween(3, 32),
			},
		},
		"region": schema.StringAttribute{
			Description: "The region to use for the Managed Database.",
			Computed:    true,
		},
		"type": schema.StringAttribute{
			Description: "The Linode Instance type used by the Managed Database for its nodes.",
			Computed:    true,
		},
		"allow_list": schema.SetAttribute{
			ElementType: types.StringType,
			Description: "A list of IP addresses that can access the Managed Database. " +
				"Each item can be a single IP address or a range in CIDR format.",
			Computed: true,
		},
		"cluster_size": schema.Int64Attribute{
			Description: "The number of Linode Instance nodes deployed to the Managed Database. Defaults to 1.",
			Computed:    true,
		},
		"encrypted": schema.BoolAttribute{
			Description: "Whether the Managed Databases is encrypted.",
			Computed:    true,
		},
		"replication_type": schema.StringAttribute{
			Description: "The replication method used for the Managed Database.",
			Computed:    true,
		},
		"ssl_connection": schema.BoolAttribute{
			Description: "Whether to require SSL credentials to establish a connection to the Managed Database.",
			Computed:    true,
		},
		"ca_cert": schema.StringAttribute{
			Description: "The base64-encoded SSL CA certificate for the Managed Database instance.",
			Computed:    true,
			Sensitive:   true,
		},
		"created": schema.StringAttribute{
			Description: "When this Managed Database was created.",
			Computed:    true,
		},
		"engine": schema.StringAttribute{
			Description: "The Managed Database engine.",
			Computed:    true,
		},
		"host_primary": schema.StringAttribute{
			Description: "The primary host for the Managed Database.",
			Computed:    true,
		},
		"host_secondary": schema.StringAttribute{
			Description: "The secondary host for the Managed Database.",
			Computed:    true,
		},
		"root_password": schema.StringAttribute{
			Description: "The randomly-generated root password for the Managed Database instance.",
			Computed:    true,
			Sensitive:   true,
		},
		"status": schema.StringAttribute{
			Description: "The operating status of the Managed Database.",
			Computed:    true,
		},
		"updated": schema.StringAttribute{
			Description: "When this Managed Database was last updated.",
			Computed:    true,
		},
		"updates": schema.ListAttribute{
			Description: "Configuration settings for automated patch update maintenance for the Managed Database.",
			Computed:    true,
			ElementType: helper.UpdateObjectType,
		},
		"root_username": schema.StringAttribute{
			Description: "The root username for the Managed Database instance.",
			Computed:    true,
			Sensitive:   true,
		},
		"version": schema.StringAttribute{
			Description: "The Managed Database engine version.",
			Computed:    true,
		},
		"id": schema.Int64Attribute{
			Description: "Unique identifier for this DataSource. The ID of the MySQL database.",
			Optional:    true,
			Validators: []validator.Int64{
				int64validator.ExactlyOneOf(path.Expressions{
					path.MatchRoot("database_id"),
				}...),
			},
		},
	},
}
