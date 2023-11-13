package token

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/helper/customtypes"
)

const (
	RequireReplacementWhenExpiryChangedDescription = "Requiring token recreation if the " +
		"Expiry time semantically changed."
	RequireReplacementWhenScopesChangedDescription = "Requiring token recreation if the " +
		"OAuth 2.0 scopes changed semantically. For example, from 'linodes:read_only lke:read_only' " +
		"to 'linodes:read_write lke:read_only' is a semantically change, but from " +
		"'linodes:read_only lke:read_only' to 'lke:read_only linodes:read_only' isn't."
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"label": schema.StringAttribute{
			Description: "The label of the Linode Token.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 100),
			},
		},
		"scopes": schema.StringAttribute{
			Description: "The scopes this token was created with. These define what parts of the Account the " +
				"token can be used to access. Many command-line tools, such as the Linode CLI, require tokens with " +
				"access to *. Tokens with more restrictive scopes are generally more secure. Multiple scopes are " +
				"separated by a space character (e.g., \"databases:read_only events:read_only\"). You can find the " +
				"list of available scopes on Linode API docs site, https://www.linode.com/docs/api#oauth-reference",
			Required: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplaceIf(
					func(
						ctx context.Context,
						sr planmodifier.StringRequest,
						rrifr *stringplanmodifier.RequiresReplaceIfFuncResponse,
					) {
						rrifr.RequiresReplace = !helper.CompareScopes(
							sr.PlanValue.ValueString(),
							sr.StateValue.ValueString(),
						)
					},
					RequireReplacementWhenExpiryChangedDescription,
					RequireReplacementWhenExpiryChangedDescription,
				),
			},
			CustomType: customtypes.LinodeScopesStringType{},
		},
		"expiry": schema.StringAttribute{
			Description: "When this token will expire. Personal Access Tokens cannot be renewed, so after " +
				"this time the token will be completely unusable and a new token will need to be generated. Tokens " +
				"may be created with 'null' as their expiry and will never expire unless revoked. Format: " +
				helper.TIME_FORMAT,
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
				stringplanmodifier.RequiresReplaceIf(
					func(
						ctx context.Context,
						sr planmodifier.StringRequest,
						rrifr *stringplanmodifier.RequiresReplaceIfFuncResponse,
					) {
						rrifr.RequiresReplace = !helper.CompareTimeStrings(
							sr.PlanValue.ValueString(),
							sr.StateValue.ValueString(),
							time.RFC3339,
						)
					},
					RequireReplacementWhenScopesChangedDescription,
					RequireReplacementWhenScopesChangedDescription,
				),
			},
			CustomType: timetypes.RFC3339Type{},
		},
		"created": schema.StringAttribute{
			Description: "The date and time this token was created.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			CustomType: timetypes.RFC3339Type{},
		},
		"token": schema.StringAttribute{
			Sensitive:   true,
			Description: "The token used to access the API.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"id": schema.StringAttribute{
			Description: "The ID of the token.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	},
}
