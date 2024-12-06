package networkingipassignment

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "The ID of the IP assignment operation.",
		},
		"region": schema.StringAttribute{
			Required:    true,
			Description: "The region for the IP assignments.",
		},
		"assignments": schema.ListAttribute{
			ElementType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"address":   types.StringType,
					"linode_id": types.Int64Type,
				},
			},
			Optional: true,
		},
	},
}
