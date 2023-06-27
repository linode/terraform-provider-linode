package helper

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NewBaseResource returns a new instance of the BaseResource
// struct for cleaner initialization.
func NewBaseResource(cfg BaseResourceConfig) BaseResource {
	return BaseResource{
		Config: cfg,
	}
}

// BaseResourceConfig contains all configurable base resource fields.
type BaseResourceConfig struct {
	Name   string
	IDAttr string
	IDType attr.Type

	// Optional
	Schema *schema.Schema
}

// BaseResource contains various re-usable fields and methods
// intended for use in resource implementations by composition.
type BaseResource struct {
	Config BaseResourceConfig
	Meta   *FrameworkProviderMeta
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
	resp.TypeName = r.Config.Name
}

func (r *BaseResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	if r.Config.Schema == nil {
		resp.Diagnostics.AddError(
			"Missing Schema",
			"Base resource was not provided a schema. "+
				"Please provide a Schema config attribute or implement, the Schema(...) function.",
		)
		return
	}

	resp.Schema = *r.Config.Schema
}

// ImportState should be overridden for resources with
// complex read logic (e.g. parent ID).
func (r *BaseResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Enforce defaults
	idAttr := r.Config.IDAttr
	if idAttr == "" {
		idAttr = "id"
	}

	idType := r.Config.IDType
	if idType == nil {
		idType = types.Int64Type
	}

	attrPath := path.Root("id")

	if attrPath.Equal(path.Empty()) {
		resp.Diagnostics.AddError(
			"Resource Import Passthrough Missing Attribute Path",
			"This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				"Resource ImportState path must be set to a valid attribute path.",
		)
		return
	}

	// Handle type conversion
	var err error
	var idValue any

	switch idType {
	case types.Int64Type:
		idValue, err = strconv.ParseInt(req.ID, 10, 64)
	case types.StringType:
		idValue = req.ID
	default:
		err = fmt.Errorf("unsupported id attribute type: %v", idType)
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to convert ID attribute",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, attrPath, idValue)...)
}
