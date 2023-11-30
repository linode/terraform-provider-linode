package nbs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v2/linode/nb"
)

// NodeBalancerFilterModel describes the Terraform resource data model to match the
// resource schema.
type NodeBalancerFilterModel struct {
	ID            types.String                     `tfsdk:"id"`
	Filters       frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order         types.String                     `tfsdk:"order"`
	OrderBy       types.String                     `tfsdk:"order_by"`
	NodeBalancers []nb.NodeBalancerDataSourceModel `tfsdk:"nodebalancers"`
}

func (data *NodeBalancerFilterModel) parseNodeBalancers(
	ctx context.Context,
	nodebalancers []linodego.NodeBalancer,
) {
	result := make([]nb.NodeBalancerDataSourceModel, len(nodebalancers))
	for i := range nodebalancers {
		var nbData nb.NodeBalancerDataSourceModel
		nbData.FlattenNodeBalancer(ctx, &nodebalancers[i])
		result[i] = nbData
	}

	data.NodeBalancers = result
}
