package volumes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v2/linode/volume"
)

type VolumeFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order   types.String                     `tfsdk:"order"`
	OrderBy types.String                     `tfsdk:"order_by"`
	Volumes []volume.VolumeDataSourceModel   `tfsdk:"volumes"`
}

func (data *VolumeFilterModel) parseVolumes(
	ctx context.Context,
	client *linodego.Client,
	volumes []linodego.Volume,
) diag.Diagnostics {
	result := make([]volume.VolumeDataSourceModel, len(volumes))
	for i := range volumes {
		var mod volume.VolumeDataSourceModel
		diags := mod.ParseComputedAttributes(ctx, &volumes[i])
		if diags.HasError() {
			return diags
		}
		diags = mod.ParseNonComputedAttributes(ctx, &volumes[i])
		if diags.HasError() {
			return diags
		}
		result[i] = mod
	}
	data.Volumes = result
	return nil
}
