package lkeversions

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var lkeVersionObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id": types.StringType,
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"versions": schema.ListAttribute{
			Description: "The Kubernetes version numbers available for deployment to a Kubernetes cluster.",
			Computed:    true,
			ElementType: lkeVersionObjectType,
		},
		"id": schema.StringAttribute{
			Description: "Unique identification field for this list of LKE Versions.",
			Computed:    true,
		},
	},
}
