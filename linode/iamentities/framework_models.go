package iamentities

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

type ListEntitiesModel struct {
	ID       types.String                     `tfsdk:"id"`
	Filters  frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order    types.String                     `tfsdk:"order"`
	OrderBy  types.String                     `tfsdk:"order_by"`
	Entities []Entity                         `tfsdk:"entities"`
}

type Entity struct {
	ID    types.Int64  `tfsdk:"id"`
	Label types.String `tfsdk:"label"`
	Type  types.String `tfsdk:"type"`
}

func (data *ListEntitiesModel) parseEntities(
	entities []linodego.LinodeEntity,
) {
	results := make([]Entity, len(entities))

	for i, r := range entities {
		results[i].ID = types.Int64Value(int64(r.ID))
		results[i].Label = types.StringValue(r.Label)
		results[i].Type = types.StringValue(r.Type)
	}

	data.Entities = results
}
