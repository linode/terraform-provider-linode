package helper

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/linode/linodego"
	"golang.org/x/crypto/sha3"
)

// FilterConfig stores a map of FilterAttributes for a resource.
type FilterConfig map[string]FilterAttribute

// FilterTypeFunc is a function that takes in a filter name and value,
// and returns the value converted to the correct filter type.
type FilterTypeFunc func(value string) (interface{}, error)

// FilterListFunc wraps a linodego list function.
type FilterListFunc func(context.Context, *schema.ResourceData,
	*linodego.Client, *linodego.ListOptions) ([]interface{}, error)

// FilterFlattenFunc flattens an object into a map[string]interface{}.
type FilterFlattenFunc func(object interface{}) map[string]interface{}

// FilterAttribute stores various configuration options about a single
// filterable field.
type FilterAttribute struct {
	// Whether this field can be filtered on at an API level.
	// If false, this filter will be handled on the client.
	APIFilterable bool

	// Converts the filter string to the correct type.
	TypeFunc FilterTypeFunc
}

// FilterSchema should be referenced in a schema configuration in order to
// enable filter functionality
func (f FilterConfig) FilterSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:             schema.TypeString,
					Description:      "The name of the attribute to filter on.",
					ValidateDiagFunc: f.ValidateDiagFunc(false),
					Required:         true,
				},
				"values": {
					Type:        schema.TypeList,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: "The value(s) to be used in the filter.",
					Required:    true,
				},
				"match_by": {
					Type:        schema.TypeString,
					Description: "The type of comparison to use for this filter.",
					Optional:    true,
					Default:     "exact",
					ValidateFunc: validation.StringInSlice([]string{"exact", "substring", "sub", "re", "regex"},
						false),
				},
			},
		},
	}
}

// OrderBySchema should be referenced in a schema configuration in order to
// enable filter ordering functionality
func (f FilterConfig) OrderBySchema() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: f.ValidateDiagFunc(true),
		Description:      "The attribute to order the results by.",
	}
}

func (f FilterConfig) OrderSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "asc",
		ValidateFunc: validation.StringInSlice([]string{"asc", "desc"}, false),
		Description:  "The order in which results should be returned.",
	}
}

// ConstructFilterString constructs a Linode filter JSON string from each filter element in the schema
func (f FilterConfig) ConstructFilterString(d *schema.ResourceData) (string, error) {
	filters := d.Get("filter").([]interface{})
	resultMap := make(map[string]interface{})

	if len(filters) < 1 {
		return "{}", nil
	}

	var rootFilter []interface{}

	for _, filter := range filters {
		filter := filter.(map[string]interface{})

		name := filter["name"].(string)
		values := filter["values"].([]interface{})
		matchBy := filter["match_by"].(string)

		// Defer this logic to the client
		if matchBy != "exact" {
			continue
		}

		// Defer this logic to the client if not API-filterable
		if cfg, ok := f[name]; !ok || !cfg.APIFilterable {
			continue
		}

		subFilter := make([]interface{}, len(values))

		for i, value := range values {
			value, err := f[name].TypeFunc(value.(string))
			if err != nil {
				return "", err
			}

			valueFilter := make(map[string]interface{})
			valueFilter[name] = value

			subFilter[i] = valueFilter
		}

		rootFilter = append(rootFilter, map[string]interface{}{
			"+or": subFilter,
		})
	}

	if len(rootFilter) < 1 {
		return "{}", nil
	}

	resultMap["+and"] = rootFilter

	if orderBy, ok := d.GetOk("order_by"); ok {
		resultMap["+order_by"] = orderBy
		resultMap["+order"] = d.Get("order")
	}

	result, err := json.Marshal(resultMap)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// FilterResults filters the given results on the client-side filters present in the resource.
func (f FilterConfig) FilterResults(
	d *schema.ResourceData,
	items []interface{},
) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)

	for _, item := range items {
		item := item.(map[string]interface{})

		match, err := f.itemMatchesFilter(d, item)
		if err != nil {
			return nil, err
		}

		if !match {
			continue
		}

		result = append(result, item)
	}

	return result, nil
}

