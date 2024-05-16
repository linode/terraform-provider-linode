//go:build unit

package region

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseRegion(t *testing.T) {
	region := &linodego.Region{
		ID:           "us-east",
		Country:      "us",
		Capabilities: []string{"Linodes", "NodeBalancers", "Block Storage", "Object Storage"},
		Status:       "ok",
		SiteType:     "core",
		Resolvers: linodego.RegionResolvers{
			IPv4: "192.0.2.0,192.0.2.1",
			IPv6: "2001:0db8::,2001:0db8::1",
		},
		PlacementGroupLimits: &linodego.RegionPlacementGroupLimits{
			MaximumLinodesPerPG:   5,
			MaximumPGsPerCustomer: 10,
		},
		Label: "Newark, NJ, USA",
	}

	model := &RegionModel{}

	model.ParseRegion(region)

	assert.Equal(t, types.StringValue("us-east"), model.ID)
	assert.Equal(t, types.StringValue("Newark, NJ, USA"), model.Label)
	assert.Equal(t, types.StringValue("ok"), model.Status)
	assert.Equal(t, types.StringValue("us"), model.Country)
	assert.Equal(t, types.StringValue("core"), model.SiteType)

	for i, capability := range region.Capabilities {
		assert.Equal(t, types.StringValue(capability), model.Capabilities[i])
	}

	assert.Equal(t, types.StringValue("192.0.2.0,192.0.2.1"), model.Resolvers[0].IPv4)
	assert.Equal(t, types.StringValue("2001:0db8::,2001:0db8::1"), model.Resolvers[0].IPv6)

	assert.Equal(t, types.Int64Value(5), model.PlacementGroupLimits[0].MaximumLinodesPerPG)
	assert.Equal(t, types.Int64Value(10), model.PlacementGroupLimits[0].MaximumPGsPerCustomer)
}
