//go:build unit

package helper_test

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/stretchr/testify/require"

	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func TestSDKv2PathToCtyPath(t *testing.T) {
	result := helper.SDKv2PathToCtyPath("wow.0.cool.2.field")

	require.Len(t, result, 5)
	require.Equal(t, "wow", result[0].(cty.GetAttrStep).Name)
	require.True(t, result[1].(cty.IndexStep).Key.Equals(cty.NumberIntVal(0)).True())
	require.Equal(t, "cool", result[2].(cty.GetAttrStep).Name)
	require.True(t, result[3].(cty.IndexStep).Key.Equals(cty.NumberIntVal(2)).True())
	require.Equal(t, "field", result[4].(cty.GetAttrStep).Name)
}
