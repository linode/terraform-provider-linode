package image

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

// ImageModel describes the Terraform resource data model to match the
// resource schema.
type ImageModel struct {
	ID           types.String   `tfsdk:"id"`
	Label        types.String   `tfsdk:"label"`
	Description  types.String   `tfsdk:"description"`
	Capabilities []types.String `tfsdk:"capabilities"`
	Created      types.String   `tfsdk:"created"`
	CreatedBy    types.String   `tfsdk:"created_by"`
	Deprecated   types.Bool     `tfsdk:"deprecated"`
	IsPublic     types.Bool     `tfsdk:"is_public"`
	Size         types.Int64    `tfsdk:"size"`
	Status       types.String   `tfsdk:"status"`
	Type         types.String   `tfsdk:"type"`
	Expiry       types.String   `tfsdk:"expiry"`
	Vendor       types.String   `tfsdk:"vendor"`
}

func (data *ImageModel) ParseImage(
	image *linodego.Image,
) {
	data.ID = types.StringValue(image.ID)
	data.Label = types.StringValue(image.Label)

	data.Description = types.StringValue(image.Description)
	if image.Created != nil {
		data.Created = types.StringValue(image.Created.Format(time.RFC3339))
	} else {
		data.Created = types.StringNull()
	}
	if image.Expiry != nil {
		data.Expiry = types.StringValue(image.Expiry.Format(time.RFC3339))
	} else {
		data.Expiry = types.StringNull()
	}
	data.Capabilities = helper.StringSliceToFramework(image.Capabilities)
	data.CreatedBy = types.StringValue(image.CreatedBy)
	data.Deprecated = types.BoolValue(image.Deprecated)
	data.IsPublic = types.BoolValue(image.IsPublic)
	data.Size = types.Int64Value(int64(image.Size))
	data.Status = types.StringValue(string(image.Status))
	data.Type = types.StringValue(image.Type)
	data.Vendor = types.StringValue(image.Vendor)
}
