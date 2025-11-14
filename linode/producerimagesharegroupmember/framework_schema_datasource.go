package producerimagesharegroupmember

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var Attributes = map[string]schema.Attribute{
	"sharegroup_id": schema.Int64Attribute{
		Description: "The ID of the Image Share Group the member belongs to.",
		Required:    true,
	},
	"token_uuid": schema.StringAttribute{
		Description: "The UUID of member's token.",
		Required:    true,
	},
	"label": schema.StringAttribute{
		Description: "The label of the member.",
		Computed:    true,
	},
	"status": schema.StringAttribute{
		Description: "The status of the member.",
		Computed:    true,
	},
	"created": schema.StringAttribute{
		Description: "When this member was created.",
		Computed:    true,
		CustomType:  timetypes.RFC3339Type{},
	},
	"updated": schema.StringAttribute{
		Description: "When this member was last updated.",
		Computed:    true,
		CustomType:  timetypes.RFC3339Type{},
	},
	"expiry": schema.StringAttribute{
		Description: "When this member will expire.",
		Computed:    true,
		CustomType:  timetypes.RFC3339Type{},
	},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: Attributes,
}
