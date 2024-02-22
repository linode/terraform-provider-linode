package firewalldevice

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique ID that represents the firewall device in the Terraform state.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
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
			// Planned breaking change: Adding RFC3339 custom type in Linode provider v3.
			// Previous SDKv2 resource didn't format the time into RFC3339 format.
			// Starting Linode provider v2.12, all time string will be converted to
			// RFC3339 format in the state for this resource. Once the time strings
			// are converted to RFC3339 format by a Linode provider >= 2.12 and < 3,
			// RFC3339 custom type with validating logic may be safely added to this
			// attribute in v3.
			Description: "When this Firewall Device was created.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated": schema.StringAttribute{
			// Planned breaking change: add RFC3339 type, similar to 'created' field above
			Description: "When this Firewall Device was updated.",
			Computed:    true,
		},
	},
}
