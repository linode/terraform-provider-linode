//go:build unit

package lke

import (
	"context"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

// minimalPoolNetworkingSchema returns a minimal schema resource with tier and a
// pool list whose only sub-attribute is a networking block. Used to test custom
// diff validators without pulling in the full resource schema.
func minimalPoolNetworkingSchema(customDiff schema.CustomizeDiffFunc) *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"tier": {Type: schema.TypeString, Optional: true},
			"pool": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"networking": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"interfaces": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type":      {Type: schema.TypeString, Required: true},
												"vpc_id":    {Type: schema.TypeInt, Optional: true},
												"subnet_id": {Type: schema.TypeInt, Optional: true},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		CustomizeDiff: customDiff,
	}
}

func buildPoolNetworkingCtyConfig(t *testing.T, r *schema.Resource, config map[string]any) cty.Value {
	t.Helper()
	return networkingTestConfig(t, r, config)
}

var testPoolWithNetworking = []any{map[string]any{
	"networking": []any{map[string]any{
		"interfaces": []any{map[string]any{
			"type":      "vpc",
			"vpc_id":    1,
			"subnet_id": 10,
		}},
	}},
}}

var testPoolWithoutNetworking = []any{map[string]any{
	"networking": []any{},
}}

func TestCustomDiffValidatePoolNetworkingRequiresEnterprise(t *testing.T) {
	for _, tc := range []struct {
		name        string
		tier        string
		pool        []any
		expectError bool
	}{
		{
			name:        "enterprise tier allows pool networking",
			tier:        TierEnterprise,
			pool:        testPoolWithNetworking,
			expectError: false,
		},
		{
			name:        "standard tier rejects pool networking",
			tier:        "standard",
			pool:        testPoolWithNetworking,
			expectError: true,
		},
		{
			name:        "omitted tier rejects pool networking",
			tier:        "",
			pool:        testPoolWithNetworking,
			expectError: true,
		},
		{
			name:        "enterprise tier without pool networking is fine",
			tier:        TierEnterprise,
			pool:        testPoolWithoutNetworking,
			expectError: false,
		},
		{
			name:        "standard tier without pool networking is fine",
			tier:        "standard",
			pool:        testPoolWithoutNetworking,
			expectError: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := minimalPoolNetworkingSchema(customDiffValidatePoolNetworkingRequiresEnterprise)
			config := map[string]any{"pool": tc.pool}
			if tc.tier != "" {
				config["tier"] = tc.tier
			}
			state := &terraform.InstanceState{RawConfig: buildPoolNetworkingCtyConfig(t, r, config)}
			_, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), nil)
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "pool networking is only available for LKE Enterprise clusters")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCustomDiffValidatePoolNetworkingRequiresBeta_UnknownDoesNotPanic(t *testing.T) {
	r := minimalPoolNetworkingSchema(customDiffValidatePoolNetworkingRequiresBeta)

	// Build a cty config where the pool's networking attribute is unknown.
	// This simulates a dynamic block with an unknown for_each during planning.
	poolType := r.CoreConfigSchema().ImpliedType()
	poolElemType := poolType.AttributeType("pool").ElementType()
	netType := poolElemType.AttributeType("networking")

	poolVal := cty.ListVal([]cty.Value{
		cty.ObjectVal(map[string]cty.Value{
			"networking": cty.UnknownVal(netType),
		}),
	})
	rawConfig := cty.ObjectVal(map[string]cty.Value{
		"tier": cty.NullVal(cty.String),
		"pool": poolVal,
	})

	// Use a non-beta API version so the validator iterates pools (and would
	// previously panic on unknown networking before the IsKnown guard was added).
	meta := &helper.ProviderMeta{Config: &helper.Config{APIVersion: "v4"}}
	state := &terraform.InstanceState{RawConfig: rawConfig}
	assert.NotPanics(t, func() {
		_, _ = r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(map[string]any{}), meta)
	})
}

func TestCustomDiffValidatePoolNetworkingFirstInterface_UnknownDoesNotPanic(t *testing.T) {
	// customDiffValidatePoolNetworkingFirstInterface does not use meta.
	r := minimalPoolNetworkingSchema(customDiffValidatePoolNetworkingFirstInterface)

	poolType := r.CoreConfigSchema().ImpliedType()
	poolElemType := poolType.AttributeType("pool").ElementType()
	netType := poolElemType.AttributeType("networking")

	poolVal := cty.ListVal([]cty.Value{
		cty.ObjectVal(map[string]cty.Value{
			"networking": cty.UnknownVal(netType),
		}),
	})
	rawConfig := cty.ObjectVal(map[string]cty.Value{
		"tier": cty.NullVal(cty.String),
		"pool": poolVal,
	})

	state := &terraform.InstanceState{RawConfig: rawConfig}
	assert.NotPanics(t, func() {
		_, _ = r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(map[string]any{}), nil)
	})
}
