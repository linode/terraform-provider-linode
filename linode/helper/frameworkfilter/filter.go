package frameworkfilter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"reflect"
	"strconv"
)

type ListFunc func(ctx context.Context, client linodego.Client, filter string) ([]any, error)
type FlattenFunc func(value any) types.Object

var Schema = schema.SetNestedBlock{
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

// FilterModel describes the Terraform resource data model to match the
// resource schema.
type FilterModel struct {
	Name    types.String   `tfsdk:"name"`
	Values  []types.String `tfsdk:"values"`
	MatchBy types.String   `tfsdk:"match_by"`
}

type FiltersModelType []FilterModel

type FilterAttribute struct {
	APIFilterable bool
	Type          attr.Type
}

type Config map[string]FilterAttribute

func (f Config) DataSourceRead(
	ctx context.Context,
	client linodego.Client,
	filters []FilterModel,
	listFunc ListFunc,
	flattenFunc FlattenFunc,
) ([]types.Object, diag.Diagnostic) {
	// Construct the API filter string
	filterStr, d := f.constructFilterString(filters)
	if d != nil {
		return nil, d
	}

	// Call the user-defined list function
	listedElems, err := listFunc(ctx, client, filterStr)
	if err != nil {
		return nil, diag.NewErrorDiagnostic(
			"Failed to list resources",
			err.Error(),
		)
	}

	// Apply local filtering
	locallyFilteredElements, d := f.applyLocalFiltering(filters, listedElems)
	if d != nil {
		return nil, d
	}

	result := make([]types.Object, len(locallyFilteredElements))
	for i, elem := range locallyFilteredElements {
		result[i] = flattenFunc(elem)
	}

	return result, nil
}

func (f Config) ParseFilterSet(
	ctx context.Context,
	filterSet types.Set,
) ([]FilterModel, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	var result []FilterModel

	// Parse out the set into an object list
	diagnostics.Append(
		filterSet.ElementsAs(ctx, &result, false)...,
	)
	if diagnostics.HasError() {
		return nil, diagnostics
	}

	return result, diagnostics
}

func (f Config) constructFilterString(
	filterSet []FilterModel,
) (string, diag.Diagnostic) {
	rootFilter := make([]map[string]any, 0)

	for _, filter := range filterSet {
		// Get string attributes
		filterFieldName := filter.Name.ValueString()

		// Is this field filterable?
		filterFieldConfig, ok := f[filterFieldName]
		if !ok {
			return "", diag.NewErrorDiagnostic(
				"Attempted to filter on non-filterable field.",
				fmt.Sprintf("Attempted to filter on non-filterable field %s.", filterFieldName),
			)
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
		return "", diag.NewErrorDiagnostic(
			"Failed to marshal api filter",
			err.Error(),
		)
	}

	return string(result), nil
}

func (f Config) applyLocalFiltering(
	filterSet []FilterModel, data []any,
) ([]any, diag.Diagnostic) {
	result := make([]any, 0)

	for _, elem := range data {
		match, d := f.matchesFilter(filterSet, elem)
		if d != nil {
			return nil, d
		}

		// This element was filtered out
		if !match {
			continue
		}

		result = append(result, elem)
	}

	return result, nil
}

func (f Config) matchesFilter(
	filterSet []FilterModel,
	elem any,
) (bool, diag.Diagnostic) {
	for _, filter := range filterSet {
		filterName := filter.Name.ValueString()

		// Skip if this field should be filtered at an API level
		if f[filterName].APIFilterable {
			continue
		}

		// Grab the field from the input struct
		matchingField, d := f.resolveStructFieldByJSON(elem, filterName)
		if d != nil {
			return false, d
		}

		// Check whether the field matches the filter
		match, d := f.checkFieldMatchesFilter(matchingField, filter)
		if d != nil {
			return false, d
		}

		// No match for this filter; return
		if !match {
			return false, nil
		}
	}

	return true, nil
}

func (f Config) checkFieldMatchesFilter(
	field any,
	filter FilterModel,
) (bool, diag.Diagnostic) {
	rField := reflect.ValueOf(field)

	// Recursively filter on list elements (tags, capabilities, etc.)
	if rField.Kind() == reflect.Slice {
		for i := 0; i < rField.Len(); i++ {
			match, d := f.checkFieldMatchesFilter(rField.Index(i).Interface(), filter)
			if d != nil {
				return false, d
			}

			if match {
				return true, nil
			}
		}

		return false, nil
	}

	normalizedValue, d := f.normalizeValue(field)
	if d != nil {
		return false, d
	}

	for _, value := range filter.Values {
		// We have a match
		// TODO: support other types of equality checks
		if normalizedValue == value.ValueString() {
			return true, nil
		}
	}

	return false, nil
}

func (f Config) normalizeValue(field any) (string, diag.Diagnostic) {
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
		return rField.String(), nil
	case int, int64:
		return strconv.FormatInt(rField.Int(), 10), nil
	case bool:
		return strconv.FormatBool(rField.Bool()), nil
	case float32, float64:
		return strconv.FormatFloat(rField.Float(), 'f', 0, 64), nil
	default:
		return "", diag.NewErrorDiagnostic(
			"Invalid field type",
			fmt.Sprintf("Invalid type for field: %s", rField.Type().String()),
		)
	}
}

func (f Config) resolveStructFieldByJSON(val any, field string) (any, diag.Diagnostic) {
	rType := reflect.TypeOf(val)

	var targetField reflect.Value

	for i := 0; i < rType.NumField(); i++ {
		currentField := rType.Field(i)
		if tag, ok := currentField.Tag.Lookup("json"); ok && tag == field {
			targetField = reflect.ValueOf(targetField).Field(i)
		}
	}

	if !targetField.IsValid() {
		return nil, diag.NewErrorDiagnostic(
			"Field not found",
			fmt.Sprintf("Could not find JSON tag in target struct: %s", field),
		)
	}

	return targetField.Interface(), nil
}
