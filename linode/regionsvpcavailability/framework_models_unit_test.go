//go:build unit

package regionsvpcavailability

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseRegionsVPCAvailability(t *testing.T) {
	regions := []linodego.RegionVPCAvailability{
		{
			Region:                     "us-east",
			Available:                  true,
			AvailableIPV6PrefixLengths: []int{64, 128},
		},
		{
			Region:                     "eu-west",
			Available:                  false,
			AvailableIPV6PrefixLengths: []int{64},
		},
		{
			Region:                     "eu-east",
			Available:                  false,
			AvailableIPV6PrefixLengths: []int{},
		},
	}

	model := &regionsVPCAvailabilityModel{}

	model.parseRegionsVPCAvailability(context.Background(), regions)

	assert.Len(t, model.RegionsVPCAvailability, len(regions))

	for i, expectedRegion := range regions {
		assert.Equal(t, types.StringValue(expectedRegion.Region), model.RegionsVPCAvailability[i].ID)
		assert.Equal(t, types.BoolValue(expectedRegion.Available), model.RegionsVPCAvailability[i].Available)
		availableIPV6PrefixLengths, diags := types.ListValueFrom(context.Background(), types.Int64Type, expectedRegion.AvailableIPV6PrefixLengths)
		if diags.HasError() {
			t.Error(diags.Errors())
		}
		assert.Equal(t, availableIPV6PrefixLengths, model.RegionsVPCAvailability[i].AvailableIPV6PrefixLengths)
	}
}
