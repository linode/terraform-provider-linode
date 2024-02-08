package stackscript

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_stackscript",
				Schema: &frameworkDataSourceSchema,
			},
		),
	}
}

type DataSource struct {
	helper.BaseDataSource
}

func (r *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	client := r.Meta.Client

	var data StackScriptModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := helper.StringToInt(data.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "stackscript_id", id)

	tflog.Debug(ctx, "Read data.linode_stackscript")

	tflog.Trace(ctx, "client.GetStackscript(...)")

	stackscript, err := client.GetStackscript(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find the Linode StackScript",
			fmt.Sprintf(
				"Error finding the specified Linode StackScript: %s",
				err.Error(),
			),
		)
		return
	}

	resp.Diagnostics.Append(data.ParseComputedAttributes(ctx, stackscript)...)
	resp.Diagnostics.Append(data.ParseNonComputedAttributes(ctx, stackscript)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
