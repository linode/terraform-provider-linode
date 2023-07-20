package nbs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/linode/nb"
)

// NodeBalancerFilterModel describes the Terraform resource data model to match the
// resource schema.
type NodeBalancerFilterModel struct {
	ID            types.String                     `tfsdk:"id"`
	Filters       frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order         types.String                     `tfsdk:"order"`
	OrderBy       types.String                     `tfsdk:"order_by"`
	NodeBalancers []nb.DataSourceModel             `tfsdk:"nodebalancers"`
}

func (data *NodeBalancerFilterModel) parseNodeBalancers(
	ctx context.Context,
	nodebalancers []linodego.NodeBalancer,
) {
	result := make([]nb.DataSourceModel, len(nodebalancers))
	for i := range nodebalancers {
		var nbData nb.DataSourceModel
		nbData.ParseNodeBalancer(ctx, &nodebalancers[i])
		result[i] = nbData
	}

	data.NodeBalancers = result
}
