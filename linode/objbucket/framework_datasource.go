package objbucket

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
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

func (data *DataSourceModel) parseObjectStorageBucket(bucket *linodego.ObjectStorageBucket) {
	data.Cluster = types.StringValue(bucket.Cluster)
	data.Region = types.StringValue(bucket.Region)
	data.Created = types.StringValue(bucket.Created.Format(time.RFC3339))
	data.Hostname = types.StringValue(bucket.Hostname)
	data.ID = types.StringValue(fmt.Sprintf("%s:%s", bucket.Cluster, bucket.Label))
	data.Label = types.StringValue(bucket.Label)
	data.Objects = types.Int64Value(int64(bucket.Objects))
	data.Size = types.Int64Value(int64(bucket.Size))
	data.EndpointType = types.StringValue(string(bucket.EndpointType))
	data.S3Endpoint = types.StringValue(bucket.S3Endpoint)
}

type DataSourceModel struct {
	Cluster      types.String `tfsdk:"cluster"`
	Region       types.String `tfsdk:"region"`
	Created      types.String `tfsdk:"created"`
	Hostname     types.String `tfsdk:"hostname"`
	ID           types.String `tfsdk:"id"`
	Label        types.String `tfsdk:"label"`
	Objects      types.Int64  `tfsdk:"objects"`
	Size         types.Int64  `tfsdk:"size"`
	EndpointType types.String `tfsdk:"endpoint_type"`
	S3Endpoint   types.String `tfsdk:"s3_endpoint"`
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
