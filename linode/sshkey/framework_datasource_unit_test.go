//go:build unit

package sshkey

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
)

func TestParseSSHKey(t *testing.T) {
	created := time.Now().UTC().Truncate(time.Second)
	key := linodego.SSHKey{
		ID:      123,
		Created: &created,
		Label:   "Test Key",
		SSHKey:  "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
	}

	data := &DataSourceModel{}
	data.ParseSSHKey(&key)

	assert.Equal(t, types.StringValue("123"), data.ID)
	assert.Equal(t, types.StringValue("Test Key"), data.Label)
	assert.Equal(t, types.StringValue("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC..."), data.SSHKey)
	assert.Equal(t, types.StringValue(created.Format(time.RFC3339)), data.Created.StringValue)
}
