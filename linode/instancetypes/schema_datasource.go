package instancetypes

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/instancetype"
)

var dataSourceSchema = map[string]*schema.Schema{
	"filter": helper.FilterSchema([]string{"class", "disk", "gpus", "label",
		"memory", "network_out", "transfer", "vcpus"}),
	"types": {
		Type:        schema.TypeList,
		Description: "The returned list of Types.",
		Computed:    true,
		Elem:        instancetype.DataSource(),
	},
}
