//go:build unit

package networkingip

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

type vpcNAT1To1Model struct {
	Address  string `tfsdk:"address"`
	SubnetID int64  `tfsdk:"subnet_id"`
	VPCID    int64  `tfsdk:"vpc_id"`
}

var testInstanceIP = &linodego.InstanceIP{
	Address:    "123.123.123.123",
	Gateway:    "123.123.123.1",
	SubnetMask: "255.255.255.0",
	Prefix:     32,
	Type:       "ipv4",
	Public:     true,
	RDNS:       "123-123-123-123.ip.linodeusercontent.com",
	LinodeID:   12345,
	Region:     "us-mia",
	Reserved:   false,
	VPCNAT1To1: &linodego.InstanceIPNAT1To1{
		Address:  "10.0.0.1",
		SubnetID: 456,
		VPCID:    123,
	},
}

func TestResourceModelParseInstanceIP(t *testing.T) {
	data := &ResourceModel{}

	var diags diag.Diagnostics
	diags.Append(data.FlattenIPAddress(testInstanceIP, false)...)
	if diags.HasError() {
		t.Fatal(diags.Errors())
	}

	assert.Equal(t, "123.123.123.123", data.ID.ValueString())
	assert.Equal(t, "123.123.123.123", data.Address.ValueString())
	assert.Equal(t, "123.123.123.1", data.Gateway.ValueString())
	assert.Equal(t, int64(32), data.Prefix.ValueInt64())
	assert.Equal(t, "ipv4", data.Type.ValueString())
	assert.Equal(t, true, data.Public.ValueBool())
	assert.Equal(t, "123-123-123-123.ip.linodeusercontent.com", data.RDNS.ValueString())
	assert.Equal(t, int64(12345), data.LinodeID.ValueInt64())
	assert.Equal(t, "us-mia", data.Region.ValueString())
	assert.Equal(t, false, data.Reserved.ValueBool())

	var vpcNAT1To1 vpcNAT1To1Model

	diags.Append(
		data.VPCNAT1To1.As(
			context.Background(),
			&vpcNAT1To1,
			basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			},
		)...,
	)
	if diags.HasError() {
		t.Fatal(diags.Errors())
	}

	assert.Equal(t, "10.0.0.1", vpcNAT1To1.Address)
	assert.Equal(t, int64(123), vpcNAT1To1.VPCID)
	assert.Equal(t, int64(456), vpcNAT1To1.SubnetID)
}

func TestDataSourceModelParseInstanceIP(t *testing.T) {
	data := &DataSourceModel{}

	var diags diag.Diagnostics
	diags.Append(data.parseIP(testInstanceIP)...)
	if diags.HasError() {
		t.Fatal(diags.Errors())
	}

	assert.Equal(t, "123.123.123.123", data.ID.ValueString())
	assert.Equal(t, "123.123.123.123", data.Address.ValueString())
	assert.Equal(t, "123.123.123.1", data.Gateway.ValueString())
	assert.Equal(t, int64(32), data.Prefix.ValueInt64())
	assert.Equal(t, "ipv4", data.Type.ValueString())
	assert.Equal(t, true, data.Public.ValueBool())
	assert.Equal(t, "123-123-123-123.ip.linodeusercontent.com", data.RDNS.ValueString())
	assert.Equal(t, int64(12345), data.LinodeID.ValueInt64())
	assert.Equal(t, "us-mia", data.Region.ValueString())
	assert.Equal(t, false, data.Reserved.ValueBool())

	var vpcNAT1To1 vpcNAT1To1Model

	diags.Append(
		data.VPCNAT1To1.As(
			context.Background(),
			&vpcNAT1To1,
			basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			},
		)...,
	)
	if diags.HasError() {
		t.Fatal(diags.Errors())
	}

	assert.Equal(t, "10.0.0.1", vpcNAT1To1.Address)
	assert.Equal(t, int64(123), vpcNAT1To1.VPCID)
	assert.Equal(t, int64(456), vpcNAT1To1.SubnetID)
}
