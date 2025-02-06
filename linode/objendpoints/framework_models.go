package objendpoints

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

type ObjectStorageEndpointFilterModel struct {
	ID        types.String                     `tfsdk:"id"`
	Filters   frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order     types.String                     `tfsdk:"order"`
	OrderBy   types.String                     `tfsdk:"order_by"`
	Endpoints []ObjectStorageEndpointModel     `tfsdk:"endpoints"`
}

type ObjectStorageEndpointModel struct {
	EndpointType types.String `tfsdk:"endpoint_type"`
	Region       types.String `tfsdk:"region"`
	S3Endpoint   types.String `tfsdk:"s3_endpoint"`
}

func (data *ObjectStorageEndpointModel) parseObjectStorageEndpoint(
	endpoint linodego.ObjectStorageEndpoint,
) {
	data.EndpointType = types.StringValue(string(endpoint.EndpointType))
	data.Region = types.StringValue(endpoint.Region)
	data.S3Endpoint = types.StringPointerValue(endpoint.S3Endpoint)
}

func (model *ObjectStorageEndpointFilterModel) parseObjectStorageEndpoints(
	endpoints []linodego.ObjectStorageEndpoint,
) {
	result := make([]ObjectStorageEndpointModel, len(endpoints))
	for i, e := range endpoints {
		result[i].parseObjectStorageEndpoint(e)
	}

	model.Endpoints = result
}
