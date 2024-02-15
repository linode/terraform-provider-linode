package sshkey

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

// ResourceModel describes the Terraform resource rm model to match the
// resource schema.
type ResourceModel struct {
	Label   types.String      `tfsdk:"label"`
	SSHKey  types.String      `tfsdk:"ssh_key"`
	Created timetypes.RFC3339 `tfsdk:"created"`
	ID      types.String      `tfsdk:"id"`
}

func (rm *ResourceModel) FlattenSSHKey(key *linodego.SSHKey, preserveKnown bool) {
	rm.Label = helper.KeepOrUpdateString(rm.Label, key.Label, preserveKnown)
	rm.SSHKey = helper.KeepOrUpdateString(rm.SSHKey, key.SSHKey, preserveKnown)
	rm.ID = helper.KeepOrUpdateString(rm.ID, strconv.Itoa(key.ID), preserveKnown)
	rm.Created = helper.KeepOrUpdateValue(
		rm.Created, timetypes.NewRFC3339TimePointerValue(key.Created), preserveKnown,
	)
}

func (rm *ResourceModel) CopyFrom(other ResourceModel, preserveKnown bool) {
	rm.Label = helper.KeepOrUpdateValue(rm.Label, other.Label, preserveKnown)
	rm.SSHKey = helper.KeepOrUpdateValue(rm.SSHKey, other.SSHKey, preserveKnown)
	rm.ID = helper.KeepOrUpdateValue(rm.ID, other.ID, preserveKnown)
	rm.Created = helper.KeepOrUpdateValue(rm.Created, other.Created, preserveKnown)
}
