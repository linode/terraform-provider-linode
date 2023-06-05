package instancetype

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var priceObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"hourly":  types.Float64Type,
		"monthly": types.Float64Type,
	},
}

var backupsObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"price": types.ListType{ElemType: priceObjectType},
	},
}

var addonsObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"backups": types.ListType{ElemType: backupsObjectType},
	},
}

var Attributes = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description: "The unique ID assigned to this Instance type.",
		Required:    true,
	},
	"label": schema.StringAttribute{
		Description: "The Linode Type's label is for display purposes only.",
		Computed:    true,
		Optional:    true,
	},
	"disk": schema.Int64Attribute{
		Description: "The Disk size, in MB, of the Linode Type.",
		Computed:    true,
	},
	"class": schema.StringAttribute{
		Description: "The class of the Linode Type. There are currently three classes of Linodes: nanode, " +
			"standard, highmem, dedicated",
		Computed: true,
	},
	"price": schema.ListAttribute{
		Description: "Cost in US dollars, broken down into hourly and monthly charges.",
		Computed:    true,
		ElementType: priceObjectType,
	},
	"addons": schema.ListAttribute{
		Description: "Information about the optional Backup service offered for Linodes.",
		Computed:    true,
		ElementType: addonsObjectType,
	},
	"network_out": schema.Int64Attribute{
		Description: "The Mbits outbound bandwidth allocation.",
		Computed:    true,
	},
	"memory": schema.Int64Attribute{
		Description: "Amount of RAM included in this Linode Type.",
		Computed:    true,
	},
	"transfer": schema.Int64Attribute{
		Description: "The monthly outbound transfer amount, in MB.",
		Computed:    true,
	},
	"vcpus": schema.Int64Attribute{
		Description: "The number of VCPU cores this Linode Type offers.",
		Computed:    true,
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: Attributes,
}
