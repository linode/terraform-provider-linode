package token

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"label": schema.StringAttribute{
			Description: "The label of the Linode Token.",
			Optional:    true,
		},
		"scopes": schema.StringAttribute{
			Description: "The scopes this token was created with. These define what parts of the Account the " +
				"token can be used to access. Many command-line tools, such as the Linode CLI, require tokens with " +
				"access to *. Tokens with more restrictive scopes are generally more secure. Multiple scopes are " +
				"separated by a space character (e.g., \"databases:read_only events:read_only\"). You can find the " +
				"list of available scopes on Linode API docs site, https://www.linode.com/docs/api#oauth-reference",
			Required: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"expiry": schema.StringAttribute{
			Description: "When this token will expire. Personal Access Tokens cannot be renewed, so after " +
				"this time the token will be completely unusable and a new token will need to be generated. Tokens " +
				"may be created with 'null' as their expiry and will never expire unless revoked. Format: " +
				time.RFC3339,
			Optional: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				helper.NewDateTimeStringValidator(time.RFC3339),
			},
		},
		"created": schema.StringAttribute{
			Description: "The date and time this token was created.",
			Computed:    true,
		},
		"token": schema.StringAttribute{
			Sensitive:   true,
			Description: "The token used to access the API.",
			Computed:    true,
		},
		"id": schema.StringAttribute{
			Description: "The ID of the token.",
			Computed:    true,
		},
	},
}
