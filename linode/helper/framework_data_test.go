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

func TestFlattenFirewallSettings(t *testing.T) {
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

	flatten := func(model *DataUnitTestAttrsModel, preserveKnown bool, diags *diag.Diagnostics) {
		model.SomeNumber = helper.KeepOrUpdateInt64(model.SomeNumber, int64(expanded.SomeNumber), preserveKnown)
		model.SomeString = helper.KeepOrUpdateString(model.SomeString, expanded.SomeString, preserveKnown)
		model.SomeNullableString = helper.KeepOrUpdateStringPointer(model.SomeNullableString, expanded.SomeNullableString, preserveKnown)
		model.SomeBool = helper.KeepOrUpdateBool(model.SomeBool, expanded.SomeBool, preserveKnown)

		flattenNestedObject := helper.KeepOrUpdateNestedObjectWithTypes(
			ctx,
			model.SomeNestedAttrsObject,
			nestedObjectType,
			preserveKnown,
			diags,
			func(nestedModel *DataUnitTestNestedAttrsModel, preserveKnown bool, _ *diag.Diagnostics) {
				nestedModel.AnotherNumber = helper.KeepOrUpdateInt64(nestedModel.AnotherNumber, int64(expanded.SomeNestedObject.AnotherNumber), preserveKnown)
			},
		)

		if diags.HasError() {
			return
		}

		model.SomeNestedAttrsObject = *flattenNestedObject
	}

	tests := map[string]struct {
		input         types.Object
		data          DataUnitTestExpandedObject
		expected      types.Object
		preserveKnown bool
	}{
		"unknown default firewall IDs with preserving known": {
			input:         types.ObjectUnknown(objectType),
			data:          expanded,
			expected:      expectedWhenOverrideAll,
			preserveKnown: true,
		},
		"null default firewall IDs with preserving known": {
			input:         types.ObjectNull(objectType),
			data:          expanded,
			expected:      types.ObjectNull(objectType),
			preserveKnown: true,
		},
		"partially known default firewall IDs with preserving known": {
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
		},

		"unknown default firewall IDs without preserving known": {
			input:         types.ObjectUnknown(objectType),
			data:          expanded,
			expected:      expectedWhenOverrideAll,
			preserveKnown: false,
		},
		"null default firewall IDs without preserving known": {
			input:         types.ObjectNull(objectType),
			data:          expanded,
			expected:      expectedWhenOverrideAll,
			preserveKnown: false,
		},
		"partially known default firewall IDs without preserving known": {
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
		},
	}

	for name, tt := range tests {
		var diags diag.Diagnostics
		t.Run(name, func(t *testing.T) {
			out := helper.KeepOrUpdateNestedObject(ctx, tt.input, tt.preserveKnown, &diags, flatten)
			if diags.HasError() {
				t.Fatalf("unexpected error: %v", diags)
			}

			assert.Truef(t, out.Equal(tt.expected),
				"Flattened object (%v) should match expected object (%v)", out, tt.expected)
		})
	}
}
