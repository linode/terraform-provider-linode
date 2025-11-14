package databaseshared

import (
	"context"

	dataSourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type PrivateNetworkModel struct {
	VPCID        types.Int64 `tfsdk:"vpc_id"`
	SubnetID     types.Int64 `tfsdk:"subnet_id"`
	PublicAccess types.Bool  `tfsdk:"public_access"`
}

func (m PrivateNetworkModel) ToLinodego(d diag.Diagnostics) *linodego.DatabasePrivateNetwork {
	return &linodego.DatabasePrivateNetwork{
		VPCID:        helper.FrameworkSafeInt64ToInt(m.VPCID.ValueInt64(), &d),
		SubnetID:     helper.FrameworkSafeInt64ToInt(m.SubnetID.ValueInt64(), &d),
		PublicAccess: m.PublicAccess.ValueBool(),
	}
}

var ResourceAttributePrivateNetwork = resourceSchema.SingleNestedAttribute{
	Description: "Restricts access to this database using a virtual private cloud (VPC) " +
		"that you've configured in the region where the database will live.",
	Optional: true,
	PlanModifiers: []planmodifier.Object{
		objectplanmodifier.UseStateForUnknown(),
	},
	Attributes: map[string]resourceSchema.Attribute{
		"vpc_id": resourceSchema.Int64Attribute{
			Description: " The ID of the virtual private cloud (VPC) " +
				"to restrict access to this database using.",
			Required: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"subnet_id": resourceSchema.Int64Attribute{
			Description: "The ID of the VPC subnet to restrict access to this database using.",
			Required:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"public_access": resourceSchema.BoolAttribute{
			Description: "Set to `true` to allow clients outside of the VPC to " +
				"connect to the database using a public IP address.",
			Optional: true,
			Computed: true,
			Default:  booldefault.StaticBool(false),
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
	},
}

var DataSourceAttributePrivateNetwork = dataSourceSchema.SingleNestedAttribute{
	Description: "Restricts access to this database using a virtual private cloud (VPC) " +
		"that you've configured in the region where the database will live.",
	Computed: true,
	Attributes: map[string]dataSourceSchema.Attribute{
		"vpc_id": dataSourceSchema.Int64Attribute{
			Description: "The ID of the virtual private cloud (VPC) " +
				"to restrict access to this database using.",
			Computed: true,
		},
		"subnet_id": dataSourceSchema.Int64Attribute{
			Description: "The ID of the VPC subnet to restrict access to this database using.",
			Computed:    true,
		},
		"public_access": dataSourceSchema.BoolAttribute{
			Description: "If true, clients outside of the VPC can " +
				"connect to the database using a public IP address.",
			Computed: true,
		},
	},
}

var ObjectTypePrivateNetwork = ResourceAttributePrivateNetwork.GetType().(types.ObjectType)

func FlattenPrivateNetwork(
	ctx context.Context,
	privateNetwork linodego.DatabasePrivateNetwork,
) (types.Object, diag.Diagnostics) {
	return types.ObjectValueFrom(
		ctx,
		ObjectTypePrivateNetwork.AttrTypes,
		&PrivateNetworkModel{
			VPCID:        types.Int64Value(int64(privateNetwork.VPCID)),
			SubnetID:     types.Int64Value(int64(privateNetwork.SubnetID)),
			PublicAccess: types.BoolValue(privateNetwork.PublicAccess),
		},
	)
}
