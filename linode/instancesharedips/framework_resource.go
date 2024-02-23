package instancesharedips

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_instance_shared_ips",
				IDType: types.StringType,
				Schema: &frameworkResourceSchema,
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func CreateOrUpdateSharedIPs(
	ctx context.Context, client *linodego.Client, plan *ResourceModel, diags *diag.Diagnostics,
) {
	linodeID := helper.FrameworkSafeInt64ToInt(plan.LinodeID.ValueInt64(), diags)
	if diags.HasError() {
		return
	}

	createOpts := linodego.IPAddressesShareOptions{
		LinodeID: linodeID,
	}

	plan.Addresses.ElementsAs(ctx, &createOpts.IPs, false)

	tflog.Debug(ctx, "client.ShareIPAddresses(...)", map[string]interface{}{
		"options": createOpts,
	})

	err := client.ShareIPAddresses(ctx, createOpts)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("failed to update ips for linode %d", linodeID),
			err.Error(),
		)
		return
	}

	plan.FlattenSharedIPs(linodeID, createOpts.IPs, true, diags)

	// IDs should always be overridden during creation (see #1085)
	// TODO: Remove when Crossplane empty string ID issue is resolved
	plan.ID = types.StringValue(strconv.Itoa(linodeID))
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	tflog.Debug(ctx, "Create "+r.Config.Name)

	var plan ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, plan)

	CreateOrUpdateSharedIPs(ctx, client, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read "+r.Config.Name)
	client := r.Meta.Client

	var state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, state.ID, resp) {
		return
	}

	ctx = populateLogAttributes(ctx, state)

	linodeID := helper.FrameworkSafeInt64ToInt(state.LinodeID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	sharedIPs, err := GetSharedIPsForLinode(ctx, client, linodeID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get shared ips for linode %d", linodeID),
			err.Error(),
		)
		return
	}

	state.FlattenSharedIPs(linodeID, sharedIPs, true, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update "+r.Config.Name)
	var plan, state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	ctx = populateLogAttributes(ctx, plan)

	client := r.Meta.Client

	if !plan.Addresses.Equal(state.Addresses) {
		CreateOrUpdateSharedIPs(ctx, client, &plan, &resp.Diagnostics)
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
	tflog.Debug(ctx, "Delete "+r.Config.Name)
	var state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, state)

	linodeID := helper.FrameworkSafeInt64ToInt(state.LinodeID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client

	options := linodego.IPAddressesShareOptions{
		LinodeID: linodeID,
		IPs:      []string{},
	}

	tflog.Debug(ctx, "client.ShareIPAddresses(...)", map[string]interface{}{
		"options": options,
	})

	err := client.ShareIPAddresses(ctx, options)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to update shared ips for linode %d", linodeID),
			err.Error(),
		)
		return
	}
}

func GetSharedIPsForLinode(ctx context.Context, client *linodego.Client, linodeID int) ([]string, error) {
	tflog.Debug(ctx, "Enter GetSharedIPsForLinode")

	tflog.Debug(ctx, "client.GetInstanceIPAddresses(...)")
	networking, err := client.GetInstanceIPAddresses(ctx, linodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance (%d) networking: %s", linodeID, err)
	}

	result := make([]string, 0)
	for _, ip := range networking.IPv4.Shared {
		result = append(result, ip.Address)
	}

	for _, ip := range networking.IPv6.Global {
		// BGP ips will not have a route target
		if ip.RouteTarget != "" {
			continue
		}

		result = append(result, ip.Range)
	}

	return result, nil
}

func populateLogAttributes(ctx context.Context, data ResourceModel) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"linode_id": data.LinodeID.ValueInt64(),
		"id":        data.ID.ValueString(),
	})
}
