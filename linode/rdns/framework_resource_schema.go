package rdns

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/customtypes"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"address": schema.StringAttribute{
			Description: "The public Linode IPv4 or IPv6 address to operate on.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			CustomType: customtypes.IPAddrStringType{},
		},
		"rdns": schema.StringAttribute{
			Description: "The reverse DNS assigned to this address. For public IPv4 addresses, this will be set " +
				"to a default value provided by Linode if not explicitly set.",
			// Required: true,
			Optional: true,
			Computed: true,
			Validators: []validator.String{
				stringvalidator.LengthBetween(3, 254),
			},
		},
		"reserved": schema.BoolAttribute{
			Description: "Whether the IP address is reserved.",
			Optional:    true,
			Computed:    true,
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
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	},
}
