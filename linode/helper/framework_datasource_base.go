package helper

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/linodego"
)

// NewBaseDataSource returns a new instance of the BaseDataSource
// struct for cleaner initialization.
func NewBaseDataSource(cfg BaseDataSourceConfig) BaseDataSource {
	return BaseDataSource{
		Config: cfg,
	}
}

// BaseDataSourceConfig contains all configurable base resource fields.
type BaseDataSourceConfig struct {
	Name string

	// Optional
	Schema        *schema.Schema
	IsEarlyAccess bool
}

// BaseDataSource contains various re-usable fields and methods
// intended for use in data source implementations by composition.
type BaseDataSource struct {
	Config BaseDataSourceConfig
	Meta   *FrameworkProviderMeta

	// Data source level linodego Client containing specific configurations.
	Client *linodego.Client
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

	r.Client = GetFwClientWithUserAgent("data."+r.Config.Name, r.Meta)

	if r.Config.IsEarlyAccess {
		resp.Diagnostics.Append(
			AttemptWarnEarlyAccessFramework(r.Meta.Config)...,
		)
	}
}

func (r *BaseDataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = r.Config.Name
}

func (r *BaseDataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	if r.Config.Schema == nil {
		resp.Diagnostics.AddError(
			"Missing Schema",
			"Base data source was not provided a schema. "+
				"Please provide a Schema config attribute or implement, the Schema(...) function.",
		)
		return
	}

	resp.Schema = *r.Config.Schema
}
