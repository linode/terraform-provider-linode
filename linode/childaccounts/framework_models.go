package childaccounts

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/account"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

type ChildAccountFilterModel struct {
	ID            types.String                     `tfsdk:"id"`
	Filters       frameworkfilter.FiltersModelType `tfsdk:"filter"`
	ChildAccounts []account.DataSourceModel        `tfsdk:"child_accounts"`
}

func (model *ChildAccountFilterModel) parseAccounts(accounts []linodego.ChildAccount) {
	result := make([]account.DataSourceModel, len(accounts))

	for i, childAccount := range accounts {
		result[i].ParseAccount(&childAccount)
	}

	model.ChildAccounts = result
}
