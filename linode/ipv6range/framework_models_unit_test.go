package ipv6range

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseIPv6RangeDataSource(t *testing.T) {
	ipRange := linodego.IPv6Range{
		Range:       "2600:3c01::",
		Region:      "us-east",
		Prefix:      64,
		RouteTarget: "2600:3c01::ffff:ffff:ffff:ffff",
		IsBGP:       true,
		Linodes:     []int{123, 456},
	}

	dataSourceModel := DataSourceModel{}
	diags := dataSourceModel.parseIPv6RangeDataSource(context.Background(), &ipRange)

	assert.Nil(t, diags)
	assert.Equal(t, types.StringValue("2600:3c01::"), dataSourceModel.Range)
	assert.Equal(t, dataSourceModel.IsBGP, types.BoolValue(true))
	assert.Equal(t, types.Int64Value(64), dataSourceModel.Prefix)
	assert.Equal(t, types.StringValue("us-east"), dataSourceModel.Region)
	assert.NotEmpty(t, dataSourceModel.ID)
}

func TestParseIPv6RangeResourceDataComputedAttrs(t *testing.T) {
	ipRange := linodego.IPv6Range{
		Range:       "2600:3c01::",
		Region:      "us-east",
		Prefix:      64,
		RouteTarget: "2600:3c01::ffff:ffff:ffff:ffff",
		IsBGP:       false,
		Linodes:     []int{789},
	}

	resourceModel := ResourceModel{}
	diags := resourceModel.parseIPv6RangeResourceDataComputedAttrs(context.Background(), &ipRange)

	assert.Nil(t, diags)
	assert.Equal(t, types.StringValue("2600:3c01::"), resourceModel.Range)
	assert.Equal(t, resourceModel.IsBGP, types.BoolValue(false))
	assert.Equal(t, types.StringValue("us-east"), resourceModel.Region)
}
