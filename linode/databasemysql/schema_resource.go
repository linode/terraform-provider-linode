package databasemysql

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var resourceSchema = map[string]*schema.Schema{
	// Required fields
	"engine": {
		Type:        schema.TypeString,
		Description: "The Managed Database engine type.",
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
		Type:        schema.TypeSet,
		Description: "A list of IP addresses that can access the Managed Database. Each item can be a single IP address or a range in CIDR format.",
		Optional:    true,
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
