package lkeversions

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/lkeversion"
)

type DataSourceModel struct {
	Versions []lkeversion.DataSourceModel `tfsdk:"versions"`
	ID       types.String                 `tfsdk:"id"`
	Tier     types.String                 `tfsdk:"tier"`
}

func (model *DataSourceModel) parseLKEVersions(lkeVersions []linodego.LKEVersion,
) diag.Diagnostics {
	result := make([]lkeversion.DataSourceModel, len(lkeVersions))

	for i := range lkeVersions {
		var m lkeversion.DataSourceModel

		diags := m.ParseLKEVersion(&lkeVersions[i])
		if diags.HasError() {
			return diags
		}

		result[i] = m
	}

	model.Versions = result

	id, _ := json.Marshal(lkeVersions)
	model.ID = types.StringValue(string(id))

	return nil
}

func (model *DataSourceModel) parseLKETierVersions(lkeTierVersions []linodego.LKETierVersion,
) diag.Diagnostics {
	result := make([]lkeversion.DataSourceModel, len(lkeTierVersions))

	for i := range lkeTierVersions {
		var m lkeversion.DataSourceModel

		diags := m.ParseLKETierVersion(&lkeTierVersions[i])
		if diags.HasError() {
			return diags
		}

		result[i] = m
	}

	model.Versions = result

	id, _ := json.Marshal(lkeTierVersions)
	model.ID = types.StringValue(string(id))

	return nil
}
