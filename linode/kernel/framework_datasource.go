package kernel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_kernel",
				Schema: &frameworkDatasourceSchema,
			},
		),
	}
}

type DataSource struct {
	helper.BaseDataSource
}

func (data *DataSourceModel) parseKernel(kernel *linodego.LinodeKernel) {
	data.ID = types.StringValue(kernel.ID)
	data.Architecture = types.StringValue(kernel.Architecture)
	data.Deprecated = types.BoolValue(kernel.Deprecated)
	data.KVM = types.BoolValue(kernel.KVM)
	data.Label = types.StringValue(kernel.Label)
	data.PVOPS = types.BoolValue(kernel.PVOPS)
	data.Version = types.StringValue(kernel.Version)
	data.XEN = types.BoolValue(kernel.XEN)
}

type DataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Architecture types.String `tfsdk:"architecture"`
	Deprecated   types.Bool   `tfsdk:"deprecated"`
	KVM          types.Bool   `tfsdk:"kvm"`
	Label        types.String `tfsdk:"label"`
	PVOPS        types.Bool   `tfsdk:"pvops"`
	Version      types.String `tfsdk:"version"`
	XEN          types.Bool   `tfsdk:"xen"`
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	client := d.Meta.Client

	var data DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	kernel, err := client.GetKernel(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get Kernel: %s", err.Error(),
		)
		return
	}

	data.parseKernel(kernel)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
