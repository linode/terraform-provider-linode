package monitoralertdefinition

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type DataSource struct {
	helper.BaseDataSource
}

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_monitor_alert_definition",
				Schema: &FrameworkDataSourceSchema,
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

	var data AlertDefinitionDataSourceModel
	client := d.Meta.Client

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueInt64()

	ctx = tflog.SetField(ctx, "alert_definition_id", id)

	alertDefinition, err := client.GetMonitorAlertDefinition(ctx, data.ServiceType.ValueString(), int(id))
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get Alert Definition %d for service type %s", id, data.ServiceType.ValueString()),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.FlattenDataSourceModel(ctx, alertDefinition, false)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
