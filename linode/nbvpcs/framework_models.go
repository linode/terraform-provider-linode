package nbvpcs

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v3/linode/nbvpc"
)

// Model describes the Terraform resource data model to match the
// resource schema.
type Model struct {
	NodeBalancerID types.Int64 `tfsdk:"nodebalancer_id"`

	ID         types.String                     `tfsdk:"id"`
	Filters    frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order      types.String                     `tfsdk:"order"`
	OrderBy    types.String                     `tfsdk:"order_by"`
	VPCConfigs []nbvpc.DataSourceModel          `tfsdk:"vpc_configs"`
}

func (data *Model) Parse(
	vpcConfigs []linodego.NodeBalancerVPCConfig,
) {
	data.VPCConfigs = helper.MapSlice(
		vpcConfigs,
		func(vpcConfig linodego.NodeBalancerVPCConfig) nbvpc.DataSourceModel {
			var result nbvpc.DataSourceModel
			return *result.Flatten(&vpcConfig)
		},
	)
}
