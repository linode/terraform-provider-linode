package domainrecord

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_domain_record",
				Schema: &frameworkDatasourceSchema,
			},
		),
	}
}

type DataSource struct {
	helper.BaseDataSource
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data.linode_domain_record")

	client := d.Meta.Client

	var data DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Name.ValueString() == "" && (data.ID.IsNull() || data.ID.IsUnknown()) {
		resp.Diagnostics.AddError("Record name or ID is required", "")
		return
	}

	var record *linodego.DomainRecord

	if !(data.ID.IsNull() || data.ID.IsUnknown()) {
		domainID := helper.FrameworkSafeInt64ToInt(
			data.DomainID.ValueInt64(),
			&resp.Diagnostics,
		)
		recordID := helper.FrameworkSafeInt64ToInt(
			data.ID.ValueInt64(),
			&resp.Diagnostics,
		)
		if resp.Diagnostics.HasError() {
			return
		}

		ctx = tflog.SetField(ctx, "domain_id", domainID)
		ctx = tflog.SetField(ctx, "record_id", recordID)

		tflog.Trace(ctx, "client.GetDomainRecord(...)")

		rec, err := client.GetDomainRecord(ctx, domainID, recordID)
		if err != nil {
			resp.Diagnostics.AddError("Error fetching domain record: %v", err.Error())
			return
		}
		record = rec
	} else if data.Name.ValueString() != "" {
		filter, _ := json.Marshal(map[string]interface{}{"name": data.Name.ValueString()})
		domainID := helper.FrameworkSafeInt64ToInt(
			data.DomainID.ValueInt64(),
			&resp.Diagnostics,
		)
		if resp.Diagnostics.HasError() {
			return
		}

		ctx = tflog.SetField(ctx, "domain_id", domainID)
		ctx = tflog.SetField(ctx, "record_name", data.Name.ValueString())

		tflog.Trace(ctx, "client.ListDomainRecords(...)")

		records, err := client.ListDomainRecords(ctx, domainID,
			linodego.NewListOptions(0, string(filter)))
		if err != nil {
			resp.Diagnostics.AddError("Error listing domain records: %v", err.Error())
			return
		}
		if len(records) > 0 {
			record = &records[0]
		}
	}

	if record != nil {
		data.FlattenDomainRecord(record)
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	} else {
		resp.Diagnostics.AddError(fmt.Sprintf(`Domain record "%s" for domain %s was not found`,
			data.Name.ValueString(), data.DomainID.String()), "")
	}
}
