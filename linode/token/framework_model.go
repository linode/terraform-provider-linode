package token

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/customtypes"
)

// ResourceModel describes the Terraform resource rm model to match the
// resource schema.
type ResourceModel struct {
	Label   types.String                        `tfsdk:"label"`
	Scopes  customtypes.LinodeScopesStringValue `tfsdk:"scopes"`
	Expiry  timetypes.RFC3339                   `tfsdk:"expiry"`
	Created timetypes.RFC3339                   `tfsdk:"created"`
	Token   types.String                        `tfsdk:"token"`
	ID      types.String                        `tfsdk:"id"`
}

func (rm *ResourceModel) FlattenToken(token *linodego.Token, refresh, preserveKnown bool) {
	rm.Label = helper.KeepOrUpdateString(rm.Label, token.Label, preserveKnown)

	rm.Created = helper.KeepOrUpdateValue(
		rm.Created, timetypes.NewRFC3339TimePointerValue(token.Created), preserveKnown,
	)
	rm.Expiry = helper.KeepOrUpdateValue(
		rm.Expiry, timetypes.NewRFC3339TimePointerValue(token.Expiry), preserveKnown,
	)

	rm.ID = helper.KeepOrUpdateString(rm.ID, strconv.Itoa(token.ID), preserveKnown)

	rm.Scopes = helper.KeepOrUpdateValue(
		rm.Scopes,
		customtypes.LinodeScopesStringValue{
			StringValue: types.StringValue(token.Scopes),
		},
		preserveKnown,
	)

	if !refresh {
		rm.Token = helper.KeepOrUpdateString(rm.Token, token.Token, preserveKnown)
	}
}

func (rm *ResourceModel) CopyFrom(other ResourceModel, preserveKnown bool) {
	rm.Label = helper.KeepOrUpdateValue(rm.Label, other.Label, preserveKnown)
	rm.Created = helper.KeepOrUpdateValue(rm.Created, other.Created, preserveKnown)
	rm.Expiry = helper.KeepOrUpdateValue(rm.Expiry, other.Expiry, preserveKnown)
	rm.ID = helper.KeepOrUpdateValue(rm.ID, other.ID, preserveKnown)
	rm.Scopes = helper.KeepOrUpdateValue(rm.Scopes, other.Scopes, preserveKnown)
	rm.Token = helper.KeepOrUpdateValue(rm.Token, other.Token, preserveKnown)
}
