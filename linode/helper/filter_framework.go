package helper

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var FrameworkFilterSchema = schema.SetNestedBlock{
	NestedObject: schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the attribute to filter on.",
			},
			"values": schema.SetAttribute{
				Required:    true,
				Description: "The value(s) to be used in the filter.",
				ElementType: types.StringType,
			},
			"match_by": schema.StringAttribute{
				Optional:    true,
				Description: "The value(s) to be used in the filter.",
			},
		},
	},
}

// FrameworkFilterModel describes the Terraform resource data model to match the
// resource schema.
type FrameworkFilterModel struct {
	Name    types.String       `tfsdk:"name"`
	Values  basetypes.SetValue `tfsdk:"values"`
	MatchBy types.String       `tfsdk:"match_by"`
}

type FrameworkFiltersType types.Set

type FrameworkFilterAttribute struct {
	APIFilterable bool
	Type          attr.Type
}

type FrameworkFilterConfig map[string]FrameworkFilterAttribute

func (f FrameworkFilterConfig) ConstructFilterString(ctx context.Context, filterSet types.Set) (string, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	var filterObjects []types.Object

	resultMap := make(map[string]any)

	diagnostics.Append(
		filterSet.ElementsAs(ctx, &filterObjects, false)...,
	)
	if diagnostics.HasError() {
		return "", diagnostics
	}

	for _, filter := range filterObjects {
		filterAttrs := filter.Attributes()
		filterFieldName := filterAttrs["name"].String()

		if _, ok := resultMap[filterFieldName]; !ok {

		}
	}

	return "cool", diagnostics
}
