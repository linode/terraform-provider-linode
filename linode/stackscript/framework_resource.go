package stackscript

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	client *linodego.Client
}

func (r *Resource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	meta := helper.GetResourceMeta(req, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	r.client = meta.Client
}

func (r *Resource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = "linode_stackscript"
}

func (r *Resource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = frameworkResourceSchema
}

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data StackScriptModel
	client := r.client

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

	stackscript, err := client.CreateStackscript(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"StackScript creation error",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.parseStackScript(ctx, stackscript)...)
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
	client := r.client

	var data StackScriptModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := helper.StringToInt64(data.ID.ValueString(), resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	stackscript, err := client.GetStackscript(ctx, int(id))
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

	resp.Diagnostics.Append(data.parseStackScript(ctx, stackscript)...)
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
	var state StackScriptModel
	var plan StackScriptModel

	// Get the state & plan
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the ID from the plan
	stackScriptID := int(helper.StringToInt64(state.ID.ValueString(), resp.Diagnostics))
	if resp.Diagnostics.HasError() {
		return
	}

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
	var data StackScriptModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	stackscriptID := int(helper.StringToInt64(data.ID.ValueString(), resp.Diagnostics))
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.client
	err := client.DeleteStackscript(ctx, stackscriptID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete the StackScript with id %v", stackscriptID),
			err.Error(),
		)
		return
	}
}

func (r *Resource) updateStackScript(
	ctx context.Context,
	resp *resource.UpdateResponse,
	plan StackScriptModel,
	stackScriptID int,
) {
	client := r.client

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

	stackScript, err := client.UpdateStackscript(ctx, stackScriptID, updateOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to update the StackScript with id %v", stackScriptID),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(plan.parseStackScript(ctx, stackScript)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
