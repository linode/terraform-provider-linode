package images

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/linode/image"
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
	images []linodego.Image,
) {
	result := make([]image.ImageModel, len(images))
	for i := range images {
		var imgData image.ImageModel
		imgData.ParseImage(&images[i])
		result[i] = imgData
	}

	data.Images = result
}
