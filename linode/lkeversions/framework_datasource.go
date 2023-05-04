package lkeversions

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	client *linodego.Client
}

func (data *DataSourceModel) parseVersions(ctx context.Context, lkeVersions []linodego.LKEVersion) diag.Diagnostics {

	versions, err := flattenVersions(lkeVersions)
	if err != nil {
		return err
	}

	data.Versions = *versions

	id, _ := json.Marshal(lkeVersions)

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

	meta := helper.GetDataSourceMeta(req, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	d.client = meta.Client
}

type DataSourceModel struct {
	Versions types.List   `tfsdk:"versions"`
	ID       types.String `tfsdk:"id"`
}

func (d *DataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "linode_lke_versions"
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
	client := d.client

	var data DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	versions, err := client.ListLKEVersions(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get LKE Versions: %s", err.Error(),
		)
		return
	}

	data.parseVersions(ctx, versions)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenVersions(versions []linodego.LKEVersion) (*basetypes.ListValue, diag.Diagnostics) {
	resultList := make([]attr.Value, len(versions))

	for i, field := range versions {
		valueMap := make(map[string]attr.Value)
		valueMap["id"] = types.StringValue(field.ID)

		obj, err := types.ObjectValue(lkeVersionObjectType.AttrTypes, valueMap)
		if err != nil {
			return nil, err
		}

		resultList[i] = obj
	}

	result, err := basetypes.NewListValue(
		lkeVersionObjectType,
		resultList,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
