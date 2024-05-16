package regions

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
	regionResource "github.com/linode/terraform-provider-linode/v2/linode/region"
)

// RegionFilterModel describes the Terraform resource data model to match the
// resource schema.
type RegionFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Regions []regionResource.RegionModel     `tfsdk:"regions"`
}

// parseRegions parses the given list of regions into the `regions` model attribute.
func (model *RegionFilterModel) parseRegions(regions []linodego.Region) {
	result := make([]regionResource.RegionModel, len(regions))

	for i, region := range regions {
		model := regionResource.RegionModel{}
		model.ParseRegion(&region)
		result[i] = model
	}

	model.Regions = result
}
