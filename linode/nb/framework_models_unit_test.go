//go:build unit

package nb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestFlattenNodeBalancerPreserveKnown(t *testing.T) {
	label := "test-nodebalancer"

	nodeBalancer := &linodego.NodeBalancer{
		ID:    123,
		Label: &label,
	}

	nodeBalancerModel := &NodeBalancerModel{
		ID:    types.StringUnknown(),
		Label: types.StringValue("another" + label),
	}

	diags := nodeBalancerModel.Flatten(
		context.Background(),
		nodeBalancer,
		nil,
		nil,
		true,
	)

	assert.False(t, diags.HasError(), "Errors should be returned due to custom context error")
	assert.False(t, types.StringValue(label).Equal(nodeBalancerModel.Label))
	assert.True(t, types.StringValue("123").Equal(nodeBalancerModel.ID))
}

func TestFlattenNodeBalancer(t *testing.T) {
	hostname := "example.nodebalancer.linode.com"
	IPv4 := "192.168.1.1"
	IPv6 := "2001:db8::1"

	createdTime := time.Date(2023, time.August, 17, 12, 0, 0, 0, time.UTC)
	updatedTime := time.Date(2023, time.August, 17, 14, 0, 0, 0, time.UTC)

	transferIn := float64(100.0)
	transferOut := float64(200.0)
	transferTotal := float64(300.0)

	label := "test-nodebalancer"

	nodeBalancer := &linodego.NodeBalancer{
		ID:                    123,
		Label:                 &label,
		Region:                "us-east",
		ClientConnThrottle:    10,
		ClientUDPSessThrottle: 5,
		Hostname:              &hostname,
		IPv4:                  &IPv4,
		IPv6:                  &IPv6,
		Created:               &createdTime,
		Updated:               &updatedTime,
		Transfer: linodego.NodeBalancerTransfer{
			In:    &transferIn,
			Out:   &transferOut,
			Total: &transferTotal,
		},
	}

	nodeBalancerModel := &NodeBalancerModel{}

	vpcConfigs := []linodego.NodeBalancerVPCConfig{
		{
			ID:             123,
			NodeBalancerID: 456,
			SubnetID:       789,
			VPCID:          321,
			IPv4Range:      "10.0.0.4/30",
		},
	}

	diags := nodeBalancerModel.Flatten(
		context.Background(),
		nodeBalancer,
		nil,
		vpcConfigs,
		false,
	)

	assert.False(t, diags.HasError())

	assert.Equal(t, types.StringValue("123"), nodeBalancerModel.ID)
	assert.Equal(t, types.StringValue("us-east"), nodeBalancerModel.Region)
	assert.Equal(t, types.Int64Value(10), nodeBalancerModel.ClientConnThrottle)
	assert.Equal(t, types.Int64Value(5), nodeBalancerModel.ClientUDPSessThrottle)
	assert.Equal(t, types.StringPointerValue(&hostname), nodeBalancerModel.Hostname)
	assert.Equal(t, types.StringPointerValue(&IPv4), nodeBalancerModel.IPv4)
	assert.Equal(t, types.StringPointerValue(&IPv6), nodeBalancerModel.IPv6)

	assert.NotNil(t, nodeBalancerModel.Created)
	assert.NotNil(t, nodeBalancerModel.Updated)

	assert.Contains(t, nodeBalancerModel.Transfer.String(), "100.0")
	assert.Contains(t, nodeBalancerModel.Transfer.String(), "200.0")
	assert.Contains(t, nodeBalancerModel.Transfer.String(), "300.0")

	var vpcConfigModel []ResourceVPCModel
	d := nodeBalancerModel.VPCs.ElementsAs(t.Context(), &vpcConfigModel, false)
	if d.HasError() {
		t.Fatal(d.Errors())
	}

	assert.Equal(t, types.Int64Value(789), vpcConfigModel[0].SubnetID)
	assert.Equal(t, types.StringValue("10.0.0.4/30"), vpcConfigModel[0].IPv4Range)

	assert.True(t, types.StringValue(label).Equal(nodeBalancerModel.Label))
}

func TestUpgradeResourceStateValue(t *testing.T) {
	t.Run("ValidFloatConversion", func(t *testing.T) {
		value := "42.5"
		result, diag := UpgradeResourceStateValue(value)

		assert.Empty(t, diag)
		assert.Equal(t, "42.500000", result.String())
	})

	t.Run("EmptyValue", func(t *testing.T) {
		value := ""
		result, diag := UpgradeResourceStateValue(value)

		assert.Empty(t, diag)
		assert.Equal(t, "0.000000", result.String())
	})

	t.Run("InvalidFloatConversion", func(t *testing.T) {
		value := "invalid"
		result, diag := UpgradeResourceStateValue(value)

		fmt.Println(diag.Detail())
		assert.Contains(t, diag.Detail(), "strconv.ParseFloat: parsing \"invalid\": invalid syntax")
		assert.Empty(t, result)
	})
}
