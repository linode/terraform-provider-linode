package nbconfigs

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v2/linode/nbconfig"
)

// NodeBalancerConfigFilterModel describes the Terraform resource data model to match the
// resource schema.
type NodeBalancerConfigFilterModel struct {
	ID                  types.String                     `tfsdk:"id"`
	Filters             frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order               types.String                     `tfsdk:"order"`
	OrderBy             types.String                     `tfsdk:"order_by"`
	NodeBalancerConfigs []nbconfig.DataSourceModel       `tfsdk:"nodebalancer_configs"`
}

func (data *NodeBalancerConfigFilterModel) parseNodeBalancerConfigs(
	nodebalancerConfigs []linodego.NodeBalancerConfig,
) {
	result := make([]nbconfig.DataSourceModel, len(nodebalancerConfigs))
	for i := range nodebalancerConfigs {
		var nbConfigData nbconfig.DataSourceModel
		nbConfigData.ParseNodebalancerConfig(&nodebalancerConfigs[i])
		result[i] = nbConfigData
	}

	data.NodeBalancerConfigs = result
}
