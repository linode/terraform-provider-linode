package frameworkfilter

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestConstructFilterString(t *testing.T) {
	testFiltersModel := []FilterModel{
		{
			Name: types.StringValue("api_foo"),
			Values: []types.String{
				types.StringValue("cool"),
				types.StringValue("wow"),
			},
			MatchBy: types.StringValue("exact"),
		},
		{
			Name: types.StringValue("api_bar"),
			Values: []types.String{
				types.StringValue("test"),
			},
		},
		{
			Name: types.StringValue("api_foo"),
			Values: []types.String{
				types.StringValue("wow"),
			},
			MatchBy: types.StringValue("sub"),
		},
		{
			Name: types.StringValue("foo"),
			Values: []types.String{
				types.StringValue("wow"),
			},
		},
		{
			Name: types.StringValue("api_foo_int"),
			Values: []types.String{
				types.StringValue("123"),
			},
		},
		{
			Name: types.StringValue("api_foo_bool"),
			Values: []types.String{
				types.StringValue("true"),
			},
		},
	}

	expectedJSONData := map[string]any{
		"+and": []map[string]any{
			{
				"+or": []map[string]any{
					{
						"api_foo": "cool",
					},
					{
						"api_foo": "wow",
					},
				},
			},
			{
				"+or": []map[string]any{
					{
						"api_bar": "test",
					},
				},
			},
			{
				"+or": []map[string]any{
					{
						"api_foo_int": 123,
					},
				},
			},
			{
				"+or": []map[string]any{
					{
						"api_foo_bool": true,
					},
				},
			},
		},
		"+order":    "api_foo",
		"+order_by": "asc",
	}
	expectedJSONBytes, _ := json.Marshal(expectedJSONData)
	expectedJSON := string(expectedJSONBytes)

	result, d := testFilterConfig.constructFilterString(
		testFiltersModel,
		types.StringValue("api_foo"),
		types.StringValue("asc"),
	)
	if d != nil {
		t.Fatal(d.Detail())
	}

	if !reflect.DeepEqual(expectedJSON, result) {
		t.Fatal(cmp.Diff(expectedJSON, result))
	}
}
