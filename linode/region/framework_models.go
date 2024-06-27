package region

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

// RegionPGLimitsModel represents the placement group limits for a .
type RegionPGLimitsModel struct {
	MaximumPGsPerCustomer types.Int64 `tfsdk:"maximum_pgs_per_customer"`
	MaximumLinodesPerPG   types.Int64 `tfsdk:"maximum_linodes_per_pg"`
}

// RegionResolversModel represents a region's resolver.
type RegionResolversModel struct {
	IPv4 types.String `tfsdk:"ipv4"`
	IPv6 types.String `tfsdk:"ipv6"`
}

// RegionModel represents a single region object.
type RegionModel struct {
	Country              types.String           `tfsdk:"country"`
	ID                   types.String           `tfsdk:"id"`
	Label                types.String           `tfsdk:"label"`
	SiteType             types.String           `tfsdk:"site_type"`
	Status               types.String           `tfsdk:"status"`
	Capabilities         []types.String         `tfsdk:"capabilities"`
	Resolvers            []RegionResolversModel `tfsdk:"resolvers"`
	PlacementGroupLimits []RegionPGLimitsModel  `tfsdk:"placement_group_limits"`
}

func (m *RegionModel) ParseRegion(region *linodego.Region) {
	m.ID = types.StringValue(region.ID)
	m.Label = types.StringValue(region.Label)
	m.Status = types.StringValue(region.Status)
	m.Country = types.StringValue(region.Country)
	m.SiteType = types.StringValue(region.SiteType)

	m.Capabilities = helper.StringSliceToFramework(region.Capabilities)

	m.Resolvers = []RegionResolversModel{
		{
			IPv4: types.StringValue(region.Resolvers.IPv4),
			IPv6: types.StringValue(region.Resolvers.IPv6),
		},
	}

	regionLimits := region.PlacementGroupLimits

	if regionLimits != nil {
		m.PlacementGroupLimits = []RegionPGLimitsModel{
			{
				MaximumPGsPerCustomer: types.Int64Value(int64(regionLimits.MaximumPGsPerCustomer)),
				MaximumLinodesPerPG:   types.Int64Value(int64(regionLimits.MaximumLinodesPerPG)),
			},
		}
	}
}
