package token

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_token",
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
	tflog.Debug(ctx, "Create linode_token")

	var data ResourceModel
	client := r.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var expiry *time.Time
	if !data.Expiry.IsNull() && !data.Expiry.IsUnknown() {
		parsedExpiry, d := data.Expiry.ValueRFC3339Time()
		resp.Diagnostics.Append(d...)
		if d.HasError() {
			return
		}
		expiry = &parsedExpiry
	}

	createOpts := linodego.TokenCreateOptions{
		Label:  data.Label.ValueString(),
		Scopes: data.Scopes.ValueString(),
		Expiry: expiry,
	}

	tflog.Debug(ctx, "client.CreateToken(...)", map[string]any{
		"options": createOpts,
	})
	token, err := client.CreateToken(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Token creation error",
			err.Error(),
		)
		return
	}

	data.FlattenToken(token, false, true)

	// IDs should always be overridden during creation (see #1085)
	// TODO: Remove when Crossplane empty string ID issue is resolved
	data.ID = types.StringValue(strconv.Itoa(token.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read linode_token")

	client := r.Client

	var data ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, data)

	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, data.ID, resp) {
		return
	}

	id := helper.StringToInt(data.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "client.GetToken(...)")
	token, err := client.GetToken(ctx, id)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"Token No Longer Exists",
				fmt.Sprintf(
					"Removing Linode Token with ID %v from state because it no longer exists",
					data.ID,
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to Refresh the Token",
			fmt.Sprintf(
				"Error finding the specified Linode Token: %s",
				err.Error(),
			),
		)
		return
	}

	data.FlattenToken(token, true, false)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update linode_token")

	var plan, state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	ctx = populateLogAttributes(ctx, state)

	if !state.Label.Equal(plan.Label) {
		tokenIDString := state.ID.ValueString()
		tokenID := helper.StringToInt(tokenIDString, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		client := r.Client

		tflog.Trace(ctx, "client.GetToken(...)")
		token, err := client.GetToken(ctx, tokenID)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to get the token with id %v", tokenID),
				err.Error(),
			)
			return
		}

		updateOpts := token.GetUpdateOptions()
		updateOpts.Label = plan.Label.ValueString()

		tflog.Debug(ctx, "client.UpdateToken(...)", map[string]any{
			"options": updateOpts,
		})
		_, err = client.UpdateToken(ctx, token.ID, updateOpts)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to update the token with id %v", tokenID),
				err.Error(),
			)
			return
		}
		plan.FlattenToken(token, true, true)
	}

	plan.CopyFrom(state, true)

	// Workaround for Crossplane issue where ID is not
	// properly populated in plan
	// See TPT-2865 for more details
	if plan.ID.ValueString() == "" {
		plan.ID = state.ID
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete linode_token")

	var data ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, data)

	tokenID := helper.StringToInt(data.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client

	tflog.Debug(ctx, "client.DeleteToken(...)")

	err := client.DeleteToken(ctx, tokenID)
	if err != nil {
		if lErr, ok := err.(*linodego.Error); (ok && lErr.Code != 404) || !ok {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to delete the token with id %v", tokenID),
				err.Error(),
			)
		}
		return
	}

	// a settling cooldown to avoid expired tokens from being returned in listings
	// may be switched to event poller later
	time.Sleep(3 * time.Second)
}

func populateLogAttributes(ctx context.Context, model ResourceModel) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"token_id": model.ID.ValueString(),
	})
}
