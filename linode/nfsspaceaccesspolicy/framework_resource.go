package nfsspaceaccesspolicy

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:          "linode_nfs_space_access_policy",
				IDAttr:        "space_id",
				IDType:        types.Int64Type,
				Schema:        &frameworkResourceSchema,
				IsEarlyAccess: true,
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
	tflog.Debug(ctx, "Create "+r.Config.Name)

	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, plan.SpaceID.ValueInt64())

	r.apply(ctx, &plan, nil, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
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

	ctx = populateLogAttributes(ctx, state.SpaceID.ValueInt64())

	spaceID := helper.FrameworkSafeInt64ToInt(state.SpaceID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "client.GetNFSSpaceAccessPolicy(...)")
	policy, err := client.GetNFSSpaceAccessPolicy(ctx, spaceID)
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"NFS Space Access Policy No Longer Exists",
				fmt.Sprintf(
					"Removing the access policy of NFS Space %d from state because it no longer exists",
					spaceID,
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get the access policy of NFS Space %d.", spaceID),
			err.Error(),
		)
		return
	}

	state.Flatten(ctx, spaceID, policy, false, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update "+r.Config.Name)

	var plan, state ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, plan.SpaceID.ValueInt64())

	r.apply(ctx, &plan, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete removes the access policy from the Terraform state only. An NFS Space
// always has exactly one access policy, so it cannot be deleted through the API.
func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete "+r.Config.Name)
	tflog.Info(
		ctx, "NFS Space access policies cannot be deleted. "+
			"The resource has been removed from the Terraform state, but the "+
			"access policy will remain configured on the NFS Space.",
	)
}

// apply sends the planned access policy fields to the API and flattens the
// settled policy into plan. When state is nil, all known planned fields are sent.
func (r *Resource) apply(
	ctx context.Context,
	plan *ResourceModel,
	state *ResourceModel,
	diags *diag.Diagnostics,
) {
	client := r.Meta.Client

	spaceID := helper.FrameworkSafeInt64ToInt(plan.SpaceID.ValueInt64(), diags)
	if diags.HasError() {
		return
	}

	opts, shouldUpdate := plan.GetUpdateOptions(ctx, state, diags)
	if diags.HasError() {
		return
	}

	// The policy may still be provisioning or applying a previous change.
	policy, err := waitForAccessPolicy(ctx, client, spaceID)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Failed to wait for the access policy of NFS Space %d to settle.", spaceID),
			err.Error(),
		)
		return
	}

	if shouldUpdate {
		tflog.Debug(ctx, "client.UpdateNFSSpaceAccessPolicy(...)", map[string]any{
			"options": opts,
		})
		// The API rejects updates until the NFS Space has an active tenant,
		// which is provisioned along with the first filesystem in the space.
		if err := helper.RetryNFSAccessPolicyUpdate(
			ctx,
			helper.NFSAccessPolicyTimeout,
			func(ctx context.Context) error {
				_, err := client.UpdateNFSSpaceAccessPolicy(ctx, spaceID, opts)
				return err
			},
		); err != nil {
			diags.AddError(
				fmt.Sprintf("Failed to update the access policy of NFS Space %d.", spaceID),
				fmt.Sprintf(
					"%s\n\nAn NFS Space access policy can only be updated once the NFS Space has "+
						"at least one filesystem.",
					err.Error(),
				),
			)
			return
		}

		policy, err = waitForAccessPolicy(ctx, client, spaceID)
		if err != nil {
			diags.AddError(
				fmt.Sprintf("Failed to wait for the access policy of NFS Space %d to settle.", spaceID),
				err.Error(),
			)
			return
		}
	}

	plan.Flatten(ctx, spaceID, policy, true, diags)
}

func waitForAccessPolicy(
	ctx context.Context,
	client *linodego.Client,
	spaceID int,
) (*linodego.NFSSpaceAccessPolicy, error) {
	return helper.WaitForNFSAccessPolicySettled(
		ctx,
		helper.NFSAccessPolicyTimeout,
		func(ctx context.Context) (*linodego.NFSSpaceAccessPolicy, linodego.NFSAccessPolicyStatus, error) {
			policy, err := client.GetNFSSpaceAccessPolicy(ctx, spaceID)
			if err != nil {
				return nil, "", err
			}
			return policy, policy.Status, nil
		},
	)
}

func populateLogAttributes(ctx context.Context, spaceID int64) context.Context {
	return tflog.SetField(ctx, "space_id", strconv.FormatInt(spaceID, 10))
}
