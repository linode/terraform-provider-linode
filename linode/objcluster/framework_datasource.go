package objcluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_object_storage_cluster",
				Schema: &frameworkDatasourceSchema,
			},
		),
	}
}

type DataSource struct {
	helper.BaseDataSource
}

func (data *DataSourceModel) parseCluster(cluster *linodego.ObjectStorageCluster) {
	data.ID = types.StringValue(cluster.ID)
	data.Domain = types.StringValue(cluster.Domain)
	data.Status = types.StringValue(cluster.Status)
	data.Region = types.StringValue(cluster.Region)
	data.StaticSiteDomain = types.StringValue(cluster.StaticSiteDomain)
}

type DataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	Domain           types.String `tfsdk:"domain"`
	Status           types.String `tfsdk:"status"`
	Region           types.String `tfsdk:"region"`
	StaticSiteDomain types.String `tfsdk:"static_site_domain"`
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

	cluster, err := client.GetObjectStorageCluster(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get LKE Versions: %s", err.Error(),
		)
		return
	}

	data.parseCluster(cluster)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
