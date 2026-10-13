package kernel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
)

type DataSourceModel struct {
	ID           types.String      `tfsdk:"id"`
	Architecture types.String      `tfsdk:"architecture"`
	Built        timetypes.RFC3339 `tfsdk:"built"`
	Deprecated   types.Bool        `tfsdk:"deprecated"`
	KVM          types.Bool        `tfsdk:"kvm"`
	Label        types.String      `tfsdk:"label"`
	PVOPS        types.Bool        `tfsdk:"pvops"`
	Version      types.String      `tfsdk:"version"`
}

func (data *DataSourceModel) ParseKernel(_ context.Context, kernel *linodego.LinodeKernel) {
	data.ID = types.StringValue(kernel.ID)
	data.Architecture = types.StringValue(kernel.Architecture)
	data.Built = timetypes.NewRFC3339TimePointerValue(kernel.Built)
	data.Deprecated = types.BoolValue(kernel.Deprecated)
	data.KVM = types.BoolValue(kernel.KVM)
	data.Label = types.StringValue(kernel.Label)
	data.PVOPS = types.BoolValue(kernel.PVOPS)
	data.Version = types.StringValue(kernel.Version)
}
