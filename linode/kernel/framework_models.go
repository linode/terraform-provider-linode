package kernel

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper/customtypes"
)

type DataSourceModel struct {
	ID           types.String                       `tfsdk:"id"`
	Architecture types.String                       `tfsdk:"architecture"`
	Built        customtypes.RFC3339TimeStringValue `tfsdk:"built"`
	Deprecated   types.Bool                         `tfsdk:"deprecated"`
	KVM          types.Bool                         `tfsdk:"kvm"`
	Label        types.String                       `tfsdk:"label"`
	PVOPS        types.Bool                         `tfsdk:"pvops"`
	Version      types.String                       `tfsdk:"version"`
	XEN          types.Bool                         `tfsdk:"xen"`
}

func (data *DataSourceModel) ParseKernel(ctx context.Context, kernel *linodego.LinodeKernel) {
	data.ID = types.StringValue(kernel.ID)
	data.Architecture = types.StringValue(kernel.Architecture)

	built := types.StringNull()
	if kernel.Built != nil {
		built = types.StringValue(kernel.Built.Format(time.RFC3339))
	}

	data.Built = customtypes.RFC3339TimeStringValue{
		StringValue: built,
	}

	data.Deprecated = types.BoolValue(kernel.Deprecated)
	data.KVM = types.BoolValue(kernel.KVM)
	data.Label = types.StringValue(kernel.Label)
	data.PVOPS = types.BoolValue(kernel.PVOPS)
	data.Version = types.StringValue(kernel.Version)
	data.XEN = types.BoolValue(kernel.XEN)
}
