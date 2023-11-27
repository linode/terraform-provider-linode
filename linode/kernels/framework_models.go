package kernels

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v2/linode/kernel"
)

// KernelFilterModel describes the Terraform resource data model to match the
// resource schema.
type KernelFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order   types.String                     `tfsdk:"order"`
	OrderBy types.String                     `tfsdk:"order_by"`
	Kernels []kernel.DataSourceModel         `tfsdk:"kernels"`
}

func (data *KernelFilterModel) parseKernels(
	ctx context.Context,
	kernels []linodego.LinodeKernel,
) {
	result := make([]kernel.DataSourceModel, len(kernels))
	for i := range kernels {
		var kernelData kernel.DataSourceModel
		kernelData.ParseKernel(ctx, &kernels[i])
		result[i] = kernelData
	}

	data.Kernels = result
}
