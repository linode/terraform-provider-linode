package domainrecord

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
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

func (data *DataSourceModel) parseDomainRecord(domainRecord *linodego.DomainRecord) {
	data.ID = types.Int64Value(int64(domainRecord.ID))
	data.Name = types.StringValue(domainRecord.Name)
	data.Type = types.StringValue(string(domainRecord.Type))
	data.TTLSec = types.Int64Value(int64(domainRecord.TTLSec))
	data.Target = types.StringValue(domainRecord.Target)
	data.Priority = types.Int64Value(int64(domainRecord.Priority))
	data.Weight = types.Int64Value(int64(domainRecord.Weight))
	data.Port = types.Int64Value(int64(domainRecord.Port))
	data.Protocol = types.StringPointerValue(domainRecord.Protocol)
	data.Service = types.StringPointerValue(domainRecord.Service)
	data.Tag = types.StringPointerValue(domainRecord.Tag)
}

type DataSourceModel struct {
	ID       types.Int64  `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	DomainID types.Int64  `tfsdk:"domain_id"`
	Type     types.String `tfsdk:"type"`
	TTLSec   types.Int64  `tfsdk:"ttl_sec"`
	Target   types.String `tfsdk:"target"`
	Priority types.Int64  `tfsdk:"priority"`
	Weight   types.Int64  `tfsdk:"weight"`
	Port     types.Int64  `tfsdk:"port"`
	Protocol types.String `tfsdk:"protocol"`
	Service  types.String `tfsdk:"service"`
	Tag      types.String `tfsdk:"tag"`
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
		data.parseDomainRecord(record)
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	} else {
		resp.Diagnostics.AddError(fmt.Sprintf(`Domain record "%s" for domain %s was not found`,
			data.Name.ValueString(), data.DomainID.String()), "")
	}
}
