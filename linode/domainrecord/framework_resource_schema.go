package domainrecord

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	linodeplanmodifiers "github.com/linode/terraform-provider-linode/v2/linode/helper/planmodifiers"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the domain record.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"domain_id": schema.Int64Attribute{
			Description: "The ID of the Domain to access.",
			Required:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
		},
		"name": schema.StringAttribute{
			Description: "The name of this Record. This field's actual usage depends " +
				"on the type of record this represents. For A and AAAA records, this is " +
				"the subdomain being associated with an IP address. Generated for SRV records.",
			Optional: true,
			Computed: true, // This is true for SRV records
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.String{
				stringvalidator.LengthBetween(0, 100),
			},
		},
		"record_type": schema.StringAttribute{
			Description: "The type of Record this is in the DNS system. " +
				"For example, A records associate a domain name with an IPv4 " +
				"address, and AAAA records associate a domain name with an IPv6 address.",
			Required: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf(
					"A", "AAAA", "NS", "MX", "CNAME", "TXT", "SRV", "PTR", "CAA",
				),
			},
		},
		"ttl_sec": schema.Int64Attribute{
			Description: "'Time to Live' - the amount of time in seconds that this " +
				"Domain's records may be cached by resolvers or other domain servers. " +
				"Valid values are 30, 120, 300, 3600, 7200, 14400, 28800, 57600, 86400, " +
				"172800, 345600, 604800, 1209600, and 2419200 - any other value will be " +
				"rounded to the nearest valid value.",
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.Int64{
				linodeplanmodifiers.DomainRecordTTLUseStateIfPlanCanBeRoundedToState(),
			},
		},
		"target": schema.StringAttribute{
			Description: "The target for this Record. This field's actual usage depends " +
				"on the type of record this represents. For A and AAAA records, this is " +
				"the address the named Domain should resolve to.",
			Required: true,
			PlanModifiers: []planmodifier.String{
				linodeplanmodifiers.DomainRecordTargetUseStateIfSematicEquals(),
			},
		},
		"priority": schema.Int64Attribute{
			Description: "The priority of the target host. Lower values are preferred.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Int64{
				int64validator.Between(0, 255),
			},
		},
		"protocol": schema.StringAttribute{
			Description: "The protocol this Record's service communicates with. " +
				"Only valid for SRV records.",
			Optional: true,
		},
		"service": schema.StringAttribute{
			Description: "The service this Record identified. Only valid for SRV records.",
			Optional:    true,
		},
		"tag": schema.StringAttribute{
			Description: "The tag portion of a CAA record. " +
				"It is invalid to set this on other record types.",
			Optional: true,
		},
		"port": schema.Int64Attribute{
			Description: "The port this Record points to.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Optional: true,
			Computed: true,
		},
		"weight": schema.Int64Attribute{
			Description: "The relative weight of this Record. Higher values are preferred.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Optional: true,
			Computed: true,
		},
	},
}
