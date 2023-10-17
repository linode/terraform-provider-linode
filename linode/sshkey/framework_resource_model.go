package sshkey

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
)

// ResourceModel describes the Terraform resource rm model to match the
// resource schema.
type ResourceModel struct {
	Label   types.String      `tfsdk:"label"`
	SSHKey  types.String      `tfsdk:"ssh_key"`
	Created timetypes.RFC3339 `tfsdk:"created"`
	ID      types.String      `tfsdk:"id"`
}

func (rm *ResourceModel) parseConfiguredAttributes(key *linodego.SSHKey) {
	rm.Label = types.StringValue(key.Label)
	rm.SSHKey = types.StringValue(key.SSHKey)
}

func (rm *ResourceModel) parseComputedAttributes(key *linodego.SSHKey) {
	rm.ID = types.StringValue(strconv.Itoa(key.ID))
	rm.Created = timetypes.NewRFC3339TimePointerValue(key.Created)
}
