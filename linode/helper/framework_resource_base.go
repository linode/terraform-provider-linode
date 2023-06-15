package helper

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// NewBaseResource returns a new instance of the BaseResource
// struct for cleaner initialization.
func NewBaseResource(name string, schemaObject schema.Schema) BaseResource {
	return BaseResource{
		TypeName:     name,
		SchemaObject: schemaObject,
	}
}

// BaseResource contains various re-usable fields and methods
// intended for use in resource implementations by composition.
type BaseResource struct {
	Meta *FrameworkProviderMeta

	SchemaObject schema.Schema
	TypeName     string
}

func (r *BaseResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	r.Meta = GetResourceMeta(req, resp)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *BaseResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = r.TypeName
}

func (r *BaseResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = r.SchemaObject
}
