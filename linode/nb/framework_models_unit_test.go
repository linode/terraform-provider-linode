//go:build unit

package nb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
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
		Type: linodego.NBTypeCommon,
		LKECluster: &linodego.NodeBalancerLKECluster{
			ID:    1234,
			Label: "test-cluster",
			Type:  "lkecluster",
			URL:   "/v4/lke/clusters/1234",
		},
	}

	nodeBalancerModel := &NodeBalancerModel{
		VPCs: types.ListValueMust(backendVPCObjType, []attr.Value{
			types.ObjectValueMust(backendVPCObjType.AttrTypes, map[string]attr.Value{
				"subnet_id":              types.Int64Value(789),
				"ipv4_range":             types.StringUnknown(),
				"ipv6_range":             types.StringNull(),
				"ipv4_range_auto_assign": types.BoolValue(true),
			}),
		}),
		BackendVPCs: types.ListValueMust(backendVPCObjType, []attr.Value{}),
	}

	vpcConfigs := []linodego.NodeBalancerVPCConfig{
		{
			ID:             123,
			NodeBalancerID: 456,
			SubnetID:       789,
			VPCID:          321,
			IPv4Range:      "10.0.0.4/30",
			Purpose:        "backend",
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
	assert.Equal(t, types.StringValue("common"), nodeBalancerModel.Type)
	assert.Equal(t, types.StringValue("us-east"), nodeBalancerModel.Region)
	assert.Equal(t, types.Int64Value(10), nodeBalancerModel.ClientConnThrottle)
	assert.Equal(t, types.Int64Value(5), nodeBalancerModel.ClientUDPSessThrottle)
	assert.Equal(t, types.StringPointerValue(&hostname), nodeBalancerModel.Hostname)
	assert.Equal(t, iptypes.NewIPv4AddressPointerValue(&IPv4), nodeBalancerModel.IPv4)
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

	var lkeClusters []LKEClusterModel
	d = nodeBalancerModel.LKECluster.ElementsAs(t.Context(), &lkeClusters, false)
	if d.HasError() {
		t.Fatal(d.Errors())
	}
	assert.Len(t, lkeClusters, 1)
	assert.Equal(t, types.Int64Value(1234), lkeClusters[0].ID)
	assert.Equal(t, types.StringValue("test-cluster"), lkeClusters[0].Label)
	assert.Equal(t, types.StringValue("lkecluster"), lkeClusters[0].Type)
	assert.Equal(t, types.StringValue("/v4/lke/clusters/1234"), lkeClusters[0].URL)
}

func TestFlattenLKECluster(t *testing.T) {
	t.Run("WithLKECluster", func(t *testing.T) {
		lkeCluster := &linodego.NodeBalancerLKECluster{
			ID:    1234,
			Label: "test-cluster",
			Type:  "lkecluster",
			URL:   "/v4/lke/clusters/1234",
		}

		result, diags := FlattenLKECluster(context.Background(), lkeCluster)
		assert.False(t, diags.HasError())
		assert.Len(t, result.Elements(), 1)

		var models []LKEClusterModel
		d := result.ElementsAs(context.Background(), &models, false)
		assert.False(t, d.HasError())
		assert.Equal(t, types.Int64Value(1234), models[0].ID)
		assert.Equal(t, types.StringValue("test-cluster"), models[0].Label)
		assert.Equal(t, types.StringValue("lkecluster"), models[0].Type)
		assert.Equal(t, types.StringValue("/v4/lke/clusters/1234"), models[0].URL)
	})

	t.Run("NilLKECluster", func(t *testing.T) {
		result, diags := FlattenLKECluster(context.Background(), nil)
		assert.False(t, diags.HasError())
		assert.Len(t, result.Elements(), 0)
	})
}

func TestFlattenNodeBalancerIPv4(t *testing.T) {
	reservedIP := "198.51.100.5"
	nodeBalancer := &linodego.NodeBalancer{
		ID:   456,
		IPv4: &reservedIP,
	}

	t.Run("preserves user-provided IPv4 when preserveKnown=true", func(t *testing.T) {
		userProvidedIP := "203.0.113.7"
		model := &NodeBalancerModel{
			IPv4: iptypes.NewIPv4AddressValue(userProvidedIP),
		}

		diags := model.Flatten(context.Background(), nodeBalancer, nil, nil, true)
		assert.False(t, diags.HasError())
		assert.Equal(t, iptypes.NewIPv4AddressValue(userProvidedIP), model.IPv4)
	})

	t.Run("updates IPv4 from API when preserveKnown=false", func(t *testing.T) {
		apiIP := "203.0.113.7"
		model := &NodeBalancerModel{
			IPv4: iptypes.NewIPv4AddressValue(apiIP),
		}

		nb := &linodego.NodeBalancer{ID: 457, IPv4: &reservedIP}
		diags := model.Flatten(context.Background(), nb, nil, nil, false)
		assert.False(t, diags.HasError())
		assert.Equal(t, iptypes.NewIPv4AddressValue(reservedIP), model.IPv4)
	})

	t.Run("sets IPv4 from API when model IPv4 is unknown", func(t *testing.T) {
		model := &NodeBalancerModel{
			IPv4: iptypes.NewIPv4AddressUnknown(),
		}

		diags := model.Flatten(context.Background(), nodeBalancer, nil, nil, true)
		assert.False(t, diags.HasError())
		assert.Equal(t, iptypes.NewIPv4AddressValue(reservedIP), model.IPv4)
	})
}

func TestFlattenNodeBalancerBackendVPCsPopulateComputedRanges(t *testing.T) {
	nodeBalancer := &linodego.NodeBalancer{ID: 123}

	nodeBalancerModel := &NodeBalancerModel{
		BackendVPCs: types.ListValueMust(backendVPCObjType, []attr.Value{
			types.ObjectValueMust(backendVPCObjType.AttrTypes, map[string]attr.Value{
				"subnet_id":              types.Int64Value(789),
				"ipv4_range":             types.StringUnknown(),
				"ipv6_range":             types.StringValue("fd00::/64"),
				"ipv4_range_auto_assign": types.BoolValue(true),
			}),
		}),
	}

	vpcConfigs := []linodego.NodeBalancerVPCConfig{
		{
			ID:             123,
			NodeBalancerID: 456,
			SubnetID:       789,
			IPv4Range:      "10.0.0.4/30",
			IPv6Range:      "fd00::4/126",
			Purpose:        linodego.NodeBalancerVPCConfigPurposeBackend,
		},
		{
			ID:             124,
			NodeBalancerID: 456,
			SubnetID:       790,
			IPv4Range:      "10.0.0.8/30",
			Purpose:        linodego.NodeBalancerVPCConfigPurposeBackend,
		},
	}

	diags := nodeBalancerModel.Flatten(
		context.Background(),
		nodeBalancer,
		nil,
		vpcConfigs,
		true,
	)

	assert.False(t, diags.HasError())

	var backendVPCModels []ResourceBackendVPCModel
	d := nodeBalancerModel.BackendVPCs.ElementsAs(context.Background(), &backendVPCModels, false)
	if d.HasError() {
		t.Fatal(d.Errors())
	}

	assert.Len(t, backendVPCModels, 1)
	assert.Equal(t, types.Int64Value(789), backendVPCModels[0].SubnetID)
	assert.Equal(t, types.StringValue("10.0.0.4/30"), backendVPCModels[0].IPv4Range)
	assert.Equal(t, types.StringValue("fd00::4/126"), backendVPCModels[0].IPv6Range)
	assert.Equal(t, types.BoolValue(true), backendVPCModels[0].IPv4RangeAutoAssign)

	options, optionDiags := backendVPCModels[0].ToLinodego()
	assert.False(t, optionDiags.HasError())
	assert.Equal(t, "fd00::4/126", options.IPv6Range)
}

func TestFlattenNodeBalancerBackendVPCsReconcileAPIEntries(t *testing.T) {
	nodeBalancer := &linodego.NodeBalancer{ID: 123}

	nodeBalancerModel := &NodeBalancerModel{
		BackendVPCs: types.ListValueMust(backendVPCObjType, []attr.Value{
			types.ObjectValueMust(backendVPCObjType.AttrTypes, map[string]attr.Value{
				"subnet_id":              types.Int64Value(788),
				"ipv4_range":             types.StringUnknown(),
				"ipv6_range":             types.StringNull(),
				"ipv4_range_auto_assign": types.BoolValue(true),
			}),
			types.ObjectValueMust(backendVPCObjType.AttrTypes, map[string]attr.Value{
				"subnet_id":              types.Int64Value(789),
				"ipv4_range":             types.StringUnknown(),
				"ipv6_range":             types.StringNull(),
				"ipv4_range_auto_assign": types.BoolValue(true),
			}),
		}),
	}

	vpcConfigs := []linodego.NodeBalancerVPCConfig{
		{
			SubnetID:  790,
			IPv4Range: "10.0.0.8/30",
			Purpose:   linodego.NodeBalancerVPCConfigPurposeBackend,
		},
		{
			SubnetID:  789,
			IPv4Range: "10.0.0.4/30",
			Purpose:   linodego.NodeBalancerVPCConfigPurposeBackend,
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

	var backendVPCModels []ResourceBackendVPCModel
	d := nodeBalancerModel.BackendVPCs.ElementsAs(context.Background(), &backendVPCModels, false)
	if d.HasError() {
		t.Fatal(d.Errors())
	}

	assert.Len(t, backendVPCModels, 2)
	assert.Equal(t, types.Int64Value(789), backendVPCModels[0].SubnetID)
	assert.Equal(t, types.StringValue("10.0.0.4/30"), backendVPCModels[0].IPv4Range)
	assert.Equal(t, types.BoolValue(true), backendVPCModels[0].IPv4RangeAutoAssign)
	assert.Equal(t, types.Int64Value(790), backendVPCModels[1].SubnetID)
	assert.Equal(t, types.StringValue("10.0.0.8/30"), backendVPCModels[1].IPv4Range)
	assert.True(t, backendVPCModels[1].IPv4RangeAutoAssign.IsNull())
}

func TestFlattenNodeBalancerFrontendVPCsPreserveConfiguredRanges(t *testing.T) {
	nodeBalancer := &linodego.NodeBalancer{
		ID: 123,
	}

	nodeBalancerModel := &NodeBalancerModel{
		FrontendVPCs: types.ListValueMust(frontendVPCObjType, []attr.Value{
			types.ObjectValueMust(frontendVPCObjType.AttrTypes, map[string]attr.Value{
				"subnet_id":            types.Int64Value(789),
				"ipv4_range":           types.StringValue("auto"),
				"allocated_ipv4_range": types.StringUnknown(),
				"ipv6_range":           types.StringValue("auto"),
				"allocated_ipv6_range": types.StringUnknown(),
			}),
		}),
	}

	vpcConfigs := []linodego.NodeBalancerVPCConfig{
		{
			ID:             123,
			NodeBalancerID: 456,
			SubnetID:       789,
			IPv4Range:      "10.0.0.4/30",
			IPv6Range:      "fd00::4/126",
			Purpose:        linodego.NodeBalancerVPCConfigPurposeFrontend,
		},
	}

	diags := nodeBalancerModel.Flatten(
		context.Background(),
		nodeBalancer,
		nil,
		vpcConfigs,
		true,
	)

	assert.False(t, diags.HasError())

	var frontendVPCModels []ResourceFrontendVPCModel
	d := nodeBalancerModel.FrontendVPCs.ElementsAs(context.Background(), &frontendVPCModels, false)
	if d.HasError() {
		t.Fatal(d.Errors())
	}

	assert.Len(t, frontendVPCModels, 1)
	assert.Equal(t, types.Int64Value(789), frontendVPCModels[0].SubnetID)
	assert.Equal(t, types.StringValue("auto"), frontendVPCModels[0].IPv4Range)
	assert.Equal(t, types.StringValue("10.0.0.4/30"), frontendVPCModels[0].AllocatedIPv4Range)
	assert.Equal(t, types.StringValue("auto"), frontendVPCModels[0].IPv6Range)
	assert.Equal(t, types.StringValue("fd00::4/126"), frontendVPCModels[0].AllocatedIPv6Range)
}

func TestFlattenNodeBalancerFrontendVPCsReconcileAPIEntries(t *testing.T) {
	nodeBalancer := &linodego.NodeBalancer{ID: 123}

	nodeBalancerModel := &NodeBalancerModel{
		FrontendVPCs: types.ListValueMust(frontendVPCObjType, []attr.Value{
			types.ObjectValueMust(frontendVPCObjType.AttrTypes, map[string]attr.Value{
				"subnet_id":            types.Int64Value(788),
				"ipv4_range":           types.StringValue("auto"),
				"allocated_ipv4_range": types.StringUnknown(),
				"ipv6_range":           types.StringValue("auto"),
				"allocated_ipv6_range": types.StringUnknown(),
			}),
			types.ObjectValueMust(frontendVPCObjType.AttrTypes, map[string]attr.Value{
				"subnet_id":            types.Int64Value(789),
				"ipv4_range":           types.StringValue("auto"),
				"allocated_ipv4_range": types.StringUnknown(),
				"ipv6_range":           types.StringValue("auto"),
				"allocated_ipv6_range": types.StringUnknown(),
			}),
		}),
	}

	vpcConfigs := []linodego.NodeBalancerVPCConfig{
		{
			ID:             123,
			NodeBalancerID: 456,
			SubnetID:       790,
			IPv4Range:      "10.0.0.8/30",
			IPv6Range:      "fd00::8/126",
			Purpose:        linodego.NodeBalancerVPCConfigPurposeFrontend,
		},
		{
			ID:             124,
			NodeBalancerID: 456,
			SubnetID:       789,
			IPv4Range:      "10.0.0.4/30",
			IPv6Range:      "fd00::4/126",
			Purpose:        linodego.NodeBalancerVPCConfigPurposeFrontend,
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

	var frontendVPCModels []ResourceFrontendVPCModel
	d := nodeBalancerModel.FrontendVPCs.ElementsAs(context.Background(), &frontendVPCModels, false)
	if d.HasError() {
		t.Fatal(d.Errors())
	}

	assert.Len(t, frontendVPCModels, 2)
	assert.Equal(t, types.Int64Value(789), frontendVPCModels[0].SubnetID)
	assert.Equal(t, types.StringValue("auto"), frontendVPCModels[0].IPv4Range)
	assert.Equal(t, types.StringValue("10.0.0.4/30"), frontendVPCModels[0].AllocatedIPv4Range)
	assert.Equal(t, types.StringValue("auto"), frontendVPCModels[0].IPv6Range)
	assert.Equal(t, types.StringValue("fd00::4/126"), frontendVPCModels[0].AllocatedIPv6Range)
	assert.Equal(t, types.Int64Value(790), frontendVPCModels[1].SubnetID)
	assert.True(t, frontendVPCModels[1].IPv4Range.IsNull())
	assert.Equal(t, types.StringValue("10.0.0.8/30"), frontendVPCModels[1].AllocatedIPv4Range)
	assert.True(t, frontendVPCModels[1].IPv6Range.IsNull())
	assert.Equal(t, types.StringValue("fd00::8/126"), frontendVPCModels[1].AllocatedIPv6Range)
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
