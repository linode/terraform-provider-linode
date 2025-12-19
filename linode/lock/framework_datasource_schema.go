package lock

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var frameworkDatasourceSchema = schema.Schema{
	Description: "Provides information about a Linode Lock.",
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The unique ID of the Lock.",
			Required:    true,
		},
		"entity_id": schema.Int64Attribute{
			Description: "The ID of the locked entity.",
			Computed:    true,
		},
		"entity_type": schema.StringAttribute{
			Description: "The type of the locked entity.",
			Computed:    true,
		},
		"lock_type": schema.StringAttribute{
			Description: "The type of lock. Possible values: 'cannot_delete' (prevents deletion, rebuild, and transfer) or 'cannot_delete_with_subresources' (also prevents deletion of subresources).",
			Computed:    true,
		},
		"entity_label": schema.StringAttribute{
			Description: "The label of the locked entity.",
			Computed:    true,
		},
		"entity_url": schema.StringAttribute{
			Description: "The URL of the locked entity.",
			Computed:    true,
		},
	},
}
