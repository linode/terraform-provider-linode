package instancetypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/linode/instancetype"
)

type InstanceTypeFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Order   types.String                     `tfsdk:"order"`
	OrderBy types.String                     `tfsdk:"order_by"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Types   []instancetype.DataSourceModel   `tfsdk:"types"`
}

func (model *InstanceTypeFilterModel) parseInstanceTypes(ctx context.Context,
	instanceTypes []linodego.LinodeType,
) diag.Diagnostics {
	result := make([]instancetype.DataSourceModel, len(instanceTypes))

	for i, instanceType := range instanceTypes {
		var m instancetype.DataSourceModel

		var diags = m.ParseLinodeType(ctx, &instanceType)
		if diags.HasError() {
			return diags
		}

		result[i] = m
	}

	model.Types = result

	return nil
}
