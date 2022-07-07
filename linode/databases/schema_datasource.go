package databases

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

var filterConfig = helper.FilterConfig{
	"engine":  {TypeFunc: helper.FilterTypeString},
	"region":  {TypeFunc: helper.FilterTypeString},
	"status":  {TypeFunc: helper.FilterTypeString},
	"type":    {TypeFunc: helper.FilterTypeString},
	"version": {TypeFunc: helper.FilterTypeString},

	"allow_list":     {TypeFunc: helper.FilterTypeString},
	"cluster_size":   {TypeFunc: helper.FilterTypeInt},
	"created":        {TypeFunc: helper.FilterTypeString},
	"encrypted":      {TypeFunc: helper.FilterTypeBool},
	"host_primary":   {TypeFunc: helper.FilterTypeString},
	"host_secondary": {TypeFunc: helper.FilterTypeString},
	"id":             {TypeFunc: helper.FilterTypeInt},
	"instance_uri":   {TypeFunc: helper.FilterTypeString},
	"updated":        {TypeFunc: helper.FilterTypeString},
	"label":          {TypeFunc: helper.FilterTypeString},
}

var dataSourceSchema = map[string]*schema.Schema{
	"latest": {
		Type:        schema.TypeBool,
		Description: "If true, only the latest created database will be returned.",
		Optional:    true,
		Default:     false,
	},
	"order_by": filterConfig.OrderBySchema(),
	"order":    filterConfig.OrderSchema(),
	"filter":   filterConfig.FilterSchema(),
	"databases": {
		Type:        schema.TypeList,
		Description: "The returned list of databases.",
		Computed:    true,
		Elem:        databaseSchema(),
	},
}

func databaseSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allow_list": {
				Type: schema.TypeList,
				Description: "A list of IP addresses that can access the Managed Database. " +
					"Each item can be a single IP address or a range in CIDR format.",
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"cluster_size": {
				Type:        schema.TypeInt,
				Description: "The number of Linode Instance nodes deployed to the Managed Database.",
				Computed:    true,
			},
			"created": {
				Type:        schema.TypeString,
				Description: "When this Managed Database was created.",
				Computed:    true,
			},
			"encrypted": {
				Type:        schema.TypeBool,
				Description: "Whether the Managed Databases is encrypted.",
				Computed:    true,
			},
			"engine": {
				Type:        schema.TypeString,
				Description: "The Managed Database engine type.",
				Computed:    true,
			},
			"host_primary": {
				Type:        schema.TypeString,
				Description: "The primary host for the Managed Database.",
				Computed:    true,
			},
			"host_secondary": {
				Type:        schema.TypeString,
				Description: "The secondary/private network host for the Managed Database.",
				Computed:    true,
			},
			"id": {
				Type:        schema.TypeInt,
				Description: "A unique ID that can be used to identify and reference the Managed Database.",
				Computed:    true,
			},
			"instance_uri": {
				Type:        schema.TypeString,
				Description: "The API route for the database instance.",
				Computed:    true,
			},
			"label": {
				Type:        schema.TypeString,
				Description: "A unique, user-defined string referring to the Managed Database.",
				Computed:    true,
			},
			"region": {
				Type:        schema.TypeString,
				Description: "The Region ID for the Managed Database.",
				Computed:    true,
			},
			"replication_type": {
				Type:        schema.TypeString,
				Description: "The replication method used for the Managed Database.",
				Computed:    true,
			},
			"ssl_connection": {
				Type:        schema.TypeBool,
				Description: "Whether to require SSL credentials to establish a connection to the Managed Database.",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "The operating status of the Managed Database.",
				Computed:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The Linode Instance type used by the Managed Database for its nodes.",
				Computed:    true,
			},
			"updated": {
				Type:        schema.TypeString,
				Description: "When this Managed Database was last updated.",
				Computed:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "The Managed Database engine version.",
				Computed:    true,
			},
		},
	}
}
