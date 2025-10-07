package objbucket

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
)

type BaseModel struct {
	ID           types.String      `tfsdk:"id"`
	Label        types.String      `tfsdk:"label"`
	Cluster      types.String      `tfsdk:"cluster"`
	Region       types.String      `tfsdk:"region"`
	EndpointType types.String      `tfsdk:"endpoint_type"`
	S3Endpoint   types.String      `tfsdk:"s3_endpoint"`
	Hostname     types.String      `tfsdk:"hostname"`
	Objects      types.Int64       `tfsdk:"objects"`
	Size         types.Int64       `tfsdk:"size"`
	Created      timetypes.RFC3339 `tfsdk:"created"`
}

type DataSourceModel struct {
	BaseModel
}

func (data *DataSourceModel) parseObjectStorageBucket(bucket *linodego.ObjectStorageBucket) {
	data.Cluster = types.StringValue(bucket.Cluster)
	data.Region = types.StringValue(bucket.Region)
	data.Created = timetypes.NewRFC3339TimePointerValue(bucket.Created)
	data.Hostname = types.StringValue(bucket.Hostname)
	data.ID = types.StringValue(fmt.Sprintf("%s:%s", bucket.Cluster, bucket.Label))
	data.Label = types.StringValue(bucket.Label)
	data.Objects = types.Int64Value(int64(bucket.Objects))
	data.Size = types.Int64Value(int64(bucket.Size))
	data.EndpointType = types.StringValue(string(bucket.EndpointType))
	data.S3Endpoint = types.StringValue(bucket.S3Endpoint)
}
