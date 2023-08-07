package token

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper/customtypes"
)

// ResourceModel describes the Terraform resource rm model to match the
// resource schema.
type ResourceModel struct {
	Label   types.String                        `tfsdk:"label"`
	Scopes  customtypes.LinodeScopesStringValue `tfsdk:"scopes"`
	Expiry  customtypes.RFC3339TimeStringValue  `tfsdk:"expiry"`
	Created customtypes.RFC3339TimeStringValue  `tfsdk:"created"`
	Token   types.String                        `tfsdk:"token"`
	ID      types.String                        `tfsdk:"id"`
}

func (rm *ResourceModel) parseToken(token *linodego.Token, refresh bool) {
	rm.Created = customtypes.RFC3339TimeStringValue{
		StringValue: types.StringValue(token.Created.Format(time.RFC3339)),
	}

	rm.Label = types.StringValue(token.Label)
	rm.Expiry = customtypes.RFC3339TimeStringValue{
		StringValue: types.StringValue(token.Expiry.Format(time.RFC3339)),
	}

	rm.Scopes = customtypes.LinodeScopesStringValue{
		StringValue: types.StringValue(token.Scopes),
	}

	// token is too sensitive and won't appear in a GET
	// method response during a refresh of this resource.
	if !refresh {
		rm.Token = types.StringValue(token.Token)
	}
	rm.ID = types.StringValue(strconv.Itoa(token.ID))
}
