package nfsfilesystems

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
				Name:   "linode_nfs_filesystems",
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

	var data NFSFilesystemFilterModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceID := helper.FrameworkSafeInt64ToInt(data.SpaceID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	id, d2 := filterConfig.GenerateID(data.Filters)
	if d2 != nil {
		resp.Diagnostics.Append(d2)
		return
	}
	data.ID = id

	result, d2 := filterConfig.GetAndFilter(
		ctx, d.Meta.Client, data.Filters, listNFSFilesystems(spaceID),
		// There are no API filterable fields so we don't need to provide
		// order and order_by.
		types.StringNull(), types.StringNull())
	if d2 != nil {
		resp.Diagnostics.Append(d2)
		return
	}

	resp.Diagnostics.Append(
		data.FlattenNFSFilesystems(ctx, helper.AnySliceToTyped[linodego.NFSFilesystem](result), false)...,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func listNFSFilesystems(spaceID int) frameworkfilter.ListFunc {
	return func(ctx context.Context, client *linodego.Client, filter string) ([]any, error) {
		tflog.Trace(ctx, "client.ListNFSFilesystems(...)", map[string]any{
			"space_id": spaceID,
			"filter":   filter,
		})
		filesystems, err := client.ListNFSFilesystems(ctx, spaceID, &linodego.ListOptions{
			Filter: filter,
		})
		if err != nil {
			return nil, err
		}

		return helper.TypedSliceToAny(filesystems), nil
	}
}
