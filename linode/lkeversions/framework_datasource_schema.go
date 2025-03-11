package lkeversions

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/linode/terraform-provider-linode/v2/linode/lkeversion"
)

var lkeVersionSchema = schema.NestedBlockObject{
	Attributes: lkeversion.Attributes,
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Unique identification field for this list of LKE Versions.",
			Computed:    true,
		},
		"tier": schema.StringAttribute{
			Description: "The tier of the LKE versions, either standard or enterprise.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.OneOf("standard", "enterprise"),
			},
		},
	},
	Blocks: map[string]schema.Block{
		"versions": schema.ListNestedBlock{
			Description:  "The Kubernetes version numbers available for deployment to a Kubernetes cluster.",
			NestedObject: lkeVersionSchema,
		},
	},
}
