package stackscript

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_stackscript",
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
	tflog.Debug(ctx, "Create linode_stackscript")
	var data StackScriptModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var images []string

	resp.Diagnostics.Append(data.Images.ElementsAs(ctx, &images, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := linodego.StackscriptCreateOptions{
		Label:       data.Label.ValueString(),
		Script:      data.Script.ValueString(),
		Description: data.Description.ValueString(),
		RevNote:     data.RevNote.ValueString(),
		IsPublic:    data.IsPublic.ValueBool(),
		Images:      images,
	}

	tflog.Debug(ctx, "client.CreateStackscript(...)", map[string]interface{}{
		"options": createOpts,
	})

	stackscript, err := client.CreateStackscript(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"StackScript creation error",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.ParseComputedAttributes(ctx, stackscript)...)
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
	tflog.Debug(ctx, "Read linode_stackscript")

	client := r.Meta.Client

	var data StackScriptModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, data.ID, resp) {
		return
	}

	id := helper.StringToInt(data.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "stackscript_id", id)

	tflog.Trace(ctx, "client.GetStackscript(...)")

	stackscript, err := client.GetStackscript(ctx, id)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"StackScript no longer exists.",
				fmt.Sprintf(
					"Removing Linode StackScript with ID %v from state because it no longer exists",
					data.ID,
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to refresh the Linode StackScript",
			fmt.Sprintf(
				"Error finding the specified Linode StackScript: %s",
				err.Error(),
			),
		)
		return
	}

	resp.Diagnostics.Append(data.ParseComputedAttributes(ctx, stackscript)...)
	resp.Diagnostics.Append(data.ParseNonComputedAttributes(ctx, stackscript)...)
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
	tflog.Debug(ctx, "Update linode_stackscript")

	var state StackScriptModel
	var plan StackScriptModel

	// Get the state & plan
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the ID from the plan
	stackScriptID := helper.StringToInt(state.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "stackscript_id", stackScriptID)

	// Check whether there were any changes
	isUnchanged := state.Label.Equal(plan.Label) &&
		state.Script.Equal(plan.Script) &&
		state.Description.Equal(plan.Description) &&
		state.RevNote.Equal(plan.RevNote) &&
		state.IsPublic.Equal(plan.IsPublic) &&
		state.Images.Equal(plan.Images)

	// Apply the change if necessary
	if !isUnchanged {
		r.updateStackScript(ctx, resp, plan, stackScriptID)
	}
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete linode_stackscript")
	var data StackScriptModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	stackscriptID := helper.StringToInt(data.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "stackscript_id", stackscriptID)

	tflog.Debug(ctx, "client.DeleteStackscript(...)")

	err := client.DeleteStackscript(ctx, stackscriptID)
	if err != nil {
		if lErr, ok := err.(*linodego.Error); (ok && lErr.Code != 404) || !ok {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to delete the StackScript with id %v", stackscriptID),
				err.Error(),
			)
		}
	}
}

func (r *Resource) updateStackScript(
	ctx context.Context,
	resp *resource.UpdateResponse,
	plan StackScriptModel,
	stackScriptID int,
) {
	client := r.Meta.Client

	updateOpts := linodego.StackscriptUpdateOptions{
		Label:       plan.Label.ValueString(),
		Script:      plan.Script.ValueString(),
		Description: plan.Description.ValueString(),
		RevNote:     plan.RevNote.ValueString(),
		IsPublic:    plan.IsPublic.ValueBool(),
	}

	// Special handling for images
	resp.Diagnostics.Append(plan.Images.ElementsAs(ctx, &updateOpts.Images, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "client.UpdateStackscript(...)", map[string]interface{}{
		"options": updateOpts,
	})

	stackscript, err := client.UpdateStackscript(ctx, stackScriptID, updateOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to update the StackScript with id %v", stackScriptID),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(plan.ParseComputedAttributes(ctx, stackscript)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
