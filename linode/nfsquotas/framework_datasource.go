package nfsquotas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/frameworkfilter"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_nfs_quotas",
				Schema: &frameworkDataSourceSchema,
			},
		),
	}
}

type DataSource struct {
	helper.BaseDataSource
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data."+d.Config.Name)

	var data NFSQuotaFilterModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceID := helper.FrameworkSafeInt64ToInt(data.SpaceID.ValueInt64(), &resp.Diagnostics)
	filesystemID := helper.FrameworkSafeInt64ToInt(data.FilesystemID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diag := filterConfig.GenerateID(data.Filters)
	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}
	data.ID = id

	result, diag := filterConfig.GetAndFilter(
		ctx, d.Meta.Client, data.Filters, listNFSQuotas(spaceID, filesystemID),
		// There are no API filterable fields so we don't need to provide
		// order and order_by.
		types.StringNull(), types.StringNull())
	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}

	data.FlattenNFSQuotas(helper.AnySliceToTyped[linodego.NFSQuota](result))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func listNFSQuotas(spaceID, filesystemID int) frameworkfilter.ListFunc {
	return func(ctx context.Context, client *linodego.Client, filter string) ([]any, error) {
		tflog.Trace(ctx, "client.ListNFSQuotas(...)", map[string]any{
			"space_id":      spaceID,
			"filesystem_id": filesystemID,
			"filter":        filter,
		})
		quotas, err := client.ListNFSQuotas(ctx, spaceID, filesystemID, &linodego.ListOptions{
			Filter: filter,
		})
		if err != nil {
			return nil, err
		}

		return helper.TypedSliceToAny(quotas), nil
	}
}
