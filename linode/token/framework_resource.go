package token

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	client *linodego.Client
}

func (data *ResourceModel) getTokenComputedAttrs(token *linodego.Token, refresh bool) {
	data.Created = types.StringValue(token.Created.Format(time.RFC3339))

	// token is too sensitive and won't appear in a GET
	// method response during a refresh of this resource.
	if !refresh {
		data.Token = types.StringValue(token.Token)
	}
	data.ID = types.StringValue(strconv.Itoa(token.ID))
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

// ResourceModel describes the Terraform resource data model to match the
// resource schema.
type ResourceModel struct {
	Label   types.String `tfsdk:"label"`
	Scopes  types.String `tfsdk:"scopes"`
	Expiry  types.String `tfsdk:"expiry"`
	Created types.String `tfsdk:"created"`
	Token   types.String `tfsdk:"token"`
	ID      types.String `tfsdk:"id"`
}

func (r *Resource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = "linode_token"
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
	var data ResourceModel
	client := r.client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	expireStr := data.Expiry.ValueString()
	dt, err := time.Parse(time.RFC3339, expireStr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid datetime string",
			fmt.Sprintf(
				"Expected expiry to be an time.RFC3339 datetime string (e.g., %s), got %s",
				time.RFC3339,
				expireStr,
			),
		)
		return
	}

	createOpts := linodego.TokenCreateOptions{
		Label:  data.Label.ValueString(),
		Scopes: data.Scopes.ValueString(),
		Expiry: &dt,
	}

	token, err := client.CreateToken(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Token creation error",
			err.Error(),
		)
		return
	}

	data.getTokenComputedAttrs(token, false)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	client := r.client

	var data ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := helper.StringToInt64(data.ID.ValueString(), resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	token, err := client.GetToken(ctx, int(id))
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

	data.getTokenComputedAttrs(token, true)

	// only update non-computed state values if not semantically equivalent
	data.Label = types.StringValue(token.Label)
	if !helper.CompareTimeWithTimeString(
		token.Expiry,
		data.Expiry.ValueString(),
		time.RFC3339,
	) {
		data.Expiry = types.StringValue(token.Expiry.Format(time.RFC3339))
	}

	if !helper.CompareScopes(
		token.Scopes,
		data.Scopes.ValueString(),
	) {
		data.Scopes = types.StringValue(token.Scopes)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if !state.Label.Equal(plan.Label) {
		tokenIDString := state.ID.ValueString()
		tokenID := int(helper.StringToInt64(tokenIDString, resp.Diagnostics))

		client := r.client
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

		_, err = client.UpdateToken(ctx, token.ID, updateOpts)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to update the token with id %v", tokenID),
				err.Error(),
			)
			return
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tokenID := int(helper.StringToInt64(data.ID.ValueString(), resp.Diagnostics))
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.client
	err := client.DeleteToken(ctx, tokenID)
	if err != nil {
		if lErr, ok := err.(*linodego.Error); (ok && lErr.Code != 404) || !ok  {
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
