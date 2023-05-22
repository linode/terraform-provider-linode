package frameworkfilter

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"golang.org/x/crypto/sha3"
	"reflect"
	"strconv"
)

// ListFunc is a wrapper for functions that will list and return values from the API.
type ListFunc func(ctx context.Context, client *linodego.Client, filter string) ([]any, error)

// FilterModel describes the Terraform resource data model to match the
// resource schema.
type FilterModel struct {
	Name    types.String   `tfsdk:"name" json:"name"`
	Values  []types.String `tfsdk:"values" json:"values"`
	MatchBy types.String   `tfsdk:"match_by" json:"match_by"`
}

// FiltersModelType should be used for the `filter` attribute in list
// data sources.
type FiltersModelType []FilterModel

// FilterAttribute is used to configure filtering for an individual
// response field.
type FilterAttribute struct {
	APIFilterable bool
}

// Config is the root configuration type for filter data sources.
type Config map[string]FilterAttribute

// Schema returns the schema that should be used for the `filter` attribute
// in list data sources.
func (f Config) Schema() schema.SetNestedBlock {
	return schema.SetNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					Required: true,
					Validators: []validator.String{
						validateFilterable(f),
					},
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
}

// GenerateID will generate a unique ID from the given filters.
func (f Config) GenerateID(filters []FilterModel) (types.String, diag.Diagnostic) {
	jsonMap := make([]map[string]any, len(filters))

	// Terraform types cannot be marshalled directly into JSON,
	// so we should convert them into their underlying primitives.
	for i, filter := range filters {
		values := make([]string, len(filter.Values))
		for i, v := range filter.Values {
			values[i] = v.ValueString()
		}

		jsonMap[i] = map[string]any{
			"name":     filter.Name.ValueString(),
			"match_by": filter.MatchBy.ValueString(),
			"values":   values,
		}
	}

	filterJSON, err := json.Marshal(jsonMap)
	if err != nil {
		return types.StringNull(), diag.NewErrorDiagnostic(
			"Failed to marshal JSON.",
			err.Error(),
		)
	}

	hash := sha3.Sum512(filterJSON)
	return types.StringValue(base64.StdEncoding.EncodeToString(hash[:])), nil
}

// GetAndFilter will run all filter operations given the parameters
// and return a list of API response objects.
func (f Config) GetAndFilter(
	ctx context.Context,
	client *linodego.Client,
	filters []FilterModel,
	listFunc ListFunc,
) ([]any, diag.Diagnostic) {
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

	return locallyFilteredElements, nil
}

// constructFilterString constructs a filter string intended to be
// used in ListFunc.
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

// applyLocalFiltering handles filtering for fields that are not
// API-filterable.
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

// matchesFilter checks whether an object matches the given filter set.
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
		matchingField, d := resolveStructValueByJSON(elem, filterName)
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

// checkFieldMatchesFilter checks whether an individual field
// meets the condition for the given filter.
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

// normalizeValue converts the given field into a comparable string.
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