// FilterDataSource should be run from inside the ReadContext function of a data source.
func (f FilterConfig) FilterDataSource(
	ctx context.Context,
	d *schema.ResourceData,
	meta interface{},
	listFunc FilterListFunc,
	flattenFunc FilterFlattenFunc,
) ([]map[string]interface{}, error) {
	client := meta.(*ProviderMeta).Client

	filterID, err := f.GetFilterID(d)
	if err != nil {
		return nil, fmt.Errorf("failed to generate filter id: %s", err)
	}

	filter, err := f.ConstructFilterString(d)
	if err != nil {
		return nil, fmt.Errorf("failed to construct filter: %s", err)
	}

	// Call linode list function defined by data source
	items, err := listFunc(ctx, d, &client, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list linode items: %s", err)
	}

	itemsFlattened := make([]interface{}, len(items))
	for i, image := range items {
		itemsFlattened[i] = flattenFunc(image)
	}

	itemsFiltered, err := f.FilterResults(d, itemsFlattened)
	if err != nil {
		return nil, fmt.Errorf("failed to filter returned data: %s", err)
	}

	d.SetId(filterID)

	return itemsFiltered, nil
}

// GetValidFilters returns a slice of valid filters for the filter config.
func (f FilterConfig) GetValidFilters(apiOnly bool) []string {
	result := make([]string, 0)

	for k, v := range f {
		if apiOnly && !v.APIFilterable {
			continue
		}

		result = append(result, k)
	}

	return result
}

// ValidateDiagFunc should be plugged into the `filter` field of a filterable data source.
func (f FilterConfig) ValidateDiagFunc(apiOnly bool) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		val := i.(string)

		cfg, ok := f[val]

		if !ok {
			return diag.Errorf("\"%s\" is not a filterable field. Valid filters: %s",
				val, strings.Join(f.GetValidFilters(false), ", "))
		}

		if apiOnly && !cfg.APIFilterable {
			return diag.Errorf("\"%s\" is an unsupported filter for this field. Valid filters: %s",
				val, strings.Join(f.GetValidFilters(true), ", "))
		}

		return nil
	}
}

// GetFilterID creates a unique ID specific to the current filter data source
func (f FilterConfig) GetFilterID(d *schema.ResourceData) (string, error) {
	idMap := map[string]interface{}{
		"filter":   d.Get("filter"),
		"order":    d.Get("order"),
		"order_by": d.Get("order_by"),
	}

	result, err := json.Marshal(idMap)
	if err != nil {
		return "", err
	}

	hash := sha3.Sum512(result)
	return base64.StdEncoding.EncodeToString(hash[:]), nil
}

// FilterLatest returns only the latest element in the given slice only if `latest` == true.
func (f FilterConfig) FilterLatest(d *schema.ResourceData, items []map[string]interface{}) []map[string]interface{} {
	if !d.Get("latest").(bool) {
		return items
	}

	if item := f.GetLatestCreated(items); item != nil {
		return []map[string]interface{}{item}
	}

	return []map[string]interface{}{}
}

// GetLatestCreated returns only the latest element in the given slice.
func (f FilterConfig) GetLatestCreated(data []map[string]interface{}) map[string]interface{} {
	var latestCreated time.Time
	var latestEntity map[string]interface{}

	for _, image := range data {
		created, ok := image["created"]
		if !ok {
			continue
		}

		createdTime, err := time.Parse(time.RFC3339, created.(string))
		if err != nil {
			return nil
		}

		if latestEntity != nil && !createdTime.After(latestCreated) {
			continue
		}

		latestCreated = createdTime
		latestEntity = image
	}

	return latestEntity
}

// FilterLatestVersion returns only the latest element in the given slice only if `latest` == true.
func (f FilterConfig) FilterLatestVersion(d *schema.ResourceData,
	items []map[string]interface{},
) ([]map[string]interface{}, error) {
	if !d.Get("latest").(bool) {
		return items, nil
	}

	item, err := f.GetLatestVersion(items)
	if err != nil {
		return nil, err
	}

	if item != nil {
		return []map[string]interface{}{item}, nil
	}

	return []map[string]interface{}{}, nil
}

