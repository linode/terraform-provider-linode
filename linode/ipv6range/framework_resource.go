package ipv6range

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

const ipv6rangeNotExist string = "IPv6 range does not exist."

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

	createOpts := linodego.IPv6RangeCreateOptions{
		PrefixLength: int(data.PrefixLength.ValueInt64()),
	}

	if !data.LinodeId.IsNull() && !data.LinodeId.IsUnknown() {
		createOpts.LinodeID = int(data.LinodeId.ValueInt64())
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
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to create ipv6 range for linode_id: %v, or route_target: %v",
				createOpts.LinodeID, createOpts.RouteTarget),
			err.Error(),
		)
		return
	}

	data.ID = types.StringValue(strings.TrimSuffix(
		ipv6range.Range,
		fmt.Sprintf("/%d", createOpts.PrefixLength)))

	// We make the GetIPv6Range API call here because the CreateIPv6Range API endpoint
	// only returns two fields for the newly created range (range and route_target).
	// We need to make a second call out to the GET endpoint to populate more
	// computed fields (region, is_bgp, linodes).
	ipv6rangeR, diag := getIPv6Range(ctx, client, data.ID.ValueString())

	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}

	resp.Diagnostics.Append(data.parseIPv6RangeResourceDataComputedAttrs(ctx, ipv6rangeR)...)
	if resp.Diagnostics.HasError() {
		return
	}
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

	ipv6range, d := getIPv6Range(ctx, client, data.ID.ValueString())

	if d != nil {
		if d.Summary() == ipv6rangeNotExist {
			resp.Diagnostics.Append(diag.NewWarningDiagnostic(
				"Removing the resource from the state.",
				fmt.Sprintf("IPv6 range \"%s\" does not exist, removing from state.", data.ID.ValueString()),
			))
			resp.State.RemoveResource(ctx)
		}
		resp.Diagnostics.Append(d)
		return
	}

	data.parseIPv6RangeResourceDataNonComputedAttrs(ipv6range)
	resp.Diagnostics.Append(data.parseIPv6RangeResourceDataComputedAttrs(ctx, ipv6range)...)
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
		err := client.InstancesAssignIPs(ctx, linodego.LinodesAssignIPsOptions{
			Region: ipv6range.Region,
			Assignments: []linodego.LinodeIPAssignment{
				{
					LinodeID: int(plan.LinodeId.ValueInt64()),
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

		resp.Diagnostics.Append(plan.parseIPv6RangeResourceDataComputedAttrs(ctx, ipv6range)...)
		if resp.Diagnostics.HasError() {
			return
		}
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	}
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
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
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

func getIPv6Range(
	ctx context.Context,
	client *linodego.Client,
	id string,
) (*linodego.IPv6Range, diag.Diagnostic) {
	ipv6range, err := client.GetIPv6Range(ctx, id)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && (lerr.Code == 404 || lerr.Code == 405) {
			return nil, diag.NewWarningDiagnostic(
				ipv6rangeNotExist,
				fmt.Sprintf("IPv6 range \"%s\" does not exist.", id),
			)
		}
		return nil, diag.NewErrorDiagnostic(
			"Failed to get ipv6 range.",
			err.Error(),
		)
	}
	return ipv6range, nil
}
