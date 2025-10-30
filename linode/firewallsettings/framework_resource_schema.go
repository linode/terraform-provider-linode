package firewallsettings

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var FrameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "A unique identifier for this resource (UUID v7).",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"default_firewall_ids": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "The default firewall ID for a linode, nodebalancer, public_interface, or vpc_interface.",
			Attributes: map[string]schema.Attribute{
				"linode": schema.Int64Attribute{
					Description: "The Linode's default firewall.",
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.UseStateForUnknown(),
					},
					Optional: true,
					Computed: true,
				},
				"nodebalancer": schema.Int64Attribute{
					Description: "The NodeBalancer's default firewall.",
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.UseStateForUnknown(),
					},
					Optional: true,
					Computed: true,
				},
				"public_interface": schema.Int64Attribute{
					Description: "The public interface's default firewall.",
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.UseStateForUnknown(),
					},
					Optional: true,
					Computed: true,
				},
				"vpc_interface": schema.Int64Attribute{
					Description: "The VPC interface's default firewall.",
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.UseStateForUnknown(),
					},
					Optional: true,
					Computed: true,
				},
			},
		},
	},
}
