package databasemysql

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

var resourceSchema = map[string]*schema.Schema{
	// Required fields
	"engine_id": {
		Type:        schema.TypeString,
		Description: "The Managed Database engine in engine/version format. (e.g. mysql/8.0.26)",
		Required:    true,
		ForceNew:    true,
	},
	"label": {
		Type:        schema.TypeString,
		Description: "A unique, user-defined string referring to the Managed Database.",
		Required:    true,
	},
	"region": {
		Type:        schema.TypeString,
		Description: "The region to use for the Managed Database.",
		Required:    true,
		ForceNew:    true,
	},
	"type": {
		Type:        schema.TypeString,
		Description: "The Linode Instance type used by the Managed Database for its nodes.",
		Required:    true,
		ForceNew:    true,
	},

	// Optional fields
	"allow_list": {
		Type: schema.TypeSet,
		Description: "A list of IP addresses that can access the Managed Database. " +
			"Each item can be a single IP address or a range in CIDR format.",
		Optional: true,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"cluster_size": {
		Type:        schema.TypeInt,
		Description: "The number of Linode Instance nodes deployed to the Managed Database. Defaults to 1.",
		Optional:    true,
		Default:     1,
		ForceNew:    true,

		// We should enforce this during plan to prevent invalid DBs from being attempted
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{1, 3})),
	},
	"encrypted": {
		Type:        schema.TypeBool,
		Description: "Whether the Managed Databases is encrypted.",
		Optional:    true,
		Default:     false,
		ForceNew:    true,
	},
	"replication_type": {
		Type:        schema.TypeString,
		Description: "The replication method used for the Managed Database.",
		Optional:    true,
		Default:     "none",
		ForceNew:    true,

		ValidateDiagFunc: validation.ToDiagFunc(
			validation.StringInSlice([]string{"none", "asynch", "semi_synch"}, false)),
	},
	"ssl_connection": {
		Type:        schema.TypeBool,
		Description: "Whether to require SSL credentials to establish a connection to the Managed Database.",
		Optional:    true,
		Default:     false,
		ForceNew:    true,
	},
	"updates": {
		Type:        schema.TypeList,
		Description: "Configuration settings for automated patch update maintenance for the Managed Database.",
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"day_of_week": {
					Type:        schema.TypeString,
					Description: "The day to perform maintenance.",
					Required:    true,
					ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
						if _, err := helper.ExpandDayOfWeek(i.(string)); err != nil {
							return diag.FromErr(err)
						}

						return nil
					},
				},
				"duration": {
					Type:             schema.TypeInt,
					Description:      "The maximum maintenance window time in hours.",
					Required:         true,
					ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 3)),
				},
				"frequency": {
					Type:        schema.TypeString,
					Description: "Whether maintenance occurs on a weekly or monthly basis.",
					Required:    true,
					ValidateDiagFunc: validation.ToDiagFunc(
						validation.StringInSlice([]string{"weekly", "monthly"}, true)),
				},
				"hour_of_day": {
					Type:             schema.TypeInt,
					Description:      "The hour to begin maintenance based in UTC time.",
					Required:         true,
					ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 23)),
				},
				"week_of_month": {
					Type: schema.TypeInt,
					Description: "The week of the month to perform monthly frequency updates." +
						" Required for monthly frequency updates.",
					ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 4)),
					Optional:         true,
					Default:          nil,
				},
			},
		},
	},

	// Computed fields
	"ca_cert": {
		Type:        schema.TypeString,
		Description: "The base64-encoded SSL CA certificate for the Managed Database instance.",
		Computed:    true,
		Sensitive:   true,
	},
	"created": {
		Type:        schema.TypeString,
		Description: "When this Managed Database was created.",
		Computed:    true,
	},
	"engine": {
		Type:        schema.TypeString,
		Description: "The Managed Database engine.",
		Computed:    true,
	},
	// In an attempt to simplify provider/config logic, these should be flattened in schema
	"host_primary": {
		Type:        schema.TypeString,
		Description: "The primary host for the Managed Database.",
		Computed:    true,
	},
	"host_secondary": {
		Type:        schema.TypeString,
		Description: "The secondary host for the Managed Database.",
		Computed:    true,
	},
	"root_password": {
		Type:        schema.TypeString,
		Description: "The randomly-generated root password for the Managed Database instance.",
		Computed:    true,
		Sensitive:   true,
	},
	"status": {
		Type:        schema.TypeString,
		Description: "The operating status of the Managed Database.",
		Computed:    true,
	},
	"updated": {
		Type:        schema.TypeString,
		Description: "When this Managed Database was last updated.",
		Computed:    true,
	},
	"root_username": {
		Type:        schema.TypeString,
		Description: "The root username for the Managed Database instance.",
		Computed:    true,
		Sensitive:   true,
	},
	"version": {
		Type:        schema.TypeString,
		Description: "The Managed Database engine version.",
		Computed:    true,
	},
}
