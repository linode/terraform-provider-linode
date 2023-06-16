package objbucket

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	client *linodego.Client
}

func (data *DataSourceModel) parseObjectStorageBucket(bucket *linodego.ObjectStorageBucket) {
	data.Cluster = types.StringValue(bucket.Cluster)
	data.Created = types.StringValue(bucket.Created.Format(time.RFC3339))
	data.Hostname = types.StringValue(bucket.Hostname)
	data.ID = types.StringValue(fmt.Sprintf("%s:%s", bucket.Cluster, bucket.Label))
	data.Label = types.StringValue(bucket.Label)
	data.Objects = types.Int64Value(int64(bucket.Objects))
	data.Size = types.Int64Value(int64(bucket.Size))
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
	Cluster  types.String `tfsdk:"cluster"`
	Created  types.String `tfsdk:"created"`
	Hostname types.String `tfsdk:"hostname"`
	ID       types.String `tfsdk:"id"`
	Label    types.String `tfsdk:"label"`
	Objects  types.Int64  `tfsdk:"objects"`
	Size     types.Int64  `tfsdk:"size"`
}

func (d *DataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "linode_object_storage_bucket"
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

	bucket, err := client.GetObjectStorageBucket(ctx, data.Cluster.ValueString(), data.Label.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to find the specified Linode ObjectStorageBucket: %s", err.Error(),
		)
		return
	}

	data.parseObjectStorageBucket(bucket)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
