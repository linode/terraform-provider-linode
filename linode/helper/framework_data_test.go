package helper_test

import (
	"context"
	//"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/stretchr/testify/assert"
)

// Testing TF schema for generating object type

var testNestedAttrsSchema = schema.SingleNestedAttribute{
	Attributes: map[string]schema.Attribute{
		"another_number": schema.Int64Attribute{},
	},
}

var testSchema = schema.SingleNestedAttribute{
	Attributes: map[string]schema.Attribute{
		"some_number":              schema.Int64Attribute{},
		"some_string":              schema.StringAttribute{},
		"some_nullable_string":     schema.StringAttribute{},
		"some_bool":                schema.BoolAttribute{},
		"some_nested_attrs_object": testNestedAttrsSchema,
	},
}

// Testing TF models

type DataUnitTestNestedAttrsModel struct {
	AnotherNumber types.Int64 `tfsdk:"another_number"`
}

type DataUnitTestAttrsModel struct {
	SomeNumber            types.Int64  `tfsdk:"some_number"`
	SomeString            types.String `tfsdk:"some_string"`
	SomeNullableString    types.String `tfsdk:"some_nullable_string"`
	SomeBool              types.Bool   `tfsdk:"some_bool"`
	SomeNestedAttrsObject types.Object `tfsdk:"some_nested_attrs_object"`
}

// Testing structs, mocking usual linodego structs

type DataUnitTestExpandedNestedObject struct {
	AnotherNumber int
}

type DataUnitTestExpandedObject struct {
	SomeNumber         int
	SomeString         string
	SomeNullableString *string
	SomeBool           bool
	SomeNestedObject   DataUnitTestExpandedNestedObject
}

