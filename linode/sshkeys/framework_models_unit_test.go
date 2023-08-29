//go:build unit

package sshkeys

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/sshkey"
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

	expectedSSHKey := []sshkey.DataSourceModel{
		{
			ID:     types.StringValue("1"),
			Label:  types.StringValue("Test Key"),
			SSHKey: types.StringValue("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC..."),
		},
		{
			ID:     types.StringValue("2"),
			Label:  types.StringValue("Test Key 2"),
			SSHKey: types.StringValue("ssh-rsa DIFFERENTKEY_EAAAADAQABAAABAQC..."),
		},
	}

	assert.Equal(t, filterModel.SSHKeys[0].SSHKey, expectedSSHKey[0].SSHKey)
	assert.Equal(t, filterModel.SSHKeys[1].SSHKey, expectedSSHKey[1].SSHKey)
}
