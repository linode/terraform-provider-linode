package frameworkfilter

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

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