func TestKeepOrUpdateSingleNestedAttribute(t *testing.T) {
	// This test covers both the general behavior of KeepOrUpdateSingleNestedAttributes
	// and the specific isNull boolean parameter behavior. The isNull parameter allows
	// the flatten function to indicate that the target should be a Terraform null value.
	ctx := context.Background()
	objectType := testSchema.GetType().(types.ObjectType).AttrTypes
	nestedObjectType := testNestedAttrsSchema.GetType().(types.ObjectType).AttrTypes

	expanded := DataUnitTestExpandedObject{
		SomeNumber:         123,
		SomeString:         "Hello, world!",
		SomeNullableString: nil,
		SomeBool:           true,
		SomeNestedObject: DataUnitTestExpandedNestedObject{
			AnotherNumber: 1234,
		},
	}

	expectedWhenOverrideAll := types.ObjectValueMust(
		objectType,
		map[string]attr.Value{
			"some_number":          types.Int64Value(int64(expanded.SomeNumber)),
			"some_string":          types.StringValue(expanded.SomeString),
			"some_nullable_string": types.StringPointerValue(expanded.SomeNullableString),
			"some_bool":            types.BoolValue(expanded.SomeBool),
			"some_nested_attrs_object": types.ObjectValueMust(
				nestedObjectType,
				map[string]attr.Value{
					"another_number": types.Int64Value(int64(expanded.SomeNestedObject.AnotherNumber)),
				},
			),
		},
	)

	tests := map[string]struct {
		input         types.Object
		data          DataUnitTestExpandedObject
		expected      types.Object
		preserveKnown bool
		setIsNull     bool
	}{
		"unknown object with preserving known": {
			input:         types.ObjectUnknown(objectType),
			data:          expanded,
			expected:      expectedWhenOverrideAll,
			preserveKnown: true,
			setIsNull:     false,
		},
		"null object with preserving known": {
			input:         types.ObjectNull(objectType),
			data:          expanded,
			expected:      types.ObjectNull(objectType),
			preserveKnown: true,
			setIsNull:     false,
		},
		"partially known object with preserving known": {
			input: types.ObjectValueMust(
				objectType,
				map[string]attr.Value{
					"some_number":          types.Int64Value(12345),
					"some_string":          types.StringValue("Hey there!"),
					"some_nullable_string": types.StringUnknown(),
					"some_bool":            types.BoolUnknown(),
					"some_nested_attrs_object": types.ObjectValueMust(
						nestedObjectType,
						map[string]attr.Value{
							"another_number": types.Int64Unknown(),
						},
					),
				},
			),
			data: expanded,
			expected: types.ObjectValueMust(
				objectType,
				map[string]attr.Value{
					"some_number":          types.Int64Value(12345),
					"some_string":          types.StringValue("Hey there!"),
					"some_nullable_string": types.StringPointerValue(expanded.SomeNullableString),
					"some_bool":            types.BoolValue(expanded.SomeBool),
					"some_nested_attrs_object": types.ObjectValueMust(
						nestedObjectType,
						map[string]attr.Value{
							"another_number": types.Int64Value(int64(expanded.SomeNestedObject.AnotherNumber)),
						},
					),
				},
			),
			preserveKnown: true,
			setIsNull:     false,
		},
		"null nested object with preserving known": {
			input: types.ObjectValueMust(
				objectType,
				map[string]attr.Value{
					"some_number":              types.Int64Value(12345),
					"some_string":              types.StringUnknown(),
					"some_nullable_string":     types.StringUnknown(),
					"some_bool":                types.BoolValue(false),
					"some_nested_attrs_object": types.ObjectNull(nestedObjectType),
				},
			),
			data: expanded,
			expected: types.ObjectValueMust(
				objectType,
				map[string]attr.Value{
					"some_number":              types.Int64Value(12345),
					"some_string":              types.StringValue(expanded.SomeString),
					"some_nullable_string":     types.StringPointerValue(expanded.SomeNullableString),
					"some_bool":                types.BoolValue(false),
					"some_nested_attrs_object": types.ObjectNull(nestedObjectType),
				},
			),
			preserveKnown: true,
			setIsNull:     false,
		},
		"unknown nested object with preserving known": {
			input: types.ObjectValueMust(
				objectType,
				map[string]attr.Value{
					"some_number":              types.Int64Value(12345),
					"some_string":              types.StringUnknown(),
					"some_nullable_string":     types.StringUnknown(),
					"some_bool":                types.BoolValue(false),
					"some_nested_attrs_object": types.ObjectUnknown(nestedObjectType),
				},
			),
			data: expanded,
			expected: types.ObjectValueMust(
				objectType,
				map[string]attr.Value{
					"some_number":          types.Int64Value(12345),
					"some_string":          types.StringValue(expanded.SomeString),
					"some_nullable_string": types.StringPointerValue(expanded.SomeNullableString),
					"some_bool":            types.BoolValue(false),
					"some_nested_attrs_object": types.ObjectValueMust(
						nestedObjectType,
						map[string]attr.Value{
							"another_number": types.Int64Value(int64(expanded.SomeNestedObject.AnotherNumber)),
						},
					),
				},
			),
			preserveKnown: true,
			setIsNull:     false,
		},
		"all known with preserving known": {
			input: types.ObjectValueMust(
				objectType,
				map[string]attr.Value{
					"some_number":          types.Int64Value(12345),
					"some_string":          types.StringValue("Hey there!"),
					"some_nullable_string": types.StringValue("I'm here!"),
					"some_bool":            types.BoolValue(true),
					"some_nested_attrs_object": types.ObjectValueMust(
						nestedObjectType,
						map[string]attr.Value{
							"another_number": types.Int64Value(321),
						},
					),
				},
			),
			data: expanded,
			expected: types.ObjectValueMust(
				objectType,
				map[string]attr.Value{
					"some_number":          types.Int64Value(12345),
					"some_string":          types.StringValue("Hey there!"),
					"some_nullable_string": types.StringValue("I'm here!"),
					"some_bool":            types.BoolValue(true),
					"some_nested_attrs_object": types.ObjectValueMust(
						nestedObjectType,
						map[string]attr.Value{
							"another_number": types.Int64Value(321),
						},
					),
				},
			),
			preserveKnown: true,
			setIsNull:     false,
		},

		"unknown object without preserving known": {
			input:         types.ObjectUnknown(objectType),
			data:          expanded,
			expected:      expectedWhenOverrideAll,
			preserveKnown: false,
			setIsNull:     false,
		},
		"null object without preserving known": {
			input:         types.ObjectNull(objectType),
			data:          expanded,
			expected:      expectedWhenOverrideAll,
			preserveKnown: false,
			setIsNull:     false,
		},
		"partially known object without preserving known": {
			input: types.ObjectValueMust(
				objectType,
				map[string]attr.Value{
					"some_number":          types.Int64Value(12345),
					"some_string":          types.StringUnknown(),
					"some_nullable_string": types.StringUnknown(),
					"some_bool":            types.BoolValue(false),
					"some_nested_attrs_object": types.ObjectValueMust(
						nestedObjectType,
						map[string]attr.Value{
							"another_number": types.Int64Unknown(),
						},
					),
				},
			),
			data:          expanded,
			expected:      expectedWhenOverrideAll,
			preserveKnown: false,
			setIsNull:     false,
		},
		"null nested object without preserving known": {
			input: types.ObjectValueMust(
				objectType,
				map[string]attr.Value{
					"some_number":              types.Int64Value(12345),
					"some_string":              types.StringUnknown(),
					"some_nullable_string":     types.StringUnknown(),
					"some_bool":                types.BoolValue(false),
					"some_nested_attrs_object": types.ObjectNull(nestedObjectType),
				},
			),
			data:          expanded,
			expected:      expectedWhenOverrideAll,
			preserveKnown: false,
			setIsNull:     false,
		},
		"unknown nested object without preserving known": {
			input: types.ObjectValueMust(
				objectType,
				map[string]attr.Value{
					"some_number":              types.Int64Value(12345),
					"some_string":              types.StringUnknown(),
					"some_nullable_string":     types.StringUnknown(),
					"some_bool":                types.BoolValue(false),
					"some_nested_attrs_object": types.ObjectUnknown(nestedObjectType),
				},
			),
			data:          expanded,
			expected:      expectedWhenOverrideAll,
			preserveKnown: false,
			setIsNull:     false,
		},
		"all known without preserving known": {
			input: types.ObjectValueMust(
				objectType,
				map[string]attr.Value{
					"some_number":          types.Int64Value(12345),
					"some_string":          types.StringValue("Hey there!"),
					"some_nullable_string": types.StringNull(),
					"some_bool":            types.BoolValue(false),
					"some_nested_attrs_object": types.ObjectValueMust(
						nestedObjectType,
						map[string]attr.Value{
							"another_number": types.Int64Null(),
						},
					),
				},
			),
			data:          expanded,
			expected:      expectedWhenOverrideAll,
			preserveKnown: false,
			setIsNull:     false,
		},

		// isNull behavior tests
		"isNull=true without preserving known should return null object": {
			input: types.ObjectValueMust(
				objectType,
				map[string]attr.Value{
					"some_number":          types.Int64Value(999),
					"some_string":          types.StringValue("Existing value"),
					"some_nullable_string": types.StringValue("Keep me"),
					"some_bool":            types.BoolValue(false),
					"some_nested_attrs_object": types.ObjectValueMust(
						nestedObjectType,
						map[string]attr.Value{
							"another_number": types.Int64Value(888),
						},
					),
				},
			),
			data:          expanded,
			expected:      types.ObjectNull(objectType),
			preserveKnown: false,
			setIsNull:     true,
		},
		"isNull=true with preserving known should ignore isNull and preserve original": {
			input: types.ObjectValueMust(
				objectType,
				map[string]attr.Value{
					"some_number":          types.Int64Value(999),
					"some_string":          types.StringValue("Existing value"),
					"some_nullable_string": types.StringValue("Keep me"),
					"some_bool":            types.BoolValue(false),
					"some_nested_attrs_object": types.ObjectValueMust(
						nestedObjectType,
						map[string]attr.Value{
							"another_number": types.Int64Value(888),
						},
					),
				},
			),
			data: expanded,
			expected: types.ObjectValueMust(
				objectType,
				map[string]attr.Value{
					"some_number":          types.Int64Value(999),
					"some_string":          types.StringValue("Existing value"),
					"some_nullable_string": types.StringValue("Keep me"),
					"some_bool":            types.BoolValue(false),
					"some_nested_attrs_object": types.ObjectValueMust(
						nestedObjectType,
						map[string]attr.Value{
							"another_number": types.Int64Value(888),
						},
					),
				},
			),
			preserveKnown: true,
			setIsNull:     true,
		},
		"isNull=true with preserving known on null input should return null": {
			input:         types.ObjectNull(objectType),
			data:          expanded,
			expected:      types.ObjectNull(objectType),
			preserveKnown: true,
			setIsNull:     true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var diags diag.Diagnostics

			flatten := func(model *DataUnitTestAttrsModel, isNull *bool, preserveKnown bool, diags *diag.Diagnostics) {
				*isNull = tt.setIsNull

				if !*isNull {
					// Only populate the model if we're not setting it to null
					model.SomeNumber = helper.KeepOrUpdateInt64(model.SomeNumber, int64(tt.data.SomeNumber), preserveKnown)
					model.SomeString = helper.KeepOrUpdateString(model.SomeString, tt.data.SomeString, preserveKnown)
					model.SomeNullableString = helper.KeepOrUpdateStringPointer(model.SomeNullableString, tt.data.SomeNullableString, preserveKnown)
					model.SomeBool = helper.KeepOrUpdateBool(model.SomeBool, tt.data.SomeBool, preserveKnown)

					flattenNestedObject := helper.KeepOrUpdateSingleNestedAttributesWithTypes(
						ctx,
						model.SomeNestedAttrsObject,
						nestedObjectType,
						preserveKnown,
						diags,
						func(nestedModel *DataUnitTestNestedAttrsModel, _ *bool, preserveKnown bool, _ *diag.Diagnostics) {
							nestedModel.AnotherNumber = helper.KeepOrUpdateInt64(
								nestedModel.AnotherNumber,
								int64(tt.data.SomeNestedObject.AnotherNumber),
								preserveKnown,
							)
						},
					)

					if diags.HasError() {
						return
					}

					model.SomeNestedAttrsObject = *flattenNestedObject
				}
			}

			out := helper.KeepOrUpdateSingleNestedAttributes(ctx, tt.input, tt.preserveKnown, &diags, flatten)
			if diags.HasError() {
				t.Fatalf("unexpected error: %v", diags)
			}

			assert.Truef(t, out.Equal(tt.expected),
				"Flattened object (%v) should match expected object (%v)", out, tt.expected)
		})
	}
}
