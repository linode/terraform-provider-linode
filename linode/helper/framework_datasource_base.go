package helper

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// NewBaseDataSource returns a new instance of the BaseDataSource
// struct for cleaner initialization.
func NewBaseDataSource(name string, schemaObject schema.Schema) BaseDataSource {
	return BaseDataSource{
		TypeName:     name,
		SchemaObject: schemaObject,
	}
}

// BaseDataSource contains various re-usable fields and methods
// intended for use in data source implementations by composition.
type BaseDataSource struct {
	Meta *FrameworkProviderMeta

	SchemaObject schema.Schema
	TypeName     string
}

func (r *BaseDataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	r.Meta = GetDataSourceMeta(req, resp)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *BaseDataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = r.TypeName
}

func (r *BaseDataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = r.SchemaObject
}
