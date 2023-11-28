package nb

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

var frameworkResourceSchema = schema.Schema{
	Version: 1,
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:   "The unique ID of the Linode NodeBalancer.",
			Computed:      true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"label": schema.StringAttribute{
			Description: "The label of the Linode NodeBalancer.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.LengthBetween(3, 32),
			},
		},
		"region": schema.StringAttribute{
			Description: "The region where this NodeBalancer will be deployed.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
				stringplanmodifier.UseStateForUnknown(),
			},
			Default: stringdefault.StaticString("us-east"),
		},
		"client_conn_throttle": schema.Int64Attribute{
			Description: "Throttle connections per second (0-20). Set to 0 (zero) to disable throttling.",
			Validators: []validator.Int64{
				int64validator.Between(0, 20),
			},
			Optional: true,
			Computed: true,
			Default:  int64default.StaticInt64(0),
		},
		"firewall_id": schema.Int64Attribute{
			Description: "ID for the firewall you'd like to use with this NodeBalancer.",
			Optional:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
		},
		"hostname": schema.StringAttribute{
			Description:   "This NodeBalancer's hostname, ending with .nodebalancer.linode.com",
			Computed:      true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"ipv4": schema.StringAttribute{
			Description:   "The Public IPv4 Address of this NodeBalancer",
			Computed:      true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"ipv6": schema.StringAttribute{
			Description:   "The Public IPv6 Address of this NodeBalancer",
			Computed:      true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"created": schema.StringAttribute{
			Description:   "When this NodeBalancer was created.",
			Computed:      true,
			CustomType:    timetypes.RFC3339Type{},
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"updated": schema.StringAttribute{
			Description: "When this NodeBalancer was last updated.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"tags": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Default:     helper.EmptySetDefault(types.StringType),
			Computed:    true,
			Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
		},
		"transfer": schema.ListAttribute{
			Description: "Information about the amount of transfer this NodeBalancer has had so far this month.",
			Computed:    true,
			ElementType: TransferObjectType,
		},
	},
}

var resourceNodebalancerV0 = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique ID of the Linode NodeBalancer.",
			Computed:    true,
		},
		"label": schema.StringAttribute{
			Description: "The label of the Linode NodeBalancer.",
			Optional:    true,
		},
		"region": schema.StringAttribute{
			Description:   "The region where this NodeBalancer will be deployed.",
			Optional:      true,
			Computed:      true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			Default:       stringdefault.StaticString("us-east"),
		},
		"client_conn_throttle": schema.Int64Attribute{
			Description: "Throttle connections per second (0-20). Set to 0 (zero) to disable throttling.",
			Validators: []validator.Int64{
				int64validator.Between(0, 20),
			},
			Optional: true,
			Computed: true,
			Default:  int64default.StaticInt64(0),
		},
		"hostname": schema.StringAttribute{
			Description: "This NodeBalancer's hostname, ending with .nodebalancer.linode.com",
			Computed:    true,
		},
		"ipv4": schema.StringAttribute{
			Description: "The Public IPv4 Address of this NodeBalancer",
			Computed:    true,
		},
		"ipv6": schema.StringAttribute{
			Description: "The Public IPv6 Address of this NodeBalancer",
			Computed:    true,
		},
		"created": schema.StringAttribute{
			Description: "When this NodeBalancer was created.",
			Computed:    true,
		},
		"updated": schema.StringAttribute{
			Description: "When this NodeBalancer was last updated.",
			Computed:    true,
		},
		"tags": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
		},
		"transfer": schema.MapAttribute{
			ElementType: types.StringType,
			Computed:    true,
		},
	},
}
