package producerimagesharegroupmembers

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v3/linode/producerimagesharegroupmember"
)

type ImageShareGroupMemberFilterModel struct {
	ID           types.String                                    `tfsdk:"id"`
	ShareGroupID types.Int64                                     `tfsdk:"sharegroup_id"`
	Filters      frameworkfilter.FiltersModelType                `tfsdk:"filter"`
	Order        types.String                                    `tfsdk:"order"`
	OrderBy      types.String                                    `tfsdk:"order_by"`
	Members      []producerimagesharegroupmember.DataSourceModel `tfsdk:"members"`
}

func (model *ImageShareGroupMemberFilterModel) ParseImageShareGroupMembers(
	members []linodego.ImageShareGroupMember,
) {
	memberModels := make([]producerimagesharegroupmember.DataSourceModel, len(members))

	for i, member := range members {
		var memberModel producerimagesharegroupmember.DataSourceModel
		memberModel.ShareGroupID = model.ShareGroupID
		memberModel.ParseImageShareGroupMember(&member)
		memberModels[i] = memberModel

	}

	model.Members = memberModels
}
