package prefixlists

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v3/linode/prefixlist"
)

type PrefixListFilterModel struct {
	ID          types.String                      `tfsdk:"id"`
	Filters     frameworkfilter.FiltersModelType   `tfsdk:"filter"`
	PrefixLists []prefixlist.PrefixListBaseModel   `tfsdk:"prefix_lists"`
}

func (data *PrefixListFilterModel) parsePrefixLists(
	ctx context.Context,
	lists []linodego.PrefixList,
	diags *diag.Diagnostics,
) {
	data.PrefixLists = make([]prefixlist.PrefixListBaseModel, len(lists))
	for i, pl := range lists {
		data.PrefixLists[i].FlattenPrefixList(ctx, pl, diags, false)
	}
}
