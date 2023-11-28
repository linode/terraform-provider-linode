package databasebackups

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"created": {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"type":    {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"label":   {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"id":      {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeInt},
}

var backupSchema = schema.NestedBlockObject{
	Attributes: map[string]schema.Attribute{
		"created": schema.StringAttribute{
			Description: "A time value given in a combined date and time format that represents when the " +
				"database backup was created.",
			Computed: true,
		},
		"id": schema.Int64Attribute{
			Description: "The ID of the database backup object.",
			Computed:    true,
		},
		"label": schema.StringAttribute{
			Description: "The database backup’s label, for display purposes only.",
			Computed:    true,
		},
		"type": schema.StringAttribute{
			Description: "The type of database backup, determined by how the backup was created.",
			Computed:    true,
		},
	},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"database_id": schema.Int64Attribute{
			Description: "The ID of the Managed Database.",
			Required:    true,
		},
		"id": schema.Int64Attribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"database_type": schema.StringAttribute{
			Description: "The type of the Managed Database",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.OneOf(helper.ValidDatabaseTypes...),
			},
		},
		"latest": schema.BoolAttribute{
			Description: "If true, only the latest engine version will be returned.",
			Optional:    true,
		},
		"order_by": filterConfig.OrderBySchema(),
		"order":    filterConfig.OrderSchema(),
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
		"backups": schema.ListNestedBlock{
			Description:  "The returned list of backups.",
			NestedObject: backupSchema,
		},
	},
}
