package token

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

// ResourceModel describes the Terraform resource rm model to match the
// resource schema.
type ResourceModel struct {
	Label   types.String `tfsdk:"label"`
	Scopes  types.String `tfsdk:"scopes"`
	Expiry  types.String `tfsdk:"expiry"`
	Created types.String `tfsdk:"created"`
	Token   types.String `tfsdk:"token"`
	ID      types.String `tfsdk:"id"`
}

func (rm *ResourceModel) parseComputedAttributes(token *linodego.Token, refresh bool) {
	rm.Created = types.StringValue(token.Created.Format(time.RFC3339))

	// token is too sensitive and won't appear in a GET
	// method response during a refresh of this resource.
	if !refresh {
		rm.Token = types.StringValue(token.Token)
	}
	rm.ID = types.StringValue(strconv.Itoa(token.ID))
}

func (rm *ResourceModel) parseNonComputedAttributes(token *linodego.Token) {
	// only update non-computed state values if not semantically equivalent
	rm.Label = types.StringValue(token.Label)
	if !helper.CompareTimeWithTimeString(
		token.Expiry,
		rm.Expiry.ValueString(),
		time.RFC3339,
	) {
		rm.Expiry = types.StringValue(token.Expiry.Format(time.RFC3339))
	}

	if !helper.CompareScopes(
		token.Scopes,
		rm.Scopes.ValueString(),
	) {
		rm.Scopes = types.StringValue(token.Scopes)
	}
}
