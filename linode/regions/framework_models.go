package regions

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
)

type RegionResolversModel struct {
	IPv4 types.String `tfsdk:"ipv4"`
	IPv6 types.String `tfsdk:"ipv6"`
}

type RegionModel struct {
	Country      types.String           `tfsdk:"country"`
	ID           types.String           `tfsdk:"id"`
	Label        types.String           `tfsdk:"label"`
	Capabilities []types.String         `tfsdk:"capabilities"`
	Status       types.String           `tfsdk:"status"`
	Resolvers    []RegionResolversModel `tfsdk:"resolvers"`
}

// RegionFilterModel describes the Terraform resource data model to match the
// resource schema.
type RegionFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Regions []RegionModel                    `tfsdk:"regions"`
}
