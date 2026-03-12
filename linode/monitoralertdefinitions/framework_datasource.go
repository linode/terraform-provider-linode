package monitoralertdefinitions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type DataSource struct {
	helper.BaseDataSource
}

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(helper.BaseDataSourceConfig{
			Name:   "linode_monitor_alert_definitions",
			Schema: &frameworkDatasourceSchema,
		}),
	}
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data."+d.Config.Name)

	var data AlertDefinitionFilterModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diag := filterConfig.GenerateID(data.Filters)
	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}
	data.ID = id

	var result []any
	if !data.ServiceType.IsUnknown() && !data.ServiceType.IsNull() {
		result, diag = filterConfig.GetAndFilter(
			ctx, d.Meta.Client, data.Filters, data.listAlertDefinitionsByServiceType,
			types.StringNull(), types.StringNull())
		if diag != nil {
			resp.Diagnostics.Append(diag)
			return
		}
	} else {
		result, diag = filterConfig.GetAndFilter(
			ctx, d.Meta.Client, data.Filters, listAllAlertDefinitions,
			types.StringNull(), types.StringNull())
		if diag != nil {
			resp.Diagnostics.Append(diag)
			return
		}
	}

	resp.Diagnostics.Append(data.parseAlertDefinitions(ctx, helper.AnySliceToTyped[linodego.AlertDefinition](result))...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func listAllAlertDefinitions(
	ctx context.Context,
	client *linodego.Client,
	_ string,
) ([]any, error) {
	tflog.Trace(ctx, "client.ListAllMonitorAlertDefinitions(...)")

	alertDefinitions, err := client.ListAllMonitorAlertDefinitions(ctx, nil)
	if err != nil {
		return nil, err
	}

	return helper.TypedSliceToAny(alertDefinitions), nil
}

func (data *AlertDefinitionFilterModel) listAlertDefinitionsByServiceType(
	ctx context.Context,
	client *linodego.Client,
	_ string,
) ([]any, error) {
	tflog.Trace(ctx, "client.ListMonitorAlertDefinitions(...)")

	alertDefinitions, err := client.ListMonitorAlertDefinitions(
		ctx,
		data.ServiceType.ValueString(),
		nil)
	if err != nil {
		return nil, err
	}

	return helper.TypedSliceToAny(alertDefinitions), nil
}
