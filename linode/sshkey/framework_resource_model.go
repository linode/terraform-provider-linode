package sshkey

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper/customtypes"
)

// ResourceModel describes the Terraform resource rm model to match the
// resource schema.
type ResourceModel struct {
	Label   types.String                       `tfsdk:"label"`
	SSHKey  types.String                       `tfsdk:"ssh_key"`
	Created customtypes.RFC3339TimeStringValue `tfsdk:"created"`
	ID      types.Int64                        `tfsdk:"id"`
}

func (rm *ResourceModel) parseSSHKey(key *linodego.SSHKey) {
	rm.Label = types.StringValue(key.Label)
	rm.SSHKey = types.StringValue(key.SSHKey)
	rm.ID = types.Int64Value(int64(key.ID))

	rm.Created = customtypes.RFC3339TimeStringValue{
		StringValue: types.StringValue(key.Created.Format(time.RFC3339)),
	}
}
