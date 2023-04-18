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

func (data *ResourceModel) parseToken(token *linodego.Token) {
	data.Created = types.StringValue(token.Created.Format(time.RFC3339))
	data.Expiry = types.StringValue(token.Expiry.Format(time.RFC3339))
	data.Label = types.StringValue(token.Label)
	data.Scopes = types.StringValue(token.Scopes)
	data.Token = types.StringValue(token.Token)
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

	meta, ok := req.ProviderData.(*helper.FrameworkProviderMeta)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected *http.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

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

	data.parseToken(token)
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
	}

	data.parseToken(token)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var data ResourceModel
	var tokenIDString string
	resp.Diagnostics.Append(
		req.State.GetAttribute(ctx, path.Root("id"), &tokenIDString)...,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tokenID := int(helper.StringToInt64(tokenIDString, resp.Diagnostics))
	if resp.Diagnostics.HasError() {
		return
	}

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
	plannedTokenLabel := data.Label.ValueString()

	if updateOpts.Label != plannedTokenLabel {
		updateOpts.Label = plannedTokenLabel
		token, err = client.UpdateToken(ctx, token.ID, updateOpts)

		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to update the token with id %v", tokenID),
				err.Error(),
			)
			return
		}

		data.parseToken(token)
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
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
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete the token with id %v", tokenID),
			err.Error(),
		)
		return
	}

	// a settling cooldown to avoid expired tokens from being returned in listings
	// may be switched to event poller later
	time.Sleep(3 * time.Second)
}
