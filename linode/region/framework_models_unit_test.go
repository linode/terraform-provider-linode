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
		Resolvers: linodego.RegionResolvers{
			IPv4: "192.0.2.0,192.0.2.1",
			IPv6: "2001:0db8::,2001:0db8::1",
		},
		Label: "Newark, NJ, USA",
	}

	model := &RegionModel{}

	model.parseRegion(region)

	assert.Equal(t, types.StringValue("us-east"), model.ID)
	assert.Equal(t, types.StringValue("Newark, NJ, USA"), model.Label)
	assert.Equal(t, types.StringValue("ok"), model.Status)
	assert.Equal(t, types.StringValue("us"), model.Country)

	for i, capability := range region.Capabilities {
		assert.Equal(t, types.StringValue(capability), model.Capabilities[i])
	}

	assert.Equal(t, types.StringValue("192.0.2.0,192.0.2.1"), model.Resolvers[0].IPv4)
	assert.Equal(t, types.StringValue("2001:0db8::,2001:0db8::1"), model.Resolvers[0].IPv6)
}
