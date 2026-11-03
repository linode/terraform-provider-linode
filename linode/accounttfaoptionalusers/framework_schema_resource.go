package accounttfaoptionalusers

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkResourceSchema = schema.Schema{
	Description: "Manages the complete account-wide set of users exempt from two-factor authentication enforcement. " +
		"Only one resource should manage this set per account. Deleting the resource makes no API calls and leaves exemptions unchanged; " +
		"apply an empty set to remove them.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The singleton identifier: account.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"usernames": schema.SetAttribute{
			Description: "The complete set of exempt usernames, regardless of the account enforcement policy. " +
				"An empty set removes all exemptions.",
			Required:    true,
			ElementType: types.StringType,
			Validators: []validator.Set{
				setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
			},
		},
	},
}
