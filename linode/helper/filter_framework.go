package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"reflect"
	"strconv"
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
	Name    types.String   `tfsdk:"name"`
	Values  []types.String `tfsdk:"values"`
	MatchBy types.String   `tfsdk:"match_by"`
}

type FrameworkFiltersType types.Set

type FrameworkFilterAttribute struct {
	APIFilterable bool
	Type          attr.Type
}

type FrameworkFilterConfig map[string]FrameworkFilterAttribute

func (f FrameworkFilterConfig) ParseFilterSet(
	ctx context.Context,
	filterSet types.Set,
) ([]FrameworkFilterModel, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	var result []FrameworkFilterModel

	// Parse out the set into an object list
	diagnostics.Append(
		filterSet.ElementsAs(ctx, &result, false)...,
	)
	if diagnostics.HasError() {
		return nil, diagnostics
	}

	return result, diagnostics
}

func (f FrameworkFilterConfig) ConstructFilterString(
	ctx context.Context,
	filterSet []FrameworkFilterModel,
) (string, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	rootFilter := make([]map[string]any, 0)

	for _, filter := range filterSet {
		// Get string attributes
		filterFieldName := filter.Name.ValueString()

		// Is this field filterable?
		filterFieldConfig, ok := f[filterFieldName]
		if !ok {
			diagnostics.AddError(
				"Attempted to filter on non-filterable field.",
				fmt.Sprintf("Attempted to filter on non-filterable field %s.", filterFieldName),
			)
			return "", diagnostics
		}

		// Skip if this field isn't API filterable
		if !filterFieldConfig.APIFilterable {
			continue
		}

		// We should only use API filters when matching on exact
		if !filter.MatchBy.IsNull() && filter.MatchBy.ValueString() != "exact" {
			continue
		}

		// Build the +or filter
		currentFilter := make([]map[string]any, len(filter.Values))

		for i, value := range filter.Values {
			currentFilter[i] = map[string]any{filterFieldName: value.ValueString()}
		}

		// Append to the root filter
		rootFilter = append(rootFilter, map[string]any{
			"+or": currentFilter,
		})
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

func (f FrameworkFilterConfig) ApplyLocalFiltering(
	ctx context.Context, filterSet []FrameworkFilterModel, data []map[string]any,
) ([]map[string]any, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	result := make([]map[string]any, 0)

	for _, elem := range data {
		match, d := f.MatchesFilter(ctx, filterSet, elem)
		diagnostics.Append(d...)
		if diagnostics.HasError() {
			return nil, diagnostics
		}

		// This element was filtered out
		if !match {
			continue
		}

		result = append(result, elem)
	}

	return result, diagnostics
}

func (f FrameworkFilterConfig) MatchesFilter(
	ctx context.Context,
	filterSet []FrameworkFilterModel,
	elem map[string]any,
) (bool, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	for _, filter := range filterSet {
		filterName := filter.Name.ValueString()

		// Skip if this field should be filtered at an API level
		if f[filterName].APIFilterable {
			continue
		}

		matchingField, ok := elem[filterName]
		if !ok {
			diagnostics.AddError(
				"Missing key",
				fmt.Sprintf("Flattened result map does not contain key %s.", filterName),
			)
		}

		match, d := f.checkFieldMatchesFilter(matchingField, filter)
		diagnostics.Append(d...)
		if diagnostics.HasError() {
			return match, diagnostics
		}

		// No match for this filter; return
		if !match {
			return false, diagnostics
		}
	}

	return true, diagnostics
}

func (f FrameworkFilterConfig) checkFieldMatchesFilter(
	field any,
	filter FrameworkFilterModel,
) (bool, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	rField := reflect.ValueOf(field)

	// Recursively filter on list elements (tags, capabilities, etc.)
	if rField.Kind() == reflect.Slice {
		for i := 0; i < rField.Len(); i++ {
			match, d := f.checkFieldMatchesFilter(rField.Index(i).Interface(), filter)
			diagnostics.Append(d...)
			if diagnostics.HasError() {
				return false, diagnostics
			}

			if match {
				return true, diagnostics
			}
		}

		return false, diagnostics
	}

	normalizedValue, d := f.normalizeValue(field)
	diagnostics.Append(d...)
	if diagnostics.HasError() {
		return false, diagnostics
	}

	for _, value := range filter.Values {
		// We have a match
		if normalizedValue == value.ValueString() {
			return true, diagnostics
		}
	}

	return false, diagnostics
}

func (f FrameworkFilterConfig) normalizeValue(field any) (string, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	rField := reflect.ValueOf(field)

	// Dereference if the value is a pointer
	for rField.Kind() == reflect.Ptr {
		// Null pointer; assume empty
		if rField.IsNil() {
			return "", nil
		}

		rField = reflect.Indirect(rField)
	}

	switch rField.Interface().(type) {
	case string:
		return rField.String(), diagnostics
	case int, int64:
		return strconv.FormatInt(rField.Int(), 10), diagnostics
	case bool:
		return strconv.FormatBool(rField.Bool()), diagnostics
	case float32, float64:
		return strconv.FormatFloat(rField.Float(), 'f', 0, 64), diagnostics
	default:
		diagnostics.AddError(
			"Invalid field type",
			fmt.Sprintf("Invalid type for field: %s", rField.Type().String()),
		)
	}

	return "", diagnostics
}
