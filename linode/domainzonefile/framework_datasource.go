package domainzonefile

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	meta *helper.FrameworkProviderMeta
}

func (data *DataSourceModel) parseDomainZoneFile(
	ctx context.Context, zone *linodego.DomainZoneFile,
) diag.Diagnostics {
	file, diags := types.ListValueFrom(ctx, types.StringType, zone.ZoneFile)
	if diags.HasError() {
		return diags
	}
	data.ZoneFile = file

	id, err := json.Marshal(zone)
	if err != nil {
		diags.AddError("Error marshalling json: %s", err.Error())
		return diags
	}

	data.ID = types.StringValue(string(id))

	return nil
}

func (d *DataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	d.meta = helper.GetDataSourceMeta(req, resp)
	if resp.Diagnostics.HasError() {
		return
	}
}

type DataSourceModel struct {
	DomainID types.Int64  `tfsdk:"domain_id"`
	ZoneFile types.List   `tfsdk:"zone_file"`
	ID       types.String `tfsdk:"id"`
}

func (d *DataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "linode_domain_zonefile"
}

func (d *DataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = frameworkDatasourceSchema
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	client := d.meta.Client

	var data DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	retry := time.Duration(d.meta.Config.EventPollMilliseconds.ValueInt64())

	zf, diags := getZoneFileRetry(ctx, client, int(data.DomainID.ValueInt64()), retry)
	if diags.HasError() {
		return
	}

	resp.Diagnostics.Append(data.parseDomainZoneFile(ctx, zf)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func getZoneFileRetry(ctx context.Context, client *linodego.Client,
	domainID int, retryDuration time.Duration,
) (*linodego.DomainZoneFile, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	ticker := time.NewTicker(retryDuration)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			zf, err := client.GetDomainZoneFile(ctx, domainID)
			if err != nil {
				diags.AddError(
					"Failed to fetch domain record: %s", err.Error(),
				)
				return nil, diags
			}
			if len(zf.ZoneFile) > 0 {
				return zf, nil
			}
		case <-ctx.Done():
			diags.AddError(
				"Unable to fetch domain record", "",
			)
			return nil, diags
		}
	}
}
