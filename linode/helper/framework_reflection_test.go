//go:build unit

package helper_test

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStructToFrameworkObjectType(t *testing.T) {
	type testModel struct {
		Field1 types.Int64 `tfsdk:"field1"`
		Field2 string
		Field3 types.String `tfsdk:"field3"`
		Field4 types.Set    `tfsdk:"field4"`
	}

	result, err := helper.FrameworkModelToObjectType[testModel](t.Context())
	require.NoError(t, err)

	assert.NotContains(t, result.AttrTypes, "field2")

	assert.True(t, reflect.TypeOf(result.AttrTypes["field1"]).AssignableTo(reflect.TypeFor[basetypes.Int64Type]()))
	assert.True(t, reflect.TypeOf(result.AttrTypes["field3"]).AssignableTo(reflect.TypeFor[basetypes.StringType]()))
	assert.True(t, reflect.TypeOf(result.AttrTypes["field4"]).AssignableTo(reflect.TypeFor[basetypes.SetType]()))
}
