package accountsettings

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The email of the current account.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"backups_enabled": schema.BoolAttribute{
			Description: "Account-wide backups default.",
			Computed:    true,
			Optional:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"network_helper": schema.BoolAttribute{
			Description: "Enables network helper across all users by default for new Linodes and Linode Configs.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"longview_subscription": schema.StringAttribute{
			Description: "The Longview Pro tier you are currently subscribed to.",
			Computed:    true,
			Optional:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"interfaces_for_new_linodes": schema.StringAttribute{
			Description: "Type of interfaces for new Linode instances.",
			Computed:    true,
			Optional:    true,
			Validators: []validator.String{stringvalidator.OneOf(
				"legacy_config_only",
				"legacy_config_default_but_linode_allowed",
				"linode_default_but_legacy_config_allowed",
				"linode_only",
			)},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"managed": schema.BoolAttribute{
			Description: "Enables monitoring for connectivity, response, and total request time.",
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"object_storage": schema.StringAttribute{
			Description: "A string describing the status of this account's Object Storage service enrollment.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	},
}
