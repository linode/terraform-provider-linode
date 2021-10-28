package instancetypes

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/instancetype"
)

var filterableFields = []string{"class", "disk", "gpus", "label",
	"memory", "network_out", "transfer", "vcpus"}

var dataSourceSchema = map[string]*schema.Schema{
	"order_by": helper.OrderBySchema(filterableFields),
	"order":    helper.OrderSchema(),
	"filter":   helper.FilterSchema(filterableFields),
	"types": {
		Type:        schema.TypeList,
		Description: "The returned list of Types.",
		Computed:    true,
		Elem:        instancetype.DataSource(),
	},
}
