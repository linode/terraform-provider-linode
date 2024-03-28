package regions

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
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
	SiteType     types.String           `tfsdk:"site_type"`
	Status       types.String           `tfsdk:"status"`
	Capabilities []types.String         `tfsdk:"capabilities"`
	Resolvers    []RegionResolversModel `tfsdk:"resolvers"`
}

// RegionFilterModel describes the Terraform resource data model to match the
// resource schema.
type RegionFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Regions []RegionModel                    `tfsdk:"regions"`
}

// parseRegions parses the given list of regions into the `regions` model attribute.
func (model *RegionFilterModel) parseRegions(regions []linodego.Region) {
	parseRegion := func(region linodego.Region) RegionModel {
		var m RegionModel
		m.ID = types.StringValue(region.ID)
		m.Label = types.StringValue(region.Label)
		m.Status = types.StringValue(region.Status)
		m.Country = types.StringValue(region.Country)
		m.SiteType = types.StringValue(region.SiteType)

		m.Capabilities = make([]types.String, len(region.Capabilities))
		for k, c := range region.Capabilities {
			m.Capabilities[k] = types.StringValue(c)
		}

		m.Resolvers = []RegionResolversModel{
			{
				IPv4: types.StringValue(region.Resolvers.IPv4),
				IPv6: types.StringValue(region.Resolvers.IPv6),
			},
		}

		return m
	}

	result := make([]RegionModel, len(regions))

	for i, region := range regions {
		result[i] = parseRegion(region)
	}

	model.Regions = result
}
