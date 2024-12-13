package databasepostgresqlv2

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

var (
	allowListDefault = helper.FrameworkMust(
		types.SetValue(types.StringType, []attr.Value{types.StringValue("0.0.0.0/0")}),
	)

	updatesAttributes = map[string]attr.Type{
		"day_of_week": types.Int64Type,
		"duration":    types.Int64Type,
		"frequency":   types.StringType,
		"hour_of_day": types.Int64Type,
	}

	pendingUpdateAttributes = map[string]attr.Type{
		"deadline":    timetypes.RFC3339Type{},
		"description": types.StringType,
		"planned_for": timetypes.RFC3339Type{},
	}
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The id of the VPC.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},

		"engine_id": schema.StringAttribute{
			Required:    true,
			Description: "The unique ID of the database engine and version to use. (e.g. postgresql/16)",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"label": schema.StringAttribute{
			Required:    true,
			Description: "A unique, user-defined string referring to the Managed Database.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"region": schema.StringAttribute{
			Required:    true,
			Description: "The Region ID for the Managed Database.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"type": schema.StringAttribute{
			Required:    true,
			Description: "The Linode Instance type used by the Managed Database for its nodes.\n\n",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},

		"allow_list": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
			Description: "A list of IP addresses that can access the Managed Database. " +
				"Each item can be a single IP address or a range in CIDR format.",
			Default: setdefault.StaticValue(allowListDefault),
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		},
		"ca_cert": schema.StringAttribute{
			Description: "The base64-encoded SSL CA certificate for the Managed Database.",
			Computed:    true,
			Sensitive:   true,
		},
		"cluster_size": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "The number of Linode instance nodes deployed to the Managed Database.",
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
			},
			Default: int64default.StaticInt64(1),
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"fork_restore_time": schema.StringAttribute{
			Description: "The database timestamp from which it was restored.",
			Optional:    true,
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
				stringplanmodifier.RequiresReplace(),
			},
		},
		"fork_source": schema.Int64Attribute{
			Description: "The ID of the database that was forked from.",
			Optional:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
				int64planmodifier.RequiresReplace(),
			},
		},
		"updates": schema.ObjectAttribute{
			Description:    "Configuration settings for automated patch update maintenance for the Managed Database.",
			AttributeTypes: updatesAttributes,
			Computed:       true,
			Optional:       true,
			PlanModifiers:  []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
		},

		"created": schema.StringAttribute{
			Description: "When this Managed Database was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"encrypted": schema.BoolAttribute{
			Description: "Whether the Managed Databases is encrypted.",
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"engine": schema.StringAttribute{
			Description: "The Managed Database engine in engine/version format.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
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
			Description:   "A set of pending updates.",
			Computed:      true,
			ElementType:   types.ObjectType{AttrTypes: pendingUpdateAttributes},
			PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
		},
		"platform": schema.StringAttribute{
			Computed:      true,
			Description:   "The back-end platform for relational databases used by the service.",
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"port": schema.Int64Attribute{
			Description:   "The access port for this Managed Database.",
			Computed:      true,
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		},
		"root_password": schema.StringAttribute{
			Description:   "The randomly generated root password for the Managed Database instance.",
			Computed:      true,
			Sensitive:     true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"root_username": schema.StringAttribute{
			Description:   "The root username for the Managed Database instance.",
			Computed:      true,
			Sensitive:     true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"ssl_connection": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether to require SSL credentials to establish a connection to the Managed Database.",
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"status": schema.StringAttribute{
			Computed:    true,
			Description: "The operating status of the Managed Database.",
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
