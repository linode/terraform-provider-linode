package placementgroups

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v3/linode/placementgroup"
)

type PlacementGroupFilterModel struct {
	ID              types.String                                   `tfsdk:"id"`
	Filters         frameworkfilter.FiltersModelType               `tfsdk:"filter"`
	Order           types.String                                   `tfsdk:"order"`
	OrderBy         types.String                                   `tfsdk:"order_by"`
	PlacementGroups []placementgroup.PlacementGroupDataSourceModel `tfsdk:"placement_groups"`
}

func (model *PlacementGroupFilterModel) ParsePlacementGroups(
	pgs []linodego.PlacementGroup,
) {
	pgModels := make([]placementgroup.PlacementGroupDataSourceModel, len(pgs))

	for i, pg := range pgs {
		pg := pg
		var pgModel placementgroup.PlacementGroupDataSourceModel
		pgModel.ID = types.Int64Value(int64(pg.ID))
		pgModel.ParsePlacementGroup(&pg)
		pgModels[i] = pgModel

	}

	model.PlacementGroups = pgModels
}
