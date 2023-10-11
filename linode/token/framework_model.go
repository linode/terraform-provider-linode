package token

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper/customtypes"
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

func (rm *ResourceModel) parseToken(token *linodego.Token, refresh bool) {
	rm.Created = timetypes.NewRFC3339TimePointerValue(token.Created)

	rm.Label = types.StringValue(token.Label)
	rm.Expiry = timetypes.NewRFC3339TimePointerValue(token.Expiry)

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
