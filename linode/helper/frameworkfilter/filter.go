package frameworkfilter

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"golang.org/x/crypto/sha3"
)

// ListFunc is a wrapper for functions that will list and return values from the API.
type ListFunc func(ctx context.Context, client *linodego.Client, filter string) ([]any, error)

// FilterModel describes the Terraform resource data model to match the
// resource schema.
//
//nolint:all
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
						f.validateFilterable(false),
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
					Description: "The type of comparison to use for this filter.",
					Validators: []validator.String{
						stringvalidator.OneOfCaseInsensitive(
							"exact", "substring", "sub", "re", "regex",
						),
					},
				},
			},
		},
	}
}

// OrderSchema returns the schema for the top-level `order` field.
func (f Config) OrderSchema() schema.StringAttribute {
	return schema.StringAttribute{
		Description: "The order in which results should be returned.",
		Validators: []validator.String{
			stringvalidator.OneOfCaseInsensitive(
				"asc", "desc",
			),
		},
		Optional: true,
	}
}

// OrderBySchema returns the schema for the top-level `order_by` field.
func (f Config) OrderBySchema() schema.StringAttribute {
	return schema.StringAttribute{
		Description: "The attribute to order the results by.",
		Validators: []validator.String{
			f.validateFilterable(true),
		},
		Optional: true,
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
	order types.String,
	orderBy types.String,
) ([]any, diag.Diagnostic) {
	// Construct the API filter string
	filterStr, d := f.constructFilterString(filters, order, orderBy)
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
