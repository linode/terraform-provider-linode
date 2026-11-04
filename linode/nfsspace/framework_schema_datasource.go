package nfsspace

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var SpaceAttributes = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description: "The unique numeric ID of the NFS Space.",
		Required:    true,
	},
	"label": schema.StringAttribute{
		Description: "The customer-provided label of the NFS Space.",
		Computed:    true,
	},
	"description": schema.StringAttribute{
		Description: "An optional free-text description of this NFS Space.",
		Computed:    true,
	},
	"status": schema.StringAttribute{
		Description: "The current lifecycle status of the NFS Space.",
		Computed:    true,
	},
	"max_capacity_bytes": schema.Int64Attribute{
		Description: "The maximum total capacity available to this NFS Space across all of its " +
			"filesystems, in bytes.",
		Computed: true,
	},
	"used_capacity_bytes": schema.Int64Attribute{
		Description: "The current used capacity across the NFS Space, in bytes.",
		Computed:    true,
	},
	"created": schema.StringAttribute{
		Description: "When this NFS Space was created.",
		Computed:    true,
		CustomType:  timetypes.RFC3339Type{},
	},
	"updated": schema.StringAttribute{
		Description: "When this NFS Space was last updated.",
		Computed:    true,
		CustomType:  timetypes.RFC3339Type{},
	},
	"tags": schema.SetAttribute{
		Description: "Customer-defined metadata tags for grouping and filtering.",
		Computed:    true,
		ElementType: types.StringType,
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: SpaceAttributes,
}
