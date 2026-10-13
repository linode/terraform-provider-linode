package tag

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var tagObjectAttrTypes = map[string]attr.Type{
	"type": types.StringType,
	"id":   types.StringType,
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The label of this Tag.",
			Computed:    true,
		},
		"label": schema.StringAttribute{
			Description: "A label used to categorize resources. For display purposes only.",
			Required:    true,
		},
		"objects": schema.ListNestedAttribute{
			Description: "The objects associated with this tag.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "The type of the tagged object " +
							"(e.g. linode, domain, volume, nodebalancer, reserved_ipv4_address).",
						Computed: true,
					},
					"id": schema.StringAttribute{
						Description: "The ID (or address for reserved_ipv4_address) of the tagged object.",
						Computed:    true,
					},
				},
			},
		},
	},
}
