package producerimagesharegroupmembers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_producer_image_share_group_members",
				Schema: &frameworkDataSourceSchema,
			},
		),
	}
}

type DataSource struct {
	helper.BaseDataSource
}

func (r *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data."+r.Config.Name)

	var data ImageShareGroupMemberFilterModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, d := filterConfig.GenerateID(data.Filters)
	if d != nil {
		resp.Diagnostics.Append(d)
		return
	}
	data.ID = id

	shareGroupID := helper.FrameworkSafeInt64ToInt(
		data.ShareGroupID.ValueInt64(),
		&resp.Diagnostics,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	result, d := filterConfig.GetAndFilter(
		ctx,
		r.Meta.Client,
		data.Filters,
		listWrapper(shareGroupID),
		data.Order,
		data.OrderBy,
	)
	if d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	data.ParseImageShareGroupMembers(helper.AnySliceToTyped[linodego.ImageShareGroupMember](result))
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func listWrapper(
	shareGroupID int,
) frameworkfilter.ListFunc {
	return func(
		ctx context.Context,
		client *linodego.Client,
		filter string,
	) ([]any, error) {
		tflog.Trace(ctx, "client.ImageShareGroupListMembers(...)")

		members, err := client.ImageShareGroupListMembers(
			ctx,
			shareGroupID,
			&linodego.ListOptions{
				Filter: filter,
			},
		)
		if err != nil {
			return nil, err
		}

		return helper.TypedSliceToAny(members), nil
	}
}
