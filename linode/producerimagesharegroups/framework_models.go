package producerimagesharegroups

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v3/linode/producerimagesharegroup"
)

type ImageShareGroupFilterModel struct {
	ID               types.String                              `tfsdk:"id"`
	Filters          frameworkfilter.FiltersModelType          `tfsdk:"filter"`
	Order            types.String                              `tfsdk:"order"`
	OrderBy          types.String                              `tfsdk:"order_by"`
	ImageShareGroups []producerimagesharegroup.DataSourceModel `tfsdk:"image_share_groups"`
}

func (model *ImageShareGroupFilterModel) ParseImageShareGroups(
	sgs []linodego.ProducerImageShareGroup,
) {
	sgModels := make([]producerimagesharegroup.DataSourceModel, len(sgs))

	for i, sg := range sgs {
		var sgModel producerimagesharegroup.DataSourceModel
		sgModel.ID = types.Int64Value(int64(sg.ID))
		sgModel.ParseImageShareGroup(&sg)
		sgModels[i] = sgModel

	}

	model.ImageShareGroups = sgModels
}
