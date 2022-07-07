package databasebackups

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

var filterConfig = helper.FilterConfig{
	"created": {APIFilterable: true, TypeFunc: helper.FilterTypeString},
	"type":    {APIFilterable: true, TypeFunc: helper.FilterTypeString},

	"id":    {APIFilterable: false, TypeFunc: helper.FilterTypeInt},
	"label": {APIFilterable: false, TypeFunc: helper.FilterTypeString},
}

var dataSourceSchema = map[string]*schema.Schema{
	"database_id": {
		Type:        schema.TypeInt,
		Description: "The ID of the Managed Database.",
		Required:    true,
	},
	"database_type": {
		Type:        schema.TypeString,
		Description: "The type of the Managed Database",
		Required:    true,
		ValidateDiagFunc: validation.ToDiagFunc(
			validation.StringInSlice(helper.ValidDatabaseTypes, true)),
	},

	"latest": {
		Type:        schema.TypeBool,
		Description: "If true, only the latest backup will be returned.",
		Optional:    true,
		Default:     false,
	},
	"order_by": filterConfig.OrderBySchema(),
	"order":    filterConfig.OrderSchema(),
	"filter":   filterConfig.FilterSchema(),
	"backups": {
		Type:        schema.TypeList,
		Description: "The returned list of backups.",
		Computed:    true,
		Elem:        backupsSchema(),
	},
}

func backupsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"created": {
				Type: schema.TypeString,
				Description: "A time value given in a combined date and time format that represents when the " +
					"database backup was created.",
				Computed: true,
			},
			"id": {
				Type:        schema.TypeInt,
				Description: "The ID of the database backup object.",
				Computed:    true,
			},
			"label": {
				Type:        schema.TypeString,
				Description: "The database backup’s label, for display purposes only.",
				Computed:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The type of database backup, determined by how the backup was created.",
				Computed:    true,
			},
		},
	}
}
