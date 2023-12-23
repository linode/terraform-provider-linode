package firewalldevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The unique ID that represents the firewall device in the Terraform state.",
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"firewall_id": schema.Int64Attribute{
			Description: "The ID of the Firewall to access.",
			Required:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
		},
		"entity_id": schema.Int64Attribute{
			Description: "The ID of the entity to create a Firewall device for.",
			Required:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
		},
		"entity_type": schema.StringAttribute{
			Description: "The type of the entity to create a Firewall device for.",
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("linode"),
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"created": schema.StringAttribute{
			Description: "When this Firewall Device was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated": schema.StringAttribute{
			Description: "When this Firewall Device was updated.",
			CustomType:  timetypes.RFC3339Type{},
			Computed:    true,
		},
	},
}
