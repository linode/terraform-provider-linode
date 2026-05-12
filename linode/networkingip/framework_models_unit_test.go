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
	Tags:       []string{"tag1", "tag2"},
	VPCNAT1To1: &linodego.InstanceIPNAT1To1{
		Address:  "10.0.0.1",
		SubnetID: 456,
		VPCID:    123,
	},
	AssignedEntity: &linodego.ReservedIPAssignedEntity{
		ID:    12345,
		Label: "my-linode",
		Type:  "linode",
		URL:   "/v4/linode/instances/12345",
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

	assert.Equal(t, 2, len(data.Tags.Elements()))
	assert.Contains(t, data.Tags.String(), "tag1")
	assert.Contains(t, data.Tags.String(), "tag2")

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

	assert.Equal(t, 2, len(data.Tags.Elements()))
	assert.Contains(t, data.Tags.String(), "tag1")
	assert.Contains(t, data.Tags.String(), "tag2")

	// assigned_entity is populated
	assert.False(t, data.AssignedEntity.IsNull())
	assert.Contains(t, data.AssignedEntity.String(), "my-linode")
	assert.Contains(t, data.AssignedEntity.String(), "/v4/linode/instances/12345")

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

func TestResourceModelParseReservedIPWithTags(t *testing.T) {
	reservedIP := &linodego.InstanceIP{
		Address:    "192.168.1.1",
		Gateway:    "192.168.1.0",
		SubnetMask: "255.255.255.0",
		Prefix:     24,
		Type:       "ipv4",
		Public:     true,
		RDNS:       "192-168-1-1.ip.linodeusercontent.com",
		LinodeID:   0,
		Region:     "us-east",
		Reserved:   true,
		Tags:       []string{"prod", "web"},
	}

	data := &ResourceModel{}
	var diags diag.Diagnostics
	diags.Append(data.FlattenIPAddress(reservedIP, false)...)
	if diags.HasError() {
		t.Fatal(diags.Errors())
	}

	assert.Equal(t, true, data.Reserved.ValueBool())
	assert.True(t, data.LinodeID.IsNull())
	assert.Equal(t, 2, len(data.Tags.Elements()))
	assert.Contains(t, data.Tags.String(), "prod")
	assert.Contains(t, data.Tags.String(), "web")
}

func TestResourceModelParseEmptyTags(t *testing.T) {
	ip := &linodego.InstanceIP{
		Address:    "10.0.0.1",
		Gateway:    "10.0.0.0",
		SubnetMask: "255.255.255.0",
		Prefix:     24,
		Type:       "ipv4",
		Public:     false,
		RDNS:       "",
		LinodeID:   99,
		Region:     "us-west",
		Reserved:   false,
		Tags:       []string{},
	}

	data := &ResourceModel{}
	var diags diag.Diagnostics
	diags.Append(data.FlattenIPAddress(ip, false)...)
	if diags.HasError() {
		t.Fatal(diags.Errors())
	}

	assert.Equal(t, 0, len(data.Tags.Elements()))
}

func TestDataSourceModelParseTags(t *testing.T) {
	ip := &linodego.InstanceIP{
		Address:    "1.2.3.4",
		Gateway:    "1.2.3.1",
		SubnetMask: "255.255.255.0",
		Prefix:     32,
		Type:       "ipv4",
		Public:     true,
		RDNS:       "rdns.example.com",
		LinodeID:   100,
		Region:     "us-mia",
		Reserved:   true,
		Tags:       []string{"alpha", "beta", "gamma"},
	}

	data := &DataSourceModel{}
	var diags diag.Diagnostics
	diags.Append(data.parseIP(ip)...)
	if diags.HasError() {
		t.Fatal(diags.Errors())
	}

	assert.Equal(t, true, data.Reserved.ValueBool())
	assert.Equal(t, 3, len(data.Tags.Elements()))
	assert.Contains(t, data.Tags.String(), "alpha")
	assert.Contains(t, data.Tags.String(), "beta")
	assert.Contains(t, data.Tags.String(), "gamma")
}
