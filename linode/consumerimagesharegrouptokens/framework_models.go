package consumerimagesharegrouptokens

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/consumerimagesharegrouptoken"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/frameworkfilter"
)

type ImageShareGroupTokenFilterModel struct {
	ID      types.String                                   `tfsdk:"id"`
	Filters frameworkfilter.FiltersModelType               `tfsdk:"filter"`
	Order   types.String                                   `tfsdk:"order"`
	OrderBy types.String                                   `tfsdk:"order_by"`
	Tokens  []consumerimagesharegrouptoken.DataSourceModel `tfsdk:"tokens"`
}

func (model *ImageShareGroupTokenFilterModel) ParseImageShareGroupTokens(
	tokens []linodego.ImageShareGroupToken,
) {
	tokenModels := make([]consumerimagesharegrouptoken.DataSourceModel, len(tokens))

	for i, token := range tokens {
		var tokenModel consumerimagesharegrouptoken.DataSourceModel
		tokenModel.ParseImageShareGroupToken(&token)
		tokenModels[i] = tokenModel

	}

	model.Tokens = tokenModels
}
