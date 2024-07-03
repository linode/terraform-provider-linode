package childaccount

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v2/linode/account"
)

func dataSourceSchema() *schema.Schema {
	result := account.DataSourceSchema()

	// This is a bit evil but allows us to avoid redefining the entire
	// account schema.
	result.Attributes["euuid"] = schema.StringAttribute{
		Description: "The unique ID of this Account.",
		Required:    true,
	}

	return &result
}
