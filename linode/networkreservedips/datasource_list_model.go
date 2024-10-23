package networkreservedips

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DataSourceListModel struct {
	ReservedIPs types.List `tfsdk:"reserved_ips"`
}
