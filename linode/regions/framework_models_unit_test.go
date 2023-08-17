package regions

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseRegions(t *testing.T) {
	regions := []linodego.Region{
		{
			ID:           "us-east",
			Label:        "Newark, NJ, USA",
			Status:       "ok",
			Country:      "us",
			Capabilities: []string{"Linodes", "NodeBalancers", "Block Storage", "Object Storage"},
			Resolvers: linodego.RegionResolvers{
				IPv4: "192.0.2.0,192.0.2.1",
				IPv6: "2001:0db8::,2001:0db8::1",
			},
		},
		{
			ID:           "ap-west",
			Label:        "Different label",
			Status:       "ok",
			Country:      "us",
			Capabilities: []string{"Linodes", "NodeBalancers", "Block Storage", "Object Storage"},
			Resolvers: linodego.RegionResolvers{
				IPv4: "192.0.2.4,192.0.2.3",
				IPv6: "2001:0db8::,2001:0db8::3",
			},
		},
	}

	model := &RegionFilterModel{}

	model.parseRegions(regions)

	assert.Len(t, model.Regions, len(regions))

	for i, expectedRegion := range regions {
		assert.Equal(t, types.StringValue(expectedRegion.ID), model.Regions[i].ID)
		assert.Equal(t, types.StringValue(expectedRegion.Label), model.Regions[i].Label)
		assert.Equal(t, types.StringValue(expectedRegion.Status), model.Regions[i].Status)
		assert.Equal(t, types.StringValue(expectedRegion.Country), model.Regions[i].Country)
		for j, capability := range regions[i].Capabilities {
			assert.Equal(t, types.StringValue(capability), model.Regions[i].Capabilities[j])
		}
		assert.Equal(t, types.StringValue(expectedRegion.Resolvers.IPv4), model.Regions[i].Resolvers[0].IPv4)
		assert.Equal(t, types.StringValue(expectedRegion.Resolvers.IPv6), model.Regions[i].Resolvers[0].IPv6)
	}
}
