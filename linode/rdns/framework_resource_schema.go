package rdns

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"address": schema.StringAttribute{
			Description: "The public Linode IPv4 or IPv6 address to operate on.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				helper.NewIPStringValidator(),
			},
		},
		"rdns": schema.StringAttribute{
			Description: "The reverse DNS assigned to this address. For public IPv4 addresses, this will be set " +
				"to a default value provided by Linode if not explicitly set.",
			Required: true,
			Validators: []validator.String{
				helper.NewStringLengthValidator(3, 254),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"wait_for_available": schema.BoolAttribute{
			Description: "If true, the RDNS assignment will be retried within the operation timeout period.",
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(false),
		},
		"id": schema.StringAttribute{
			Description: "Unique identification field for this RDNS Resource. " +
				"The public Linode IPv4 or IPv6 address to operate on. ",
			Computed: true,
		},
	},
}
