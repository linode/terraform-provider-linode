package accountlogin

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var frameworkDataSourceSchema = schema.Schema{
	Attributes: Attributes,
}

var Attributes = map[string]schema.Attribute{
	"datetime": schema.StringAttribute{
		Description: "The time when the login was initiated.",
		Computed:    true,
	},
	"id": schema.Int64Attribute{
		Description: "The unique ID of this login object.",
		Required:    true,
	},
	"ip": schema.StringAttribute{
		Description: "The remote IP address that requested the login.",
		Computed:    true,
	},
	"restricted": schema.BoolAttribute{
		Description: "True if the User that was logged into was a restricted User, false otherwise.",
		Computed:    true,
	},
	"username": schema.StringAttribute{
		Description: "The username of the User that was logged into.",
		Computed:    true,
	},
	"status": schema.StringAttribute{
		Description: "Whether the login attempt succeeded or failed.",
		Computed:    true,
	},
}
