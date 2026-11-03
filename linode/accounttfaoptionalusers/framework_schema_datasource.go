package accounttfaoptionalusers

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkDataSourceSchema = schema.Schema{
	Description: "Lists all users exempt from two-factor authentication enforcement, regardless of the account policy.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The singleton identifier: account.",
			Computed:    true,
		},
		"usernames": schema.SetAttribute{
			Description: "The complete set of exempt usernames.",
			Computed:    true,
			ElementType: types.StringType,
		},
	},
}
