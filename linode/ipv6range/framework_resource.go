package ipv6range

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_ipv6_range",
				IDType: types.StringType,
				Schema: &frameworkResourceSchema,
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	prefixLength := helper.FrameworkSafeInt64ToInt(
		data.PrefixLength.ValueInt64(),
		&resp.Diagnostics,
	)
	if resp.Diagnostics.HasError() {
		return
	}
	createOpts := linodego.IPv6RangeCreateOptions{
		PrefixLength: prefixLength,
	}

	linodeIdConfigured := false
	linodeID := helper.FrameworkSafeInt64ToInt(
		data.LinodeId.ValueInt64(),
		&resp.Diagnostics,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.LinodeId.IsNull() && !data.LinodeId.IsUnknown() {
		createOpts.LinodeID = linodeID
		linodeIdConfigured = true
	} else if !data.RouteTarget.IsNull() && !data.RouteTarget.IsUnknown() {
		createOpts.RouteTarget = strings.Split(data.RouteTarget.ValueString(), "/")[0]
	} else {
		resp.Diagnostics.AddError(
			"Failed to create ipv6 range.",
			"Either linode_id or route_target must be specified.",
		)
		return
	}

	ipv6range, err := client.CreateIPv6Range(ctx, createOpts)
	if err != nil {
		if linodeIdConfigured {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to create ipv6 range for linode_id: %v", createOpts.LinodeID),
				err.Error(),
			)
		} else {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to create ipv6 range for route_target: %v", createOpts.RouteTarget),
				err.Error(),
			)
		}
		return
	}

	data.ID = types.StringValue(strings.TrimSuffix(
		ipv6range.Range,
		fmt.Sprintf("/%d", createOpts.PrefixLength)))

	// We make the GetIPv6Range API call here because the CreateIPv6Range API endpoint
	// only returns two fields for the newly created range (range and route_target).
	// We need to make a second call out to the GET endpoint to populate more
	// computed fields (region, is_bgp, linodes).
	ipv6rangeR, err := client.GetIPv6Range(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get ipv6 range when create.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.FlattenIPv6Range(ctx, ipv6rangeR, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// IDs need to always be set in the state (see #1085).
	data.ID = types.StringValue(ipv6rangeR.Range)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, data.ID, resp) {
		return
	}

	ipv6range, err := client.GetIPv6Range(ctx, data.ID.ValueString())
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && (lerr.Code == 404 || lerr.Code == 405) {
			resp.Diagnostics.AddWarning(
				"IPv6 range does not exist.",
				fmt.Sprintf("IPv6 range \"%s\" does not exist, removing from state.", data.ID.ValueString()),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Failed to get ipv6 range.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.FlattenIPv6Range(ctx, ipv6range, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ipv6range, err := client.GetIPv6Range(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get ipv6 range when update.",
			err.Error(),
		)
		return
	}

	if !state.LinodeId.Equal(plan.LinodeId) {
		linodeID := helper.FrameworkSafeInt64ToInt(
			plan.LinodeId.ValueInt64(),
			&resp.Diagnostics,
		)
		if resp.Diagnostics.HasError() {
			return
		}
		err := client.InstancesAssignIPs(ctx, linodego.LinodesAssignIPsOptions{
			Region: ipv6range.Region,
			Assignments: []linodego.LinodeIPAssignment{
				{
					LinodeID: linodeID,
					Address:  fmt.Sprintf("%s/%d", ipv6range.Range, ipv6range.Prefix),
				},
			},
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to assign ipv6 address to instance.",
				err.Error(),
			)
			return
		}

		resp.Diagnostics.Append(plan.FlattenIPv6Range(ctx, ipv6range, true)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	plan.CopyFrom(state, true)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := client.DeleteIPv6Range(ctx, data.ID.ValueString()); err != nil {
		if lerr, ok := err.(*linodego.Error); ok && (lerr.Code == 404 || lerr.Code == 405) {
			resp.Diagnostics.AddWarning(
				"IPv6 range does not exist.",
				fmt.Sprintf("IPv6 range \"%s\" does not exist, removing from state.", data.ID.ValueString()),
			)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete ipv6 range: %s", data.ID.ValueString()),
			err.Error(),
		)
		return
	}
}
