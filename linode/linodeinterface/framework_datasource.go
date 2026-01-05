package linodeinterface

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type DataSource struct {
	helper.BaseDataSource
}

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_interface",
				Schema: &frameworkDataSourceSchema,
			},
		),
	}
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data."+d.Config.Name)

	var data DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Both id and linode_id are required
	if data.ID.IsNull() || data.LinodeID.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Required Configuration",
			"Both id and linode_id must be specified",
		)
		return
	}

	// Convert string ID to int for API call
	interfaceID, parseErr := strconv.Atoi(data.ID.ValueString())
	if parseErr != nil {
		resp.Diagnostics.AddError(
			"Invalid Interface ID",
			fmt.Sprintf("Could not parse interface ID: %s", parseErr),
		)
		return
	}

	linodeID := helper.FrameworkSafeInt64ToInt(data.LinodeID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "interface_id", interfaceID)
	ctx = tflog.SetField(ctx, "linode_id", linodeID)

	linodeInterface, err := d.getInterfaceByID(ctx, linodeID, interfaceID)
	if err != nil {
		resp.Diagnostics.Append(err)
		return
	}

	// Flatten the interface into the data source model
	data.FlattenInterface(ctx, *linodeInterface, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *DataSource) getInterfaceByID(ctx context.Context, linodeID, interfaceID int) (*linodego.LinodeInterface, diag.Diagnostic) {
	tflog.Trace(ctx, "client.GetInterface(...)")
	linodeInterface, err := d.Meta.Client.GetInterface(ctx, linodeID, interfaceID)
	if err != nil {
		return nil, diag.NewErrorDiagnostic(
			fmt.Sprintf("Failed to get Interface %d for Linode %d", interfaceID, linodeID),
			err.Error(),
		)
	}

	return linodeInterface, nil
}
