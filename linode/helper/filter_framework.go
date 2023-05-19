package helper

import (
	"context"
	"encoding/json"
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
	Name    types.String `tfsdk:"name"`
	Values  types.Set    `tfsdk:"values"`
	MatchBy types.String `tfsdk:"match_by"`
}

type FrameworkFiltersType types.Set

type FrameworkFilterAttribute struct {
	APIFilterable bool
	Type          attr.Type
}

type FrameworkFilterConfig map[string]FrameworkFilterAttribute

func ConstructFilterString(ctx context.Context, filterSet types.Set) (string, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	var filterObjects []types.Object

	diagnostics.Append(
		filterSet.ElementsAs(ctx, &filterObjects, false)...,
	)
	if diagnostics.HasError() {
		return "", diagnostics
	}

	rootFilter := make([]map[string]any, len(filterObjects))

	for filterIndex, filter := range filterObjects {
		// Parse the model
		var filterModel FrameworkFilterModel

		diagnostics.Append(
			filter.As(ctx, &filterModel, basetypes.ObjectAsOptions{})...,
		)
		if diagnostics.HasError() {
			return "", diagnostics
		}

		// Parse out the accepted values
		var filterFieldValues []types.String

		diagnostics.Append(
			filterModel.Values.ElementsAs(ctx, &filterFieldValues, false)...,
		)
		if diagnostics.HasError() {
			return "", diagnostics
		}

		// Get other attributes
		filterFieldName := filterModel.Name.ValueString()

		// Build the +or filter
		currentFilter := make([]map[string]any, len(filterFieldValues))

		for i, value := range filterFieldValues {
			currentFilter[i] = map[string]any{filterFieldName: value.ValueString()}
		}

		// Append to the root filter
		rootFilter[filterIndex] = map[string]any{
			"+or": currentFilter,
		}
	}

	resultFilter := map[string]any{
		"+and": rootFilter,
	}

	result, err := json.Marshal(resultFilter)
	if err != nil {
		diagnostics.AddError(
			"failed to marshal api filter",
			err.Error(),
		)
	}

	return string(result), diagnostics
}
