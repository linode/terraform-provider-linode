package lkeversion

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
)

type DataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Tier types.String `tfsdk:"tier"`
}

func (data *DataSourceModel) ParseLKEVersion(lkeVersion *linodego.LKEVersion,
) diag.Diagnostics {
	data.ID = types.StringValue(lkeVersion.ID)

	return nil
}

func (data *DataSourceModel) ParseLKETierVersion(lkeTierVersion *linodego.LKETierVersion,
) diag.Diagnostics {
	data.ID = types.StringValue(lkeTierVersion.ID)
	data.Tier = types.StringValue(string(lkeTierVersion.Tier))

	return nil
}
