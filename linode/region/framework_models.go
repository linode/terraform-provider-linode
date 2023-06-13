package region

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

// RegionResolversModel represents a region's resolver.
type RegionResolversModel struct {
	IPv4 types.String `tfsdk:"ipv4"`
	IPv6 types.String `tfsdk:"ipv6"`
}

// RegionModel represents a single region object.
type RegionModel struct {
	Country      types.String           `tfsdk:"country"`
	ID           types.String           `tfsdk:"id"`
	Label        types.String           `tfsdk:"label"`
	Capabilities []types.String         `tfsdk:"capabilities"`
	Status       types.String           `tfsdk:"status"`
	Resolvers    []RegionResolversModel `tfsdk:"resolvers"`
}

func (m *RegionModel) parseRegion(region *linodego.Region) {
	m.ID = types.StringValue(region.ID)
	m.Label = types.StringValue(region.Label)
	m.Status = types.StringValue(region.Status)
	m.Country = types.StringValue(region.Country)

	m.Capabilities = helper.StringSliceToFramework(region.Capabilities)

	m.Resolvers = []RegionResolversModel{
		{
			IPv4: types.StringValue(region.Resolvers.IPv4),
			IPv6: types.StringValue(region.Resolvers.IPv6),
		},
	}
}
