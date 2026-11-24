package nfsspaceaccesspolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkDatasourceSchema = schema.Schema{
	Description: "Provides information about the access policy of a Linode NFS Space.",
	Attributes: map[string]schema.Attribute{
		"space_id": schema.Int64Attribute{
			Description: "The ID of the NFS Space this access policy belongs to.",
			Required:    true,
		},
		"id": schema.StringAttribute{
			Description: "The unique ID of this access policy. Equal to the ID of the NFS Space.",
			Computed:    true,
		},
		"label": schema.StringAttribute{
			Description: "The descriptive label of this access policy.",
			Computed:    true,
		},
		"enabled": schema.BoolAttribute{
			Description: "Whether this access policy is active. When false, the NFS Space is inaccessible.",
			Computed:    true,
		},
		"vpcs": schema.SetNestedAttribute{
			Description: "The VPCs (and optionally specific subnets) permitted to access this NFS Space.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Description: "The ID of the VPC.",
						Computed:    true,
					},
					"subnets": schema.SetAttribute{
						Description: "The IDs of the permitted subnets within the VPC. Null when all " +
							"subnets in the VPC are permitted.",
						Computed:    true,
						ElementType: types.Int64Type,
					},
				},
			},
		},
		"mtls_ca_cert": schema.StringAttribute{
			Description: "The PEM-encoded CA certificate used to validate mTLS client certificates.",
			Computed:    true,
		},
		"mtls_mode": schema.StringAttribute{
			Description: "Whether mTLS client certificate presentation is enforced.",
			Computed:    true,
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
		},
		"created": schema.StringAttribute{
			Description: "When this access policy was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"updated": schema.StringAttribute{
			Description: "When this access policy was last updated.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
	},
}
