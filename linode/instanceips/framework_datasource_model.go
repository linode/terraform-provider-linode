package instanceips

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type InstanceIPDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	IPv4 types.Object `tfsdk:"ipv4"`
	IPv6 types.Object `tfsdk:"ipv6"`
}
