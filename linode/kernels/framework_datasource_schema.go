package kernels

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v2/linode/kernel"
)

var filterConfig = frameworkfilter.Config{
	"id":           {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeString},
	"architecture": {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"deprecated":   {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeBool},
	"kvm":          {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeBool},
	"label":        {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"pvops":        {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeBool},
	"version":      {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
	"xen":          {APIFilterable: false, TypeFunc: frameworkfilter.FilterTypeBool},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"order":    filterConfig.OrderSchema(),
		"order_by": filterConfig.OrderBySchema(),
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
		"kernels": schema.ListNestedBlock{
			Description: "The returned list of Kernels.",
			NestedObject: schema.NestedBlockObject{
				Attributes: kernel.KernelAttributes,
			},
		},
	},
}
