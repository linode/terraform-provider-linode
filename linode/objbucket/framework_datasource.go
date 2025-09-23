package objbucket

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_object_storage_bucket",
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
	tflog.Debug(ctx, "Read data."+d.Config.Name)
	client := d.Meta.Client

	var data DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cluster := data.Cluster.ValueString()
	region := data.Region.ValueString()

	var regionOrCluster string
	if region != "" {
		regionOrCluster = region
	} else {
		regionOrCluster = cluster
	}
	bucketLabel := data.Label.ValueString()

	ctx = helper.SetLogFieldBulk(ctx, map[string]any{
		"region_or_cluster": regionOrCluster,
		"bucket":            bucketLabel,
	})

	bucket, err := client.GetObjectStorageBucket(ctx, regionOrCluster, bucketLabel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to find the specified Linode ObjectStorageBucket: %s", err.Error(),
		)
		return
	}

	data.parseObjectStorageBucket(bucket)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
