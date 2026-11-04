package natgateway

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// updatedPlanModifier keeps the prior `updated` value unless a user-settable attribute changes.
// `updated` also changes when another resource (e.g. linode_vpc_subnet) modifies the NAT Gateway,
// which would otherwise show as a change in every plan.
type updatedPlanModifier struct{}

func (updatedPlanModifier) Description(_ context.Context) string {
	return "Keeps the prior value unless a user-settable attribute changes."
}

func (m updatedPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (updatedPlanModifier) PlanModifyString(
	ctx context.Context,
	req planmodifier.StringRequest,
	resp *planmodifier.StringResponse,
) {
	// Keep the framework's unknown value on create and destroy.
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() || req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}

	var plan, state, config NATGatewayModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Unknown addresses are unchanged only when the user omitted them.
	addressesUnchanged := plan.Addresses.Equal(state.Addresses) ||
		(plan.Addresses.IsUnknown() && config.Addresses.IsNull())

	if plan.Region.Equal(state.Region) &&
		plan.Label.Equal(state.Label) &&
		plan.DefaultPortsPerInterface.Equal(state.DefaultPortsPerInterface) &&
		plan.UseAutoscaling.Equal(state.UseAutoscaling) &&
		plan.VPCSubnetID.Equal(state.VPCSubnetID) &&
		addressesUnchanged {
		resp.PlanValue = req.StateValue
	}
}

// portsetPlanModifier keeps the prior portset value unless addresses change.
// Portset values also change when a VPC subnet attaches or detaches, which
// must not show as a diff. Address changes can alter them, so those plans
// leave the value unknown and Update stores the refreshed API value.
type portsetPlanModifier struct{}

func (portsetPlanModifier) Description(_ context.Context) string {
	return "Keeps the prior value unless addresses change."
}

func (m portsetPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (portsetPlanModifier) PlanModifyInt64(
	ctx context.Context,
	req planmodifier.Int64Request,
	resp *planmodifier.Int64Response,
) {
	// Keep the framework's unknown value on create and destroy.
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() || req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}

	var planAddrs, stateAddrs, configAddrs types.Set
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("addresses"), &planAddrs)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("addresses"), &stateAddrs)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("addresses"), &configAddrs)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Unknown addresses are unchanged only when the user omitted them.
	addressesUnchanged := planAddrs.Equal(stateAddrs) ||
		(planAddrs.IsUnknown() && configAddrs.IsNull())

	if addressesUnchanged {
		resp.PlanValue = req.StateValue
	}
}

var vpcSubnet = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique ID of the attached VPC subnet.",
			Computed:    true,
		},
		"type": schema.StringAttribute{
			Description: "The type of the attached vpc_subnet.",
			Computed:    true,
		},
		"url": schema.StringAttribute{
			Description: "The API URL of the attached VPC subnet.",
			Computed:    true,
		},
		"label": schema.StringAttribute{
			Description: "The label of the VPC subnet.",
			Computed:    true,
		},
		"vpc_id": schema.StringAttribute{
			Description: "The unique ID of the parent VPC.",
			Computed:    true,
		},
		"vpc_label": schema.StringAttribute{
			Description: "The label of the parent VPC.",
			Computed:    true,
		},
	},
}

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Computed:    true,
			Description: "The unique ID of the NAT Gateway.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"region": schema.StringAttribute{
			Required:    true,
			Description: "The region where the NAT Gateway will be created.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"label": schema.StringAttribute{
			Required:    true,
			Description: "The label for the NAT Gateway.",
		},
		"default_ports_per_interface": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "The default number of ports per interface for the NAT Gateway.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"use_autoscaling": schema.BoolAttribute{
			Optional: true,
			Description: "Whether the NAT Gateway automatically scales its public IPv4 addresses. " +
				"If omitted, defaults to true when addresses is empty, otherwise false. Create-only and Write-only.",
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.RequiresReplace(),
			},
		},
		"vpc_subnet_id": schema.Int64Attribute{
			Optional:    true,
			Description: "The ID of the VPC subnet for this NAT Gateway. Create-only.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"address_autoscale_max": schema.Int64Attribute{
			Computed:    true,
			Description: "The maximum number of auto-allocated IP addresses allowed for this NAT Gateway instance. Reflects the state of use_autoscaling.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"portset_assignments": schema.Int64Attribute{
			Computed:    true,
			Description: "The number of portsets currently assigned to interfaces behind this NAT Gateway, excluding system-reserved ports.",
			PlanModifiers: []planmodifier.Int64{
				portsetPlanModifier{},
			},
		},
		"portset_capacity": schema.Int64Attribute{
			Computed:    true,
			Description: "The total number of portsets available to allocate across all current and potential autoscaled public IP addresses assigned to this NAT Gateway, excluding system-reserved ports.",
			PlanModifiers: []planmodifier.Int64{
				portsetPlanModifier{},
			},
		},
		"vpc_subnet": schema.SingleNestedAttribute{
			Computed:    true,
			Attributes:  vpcSubnet.Attributes,
			Description: "Information about the VPC subnet attached to this NAT Gateway, if there is one attached.",
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
		},
		"created": schema.StringAttribute{
			Description: "When this NAT Gateway was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated": schema.StringAttribute{
			Description: "When this NAT Gateway was last updated.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			PlanModifiers: []planmodifier.String{
				updatedPlanModifier{},
			},
		},
		"addresses": schema.SetNestedAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The public IP addresses assigned to this NAT Gateway. Must be empty when use_autoscaling is true.",
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"address": schema.StringAttribute{
						Description: "The public IP address assigned to this NAT Gateway.",
						Required:    true,
					},
				},
			},
		},
	},
}
