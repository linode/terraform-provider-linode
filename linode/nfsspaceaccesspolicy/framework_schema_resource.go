package nfsspaceaccesspolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
)

var frameworkResourceSchema = schema.Schema{
	Description: "Manages the access policy of a Linode NFS Space. The access policy controls which " +
		"VPCs can reach the NFS Space and the optional mTLS configuration, and applies to all " +
		"filesystems in the NFS Space.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique ID of this resource. Equal to the ID of the NFS Space.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"space_id": schema.Int64Attribute{
			Description: "The ID of the NFS Space this access policy belongs to.",
			Required:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
		},
		"label": schema.StringAttribute{
			Description: "The descriptive label of this access policy.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"enabled": schema.BoolAttribute{
			Description: "Whether this access policy is active.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"vpcs": schema.SetNestedAttribute{
			Description: "The VPCs (and optionally specific subnets) permitted to access this NFS Space.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Description: "The ID of the VPC.",
						Required:    true,
					},
					"subnets": schema.SetAttribute{
						Description: "The IDs of the subnets within the VPC to permit. When unset, " +
							"all subnets in the VPC are permitted.",
						Optional:    true,
						ElementType: types.Int64Type,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
					},
				},
			},
		},
		"mtls_ca_cert": schema.StringAttribute{
			Description: "The PEM-encoded CA certificate used to validate mTLS client certificates.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"mtls_mode": schema.StringAttribute{
			Description: "Whether mTLS client certificate presentation is enforced. " +
				"(`required`, `optional`, `disabled`)",
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf(
					string(linodego.NFSMTLSModeRequired),
					string(linodego.NFSMTLSModeOptional),
					string(linodego.NFSMTLSModeDisabled),
				),
			},
		},
		"vpc_acl": schema.ListNestedAttribute{
			Description: "The resolved VPC ACL entries, one per permitted VPC.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Description: "The ID of the VPC.",
						Computed:    true,
					},
					"label": schema.StringAttribute{
						Description: "The label of the VPC.",
						Computed:    true,
					},
					"url": schema.StringAttribute{
						Description: "The relative URL of the VPC.",
						Computed:    true,
					},
					"range": schema.StringAttribute{
						Description: "The resolved CIDR range of the entire VPC. Only set when no " +
							"specific subnets were requested.",
						Computed: true,
					},
					"subnets": schema.ListNestedAttribute{
						Description: "The resolved subnet entries within the VPC. Only populated " +
							"when specific subnets were requested.",
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.Int64Attribute{
									Description: "The ID of the subnet.",
									Computed:    true,
								},
								"label": schema.StringAttribute{
									Description: "The label of the subnet.",
									Computed:    true,
								},
								"url": schema.StringAttribute{
									Description: "The relative URL of the subnet.",
									Computed:    true,
								},
								"range": schema.StringAttribute{
									Description: "The resolved CIDR range of the subnet.",
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
		"status": schema.StringAttribute{
			Description: "The current lifecycle status of the access policy.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"created": schema.StringAttribute{
			Description: "When this access policy was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated": schema.StringAttribute{
			Description: "When this access policy was last updated.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
	},
}
