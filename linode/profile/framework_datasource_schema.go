package profile

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var referralObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"total":     types.Int64Type,
		"credit":    types.Float64Type,
		"completed": types.Int64Type,
		"pending":   types.Int64Type,
		"code":      types.StringType,
		"url":       types.StringType,
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"email": schema.StringAttribute{
			Description: "The profile email address. This address will be used for communication with Linode as necessary.",
			Computed:    true,
		},
		"timezone": schema.StringAttribute{
			Description: "The profile's preferred timezone. This is not used by the API, and is for the benefit of " +
				"clients only. All times the API returns are in UTC.",
			Computed: true,
		},
		"email_notifications": schema.BoolAttribute{
			Description: "If true, email notifications will be sent about account activity. If false, when false " +
				"business-critical communications may still be sent through email.",
			Computed: true,
		},
		"username": schema.StringAttribute{
			Description: "The username for logging in to Linode services.",
			Computed:    true,
		},
		"ip_whitelist_enabled": schema.BoolAttribute{
			Description: "If true, logins for the user will only be allowed from whitelisted IPs. " +
				"This setting is currently deprecated, and cannot be enabled.",
			Computed: true,
		},
		"lish_auth_method": schema.StringAttribute{
			Description: "The methods of authentication allowed when connecting via Lish. 'keys_only' is the most " +
				"secure with the intent to use Lish, and 'disabled' is recommended for users that will not use Lish at all.",
			Computed: true,
		},
		"authorized_keys": schema.ListAttribute{
			Description: "The list of SSH Keys authorized to use Lish for this user. This value is ignored if " +
				"lish_auth_method is 'disabled'.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"two_factor_auth": schema.BoolAttribute{
			Description: "If true, logins from untrusted computers will require Two Factor Authentication.",
			Computed:    true,
		},
		"restricted": schema.BoolAttribute{
			Description: "If true, the user has restrictions on what can be accessed on the Account.",
			Computed:    true,
		},
		"referrals": schema.ListAttribute{
			Description: "Credit Card information associated with this Account.",
			Computed:    true,
			ElementType: referralObjectType,
		},
		"id": schema.StringAttribute{
			Description: "Unique identification field for this datasource.",
			Computed:    true,
		},
	},
}
