package images

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/image"
)

var filterableFields = []string{"deprecated", "is_public", "label", "size", "vendor"}

var dataSourceSchema = map[string]*schema.Schema{
	"latest": {
		Type:        schema.TypeBool,
		Description: "If true, only the latest image will be returned.",
		Optional:    true,
		Default:     false,
	},
	"order_by": helper.OrderBySchema(filterableFields),
	"order":    helper.OrderSchema(),
	"filter":   helper.FilterSchema(filterableFields),
	"images": {
		Type:        schema.TypeList,
		Description: "The returned list of Images.",
		Computed:    true,
		Elem:        image.DataSource(),
	},
}
