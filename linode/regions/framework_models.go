package regions

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RegionFilterModel describes the Terraform resource data model to match the
// resource schema.
type RegionFilterModel struct {
	ID      types.String `tfsdk:"id"`
	Filters types.Set    `tfsdk:"filter"`
}
