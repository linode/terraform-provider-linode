package images

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v2/linode/image"
)

// ImageFilterModel describes the Terraform resource data model to match the
// resource schema.
type ImageFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Latest  types.Bool                       `tfsdk:"latest"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order   types.String                     `tfsdk:"order"`
	OrderBy types.String                     `tfsdk:"order_by"`
	Images  []image.ImageModel               `tfsdk:"images"`
}

func (data *ImageFilterModel) parseImages(
	ctx context.Context,
	images []linodego.Image,
) diag.Diagnostics {
	result := make([]image.ImageModel, len(images))
	for i := range images {
		var imgData image.ImageModel
		diags := imgData.ParseImage(ctx, &images[i])
		if diags.HasError() {
			return diags
		}
		result[i] = imgData
	}

	data.Images = result

	return nil
}
