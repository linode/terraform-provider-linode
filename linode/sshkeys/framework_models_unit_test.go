//go:build unit

package sshkeys

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
)

func TestParseSSHKeys(t *testing.T) {
	sshKeys := []linodego.SSHKey{
		{
			ID:     1,
			Label:  "Test Key",
			SSHKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
		},
		{
			ID:     2,
			Label:  "Test Key 2",
			SSHKey: "ssh-rsa DIFFERENTKEY_EAAAADAQABAAABAQC...",
		},
	}

	filterModel := &SSHKeyFilterModel{}
	filterModel.parseSSHKeys(context.Background(), sshKeys)

	assert.Equal(t, types.StringValue("1"), filterModel.SSHKeys[0].ID)
	assert.Equal(t, types.StringValue("Test Key"), filterModel.SSHKeys[0].Label)
	assert.Equal(t, types.StringValue("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC..."), filterModel.SSHKeys[0].SSHKey)

	assert.Equal(t, types.StringValue("2"), filterModel.SSHKeys[1].ID)
	assert.Equal(t, types.StringValue("Test Key 2"), filterModel.SSHKeys[1].Label)
	assert.Equal(t, types.StringValue("ssh-rsa DIFFERENTKEY_EAAAADAQABAAABAQC..."), filterModel.SSHKeys[1].SSHKey)
}
