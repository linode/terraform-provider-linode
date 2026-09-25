//go:build unit

package nb

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
		LKECluster: &linodego.NodeBalancerLKECluster{
			ID:    1234,
			Label: "test-cluster",
			Type:  "lkecluster",
			URL:   "/v4/lke/clusters/1234",
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

func TestNodeBalancerCreateOptionsConnectivity(t *testing.T) {
	model := NodeBalancerModel{
		Region:              types.StringValue("us-east"),
		Type:                types.StringValue(string(linodego.NBTypePremium)),
		BackendConnectivity: types.StringValue(string(linodego.NBBackendConnectivityIPv6)),
	}

	opts := model.GetCreateOptions(10, 5)
	assert.Equal(t, linodego.NBTypePremium, opts.Type)
	if assert.NotNil(t, opts.BackendConnectivity) {
		assert.Equal(t, linodego.NBBackendConnectivityIPv6, *opts.BackendConnectivity)
	}
	payload, err := json.Marshal(opts)
	assert.NoError(t, err)
	assert.Contains(t, string(payload), `"type":"premium"`)
	assert.Contains(t, string(payload), `"backend_connectivity":"ipv6"`)

	model.Type = types.StringUnknown()
	model.BackendConnectivity = types.StringUnknown()
	opts = model.GetCreateOptions(10, 5)
	assert.Empty(t, opts.Type)
	assert.Nil(t, opts.BackendConnectivity)
	payload, err = json.Marshal(opts)
	assert.NoError(t, err)
	assert.NotContains(t, string(payload), `"type"`)
	assert.NotContains(t, string(payload), `"backend_connectivity"`)

	model.Type = types.StringNull()
	model.BackendConnectivity = types.StringNull()
	opts = model.GetCreateOptions(10, 5)
	assert.Empty(t, opts.Type)
	assert.Nil(t, opts.BackendConnectivity)
}

func TestNodeBalancerCreateOnlyFieldValidators(t *testing.T) {
	tests := []struct {
		field     string
		value     string
		wantError bool
	}{
		{"type", "common", false},
		{"type", "premium", false},
		{"type", "enterprise", false},
		{"type", "premium_40gb", true},
		{"backend_connectivity", "legacy", false},
		{"backend_connectivity", "ipv6", false},
		{"backend_connectivity", "vpc", false},
		{"backend_connectivity", "undefined", true},
		{"backend_connectivity", "ipv6_and_vpc", true},
	}

	for _, tt := range tests {
		t.Run(tt.field+"/"+tt.value, func(t *testing.T) {
			attr := frameworkResourceSchema.Attributes[tt.field].(schema.StringAttribute)
			var resp validator.StringResponse
			attr.Validators[0].ValidateString(t.Context(), validator.StringRequest{
				Path:        path.Root(tt.field),
				ConfigValue: types.StringValue(tt.value),
			}, &resp)
			assert.Equal(t, tt.wantError, resp.Diagnostics.HasError())
		})
	}
}

func TestFlattenNodeBalancerConnectivity(t *testing.T) {
	connectivity := linodego.NBBackendConnectivityIPv6
	nodeBalancer := &linodego.NodeBalancer{
		ID:                  123,
		Type:                linodego.NBTypePremium,
		BackendConnectivity: &connectivity,
	}

	model := &NodeBalancerModel{
		Type:                types.StringUnknown(),
		BackendConnectivity: types.StringUnknown(),
	}
	diags := model.Flatten(t.Context(), nodeBalancer, nil, nil, true)
	assert.False(t, diags.HasError())
	assert.Equal(t, types.StringValue("premium"), model.Type)
	assert.Equal(t, types.StringValue("ipv6"), model.BackendConnectivity)

	model.Type = types.StringValue("common")
	model.BackendConnectivity = types.StringValue("vpc")
	diags = model.Flatten(t.Context(), nodeBalancer, nil, nil, true)
	assert.False(t, diags.HasError())
	assert.Equal(t, types.StringValue("common"), model.Type)
	assert.Equal(t, types.StringValue("vpc"), model.BackendConnectivity)

	diags = model.Flatten(t.Context(), nodeBalancer, nil, nil, false)
	assert.False(t, diags.HasError())
	assert.Equal(t, types.StringValue("premium"), model.Type)
	assert.Equal(t, types.StringValue("ipv6"), model.BackendConnectivity)

	nodeBalancer.BackendConnectivity = nil
	diags = model.Flatten(t.Context(), nodeBalancer, nil, nil, false)
	assert.False(t, diags.HasError())
	assert.True(t, model.BackendConnectivity.IsNull())
}

func TestFlattenNodeBalancerDataSourceConnectivity(t *testing.T) {
	connectivity := linodego.NBBackendConnectivityVPC
	model := &NodeBalancerDataSourceModel{}
	diags := model.Flatten(t.Context(), &linodego.NodeBalancer{
		Type:                linodego.NBTypeEnterprise,
		BackendConnectivity: &connectivity,
	}, nil, nil)
	assert.False(t, diags.HasError())
	assert.Equal(t, types.StringValue("enterprise"), model.Type)
	assert.Equal(t, types.StringValue("vpc"), model.BackendConnectivity)
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
