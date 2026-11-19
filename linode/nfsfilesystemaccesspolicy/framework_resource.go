package nfsfilesystemaccesspolicy

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
				Name:          "linode_nfs_filesystem_access_policy",
				IDType:        types.StringType,
				Schema:        &frameworkResourceSchema,
				IsEarlyAccess: true,
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	tflog.Debug(ctx, "Import "+r.Config.Name)
	helper.ImportStateWithMultipleIDs(
		ctx,
		req,
		resp,
		[]helper.ImportableID{
			{
				Name:          "space_id",
				TypeConverter: helper.IDTypeConverterInt64,
			},
			{
				Name:          "filesystem_id",
				TypeConverter: helper.IDTypeConverterInt64,
			},
		})
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

	ctx = populateLogAttributes(ctx, plan)

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

	ctx = populateLogAttributes(ctx, state)

	spaceID, filesystemID := getIDs(state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "client.GetNFSFilesystemAccessPolicy(...)")
	policy, err := client.GetNFSFilesystemAccessPolicy(ctx, spaceID, filesystemID)
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"NFS Filesystem Access Policy No Longer Exists",
				fmt.Sprintf(
					"Removing the access policy of NFS Filesystem %d (NFS Space %d) from state "+
						"because it no longer exists",
					filesystemID, spaceID,
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Failed to get the access policy of NFS Filesystem %d (NFS Space %d).",
				filesystemID, spaceID,
			),
			err.Error(),
		)
		return
	}

	state.Flatten(ctx, spaceID, filesystemID, policy, false, &resp.Diagnostics)
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

	ctx = populateLogAttributes(ctx, plan)

	r.apply(ctx, &plan, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete removes the access policy from the Terraform state only. An NFS
// Filesystem always has exactly one access policy, so it cannot be deleted
// through the API.
func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete "+r.Config.Name)
	tflog.Info(
		ctx, "NFS Filesystem access policies cannot be deleted. "+
			"The resource has been removed from the Terraform state, but the "+
			"access policy will remain configured on the NFS Filesystem.",
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

	spaceID, filesystemID := getIDs(*plan, diags)
	if diags.HasError() {
		return
	}

	opts, shouldUpdate := plan.GetUpdateOptions(ctx, state, diags)
	if diags.HasError() {
		return
	}

	// The policy may still be provisioning or applying a previous change.
	policy, err := waitForAccessPolicy(ctx, client, spaceID, filesystemID)
	if err != nil {
		diags.AddError(
			fmt.Sprintf(
				"Failed to wait for the access policy of NFS Filesystem %d (NFS Space %d) to settle.",
				filesystemID, spaceID,
			),
			err.Error(),
		)
		return
	}

	if shouldUpdate {
		tflog.Debug(ctx, "client.UpdateNFSFilesystemAccessPolicy(...)", map[string]any{
			"options": opts,
		})
		if err := helper.RetryNFSAccessPolicyUpdate(
			ctx,
			helper.NFSAccessPolicyTimeout,
			func(ctx context.Context) error {
				_, err := client.UpdateNFSFilesystemAccessPolicy(ctx, spaceID, filesystemID, opts)
				return err
			},
		); err != nil {
			diags.AddError(
				fmt.Sprintf(
					"Failed to update the access policy of NFS Filesystem %d (NFS Space %d).",
					filesystemID, spaceID,
				),
				err.Error(),
			)
			return
		}

		policy, err = waitForAccessPolicy(ctx, client, spaceID, filesystemID)
		if err != nil {
			diags.AddError(
				fmt.Sprintf(
					"Failed to wait for the access policy of NFS Filesystem %d (NFS Space %d) to settle.",
					filesystemID, spaceID,
				),
				err.Error(),
			)
			return
		}
	}

	plan.Flatten(ctx, spaceID, filesystemID, policy, true, diags)
}

func getIDs(m ResourceModel, diags *diag.Diagnostics) (spaceID, filesystemID int) {
	spaceID = helper.FrameworkSafeInt64ToInt(m.SpaceID.ValueInt64(), diags)
	filesystemID = helper.FrameworkSafeInt64ToInt(m.FilesystemID.ValueInt64(), diags)
	return spaceID, filesystemID
}

func waitForAccessPolicy(
	ctx context.Context,
	client *linodego.Client,
	spaceID, filesystemID int,
) (*linodego.NFSFilesystemAccessPolicy, error) {
	return helper.WaitForNFSAccessPolicySettled(
		ctx,
		helper.NFSAccessPolicyTimeout,
		func(ctx context.Context) (*linodego.NFSFilesystemAccessPolicy, linodego.NFSAccessPolicyStatus, error) {
			policy, err := client.GetNFSFilesystemAccessPolicy(ctx, spaceID, filesystemID)
			if err != nil {
				return nil, "", err
			}
			return policy, policy.Status, nil
		},
	)
}

func populateLogAttributes(ctx context.Context, m ResourceModel) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"space_id":      strconv.FormatInt(m.SpaceID.ValueInt64(), 10),
		"filesystem_id": strconv.FormatInt(m.FilesystemID.ValueInt64(), 10),
	})
}
