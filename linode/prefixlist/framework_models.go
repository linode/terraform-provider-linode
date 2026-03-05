package prefixlist

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

// PrefixListBaseModel contains the shared fields for both single and list data sources.
type PrefixListBaseModel struct {
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	Visibility         types.String `tfsdk:"visibility"`
	SourcePrefixListID types.Int64  `tfsdk:"source_prefixlist_id"`
	IPv4               types.List   `tfsdk:"ipv4"`
	IPv6               types.List   `tfsdk:"ipv6"`
	Version            types.Int64  `tfsdk:"version"`
	Created            types.String `tfsdk:"created"`
	Updated            types.String `tfsdk:"updated"`
}

// FlattenPrefixList maps a linodego.PrefixList into the base model fields.
func (data *PrefixListBaseModel) FlattenPrefixList(
	ctx context.Context,
	pl linodego.PrefixList,
	diags *diag.Diagnostics,
	preserveKnown bool,
) {
	data.Name = helper.KeepOrUpdateString(data.Name, pl.Name, preserveKnown)
	data.Description = helper.KeepOrUpdateString(data.Description, pl.Description, preserveKnown)
	data.Visibility = helper.KeepOrUpdateString(data.Visibility, pl.Visibility, preserveKnown)
	data.Version = helper.KeepOrUpdateInt64(data.Version, int64(pl.Version), preserveKnown)

	if pl.SourcePrefixListID != nil {
		data.SourcePrefixListID = helper.KeepOrUpdateInt64(data.SourcePrefixListID, int64(*pl.SourcePrefixListID), preserveKnown)
	} else {
		data.SourcePrefixListID = helper.KeepOrUpdateValue(data.SourcePrefixListID, types.Int64Null(), preserveKnown)
	}

	if pl.Created != nil {
		data.Created = helper.KeepOrUpdateString(data.Created, pl.Created.Format("2006-01-02T15:04:05"), preserveKnown)
	}
	if pl.Updated != nil {
		data.Updated = helper.KeepOrUpdateString(data.Updated, pl.Updated.Format("2006-01-02T15:04:05"), preserveKnown)
	}

	if pl.IPv4 != nil {
		ipv4, newDiags := types.ListValueFrom(ctx, types.StringType, *pl.IPv4)
		diags.Append(newDiags...)
		if diags.HasError() {
			return
		}
		data.IPv4 = helper.KeepOrUpdateValue(data.IPv4, ipv4, preserveKnown)
	} else {
		data.IPv4 = helper.KeepOrUpdateValue(data.IPv4, types.ListNull(types.StringType), preserveKnown)
	}

	if pl.IPv6 != nil {
		ipv6, newDiags := types.ListValueFrom(ctx, types.StringType, *pl.IPv6)
		diags.Append(newDiags...)
		if diags.HasError() {
			return
		}
		data.IPv6 = helper.KeepOrUpdateValue(data.IPv6, ipv6, preserveKnown)
	} else {
		data.IPv6 = helper.KeepOrUpdateValue(data.IPv6, types.ListNull(types.StringType), preserveKnown)
	}
}

// PrefixListDataSourceModel is the data source read model for a single prefix list.
type PrefixListDataSourceModel struct {
	ID types.String `tfsdk:"id"`
	PrefixListBaseModel
}

func (data *PrefixListDataSourceModel) parsePrefixList(
	ctx context.Context,
	pl linodego.PrefixList,
	diags *diag.Diagnostics,
) {
	data.ID = helper.KeepOrUpdateString(data.ID, strconv.Itoa(pl.ID), false)
	data.FlattenPrefixList(ctx, pl, diags, false)
}