// GetLatestVersion returns only the latest version string (e.g. `8.0.26`) in the given slice.
func (f FilterConfig) GetLatestVersion(data []map[string]interface{}) (map[string]interface{}, error) {
	var latestVersion []int

	var latestEntity map[string]interface{}

	for _, entity := range data {
		version, ok := entity["version"]
		if !ok {
			continue
		}

		versionSplit := strings.Split(version.(string), ".")
		versionSplitInt := make([]int, len(versionSplit))

		// Parse the versions to []int
		for i, v := range versionSplit {
			versionInt, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("failed to parse version segment: %s", v)
			}

			versionSplitInt[i] = versionInt
		}

		// We should allocate the number of version segments dynamically
		if latestVersion == nil {
			latestVersion = make([]int, len(versionSplit))
			// Just in case :)
			for i := range latestVersion {
				latestVersion[i] = 0
			}
		}

		for i, seg := range versionSplitInt {
			if seg < latestVersion[i] {
				break
			}

			if seg > latestVersion[i] {
				latestEntity = entity
				latestVersion = versionSplitInt
				break
			}
		}
	}

	return latestEntity, nil
}

func (f FilterConfig) itemMatchesFilter(
	d *schema.ResourceData,
	item map[string]interface{},
) (bool, error) {
	filters := d.Get("filter").([]interface{})

	for _, filter := range filters {
		filter := filter.(map[string]interface{})

		name := filter["name"].(string)
		values := filter["values"].([]interface{})
		matchBy := filter["match_by"].(string)

		itemValue, ok := item[name]
		if !ok {
			return false, fmt.Errorf("\"%v\" is not a valid attribute", name)
		}

		valid, err := f.validateFilter(matchBy, name, ExpandStringList(values), itemValue)
		if err != nil {
			return false, err
		}

		if !valid {
			return false, nil
		}
	}

	return true, nil
}

func (f FilterConfig) validateFilter(
	matchBy, name string,
	values []string,
	itemValue interface{},
) (bool, error) {
	// Filter recursively on lists (tags, etc.)
	if items, ok := itemValue.([]string); ok {
		for _, item := range items {
			valid, err := f.validateFilter(matchBy, name, values, item)
			if err != nil {
				return false, err
			}

			if valid {
				return true, nil
			}
		}

		return false, nil
	}

	cfg := f[name]

	valuesNormalized := make([]interface{}, len(values))
	for i := range valuesNormalized {
		n, err := cfg.TypeFunc(values[i])
		if err != nil {
			return false, err
		}

		valuesNormalized[i] = n
	}

	switch matchBy {
	case "exact":
		return validateFilterExact(valuesNormalized, itemValue)
	case "substring", "sub":
		return validateFilterSubstring(name, valuesNormalized, itemValue)
	case "re", "regex":
		return validateFilterRegex(name, valuesNormalized, itemValue)
	}

	return true, nil
}

func validateFilterExact(values []interface{}, result interface{}) (bool, error) {
	for _, value := range values {
		if reflect.DeepEqual(result, value) {
			return true, nil
		}
	}

	return false, nil
}

func validateFilterSubstring(name string, values []interface{}, result interface{}) (bool, error) {
	itemValueStr, ok := result.(string)
	if !ok {
		return false, fmt.Errorf("\"%s\" is not a string (type %s) and cannot be filtered on substring",
			name, reflect.TypeOf(result))
	}

	for _, value := range values {
		if strings.Contains(itemValueStr, value.(string)) {
			return true, nil
		}
	}

	return false, nil
}

func validateFilterRegex(name string, values []interface{}, result interface{}) (bool, error) {
	itemValueStr, ok := result.(string)
	if !ok {
		return false, fmt.Errorf("\"%s\" is not a string (type %s) and cannot be filtered on regex",
			name, reflect.TypeOf(result))
	}

	for _, value := range values {
		r, err := regexp.Compile(value.(string))
		if err != nil {
			return false, fmt.Errorf("failed to compile regex: %s", err)
		}

		if r.MatchString(itemValueStr) {
			return true, nil
		}
	}

	return false, nil
}

func FilterTypeString(value string) (interface{}, error) {
	return value, nil
}

func FilterTypeInt(value string) (interface{}, error) {
	return strconv.Atoi(value)
}

func FilterTypeBool(value string) (interface{}, error) {
	return strconv.ParseBool(value)
}
