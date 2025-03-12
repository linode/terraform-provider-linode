package lkeversions

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_lke_versions",
				Schema: &frameworkDatasourceSchema,
			},
		),
	}
}

type DataSource struct {
	helper.BaseDataSource
}

func (data *DataSourceModel) parseVersions(ctx context.Context, lkeVersions []linodego.LKEVersion) diag.Diagnostics {
	versions, diag := flattenVersions(lkeVersions)
	if diag.HasError() {
		return diag
	}

	data.Versions = *versions

	id, _ := json.Marshal(lkeVersions)

	data.ID = types.StringValue(string(id))

	return nil
}

type DataSourceModel struct {
	Versions types.List   `tfsdk:"versions"`
	ID       types.String `tfsdk:"id"`
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data.linode_lke_versions")

	client := d.Meta.Client

	var data DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "client.ListLKEVersions(...)")
	versions, err := client.ListLKEVersions(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get LKE Versions: %s", err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.parseVersions(ctx, versions)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenVersions(versions []linodego.LKEVersion) (*basetypes.ListValue, diag.Diagnostics) {
	resultList := make([]attr.Value, len(versions))

	for i, field := range versions {
		valueMap := make(map[string]attr.Value)
		valueMap["id"] = types.StringValue(field.ID)

		obj, diag := types.ObjectValue(lkeVersionObjectType.AttrTypes, valueMap)
		if diag.HasError() {
			return nil, diag
		}

		resultList[i] = obj
	}

	result, diag := basetypes.NewListValue(
		lkeVersionObjectType,
		resultList,
	)
	if diag.HasError() {
		return nil, diag
	}
	return &result, nil
}
