package networkreservedips

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DataSourceModel struct {
	ReservedIPs types.List   `tfsdk:"reserved_ips"`
	Region      types.String `tfsdk:"region"`
}
