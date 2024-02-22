//go:build unit

package sshkey

import (
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseConfiguredAttributes(t *testing.T) {
	created := time.Now()
	key := linodego.SSHKey{
		ID:      123,
		Created: &created,
		Label:   "Test Key",
		SSHKey:  "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
	}

	rm := &ResourceModel{}
	rm.FlattenSSHKey(&key, false)

	assert.Equal(t, types.StringValue(key.Label), rm.Label)
	assert.Equal(t, types.StringValue(key.SSHKey), rm.SSHKey)
	assert.Equal(t, types.StringValue(strconv.Itoa(key.ID)), rm.ID)
	assert.Equal(t, types.StringValue(created.Format(time.RFC3339)), rm.Created.StringValue)
}

func TestParseComputedAttributes(t *testing.T) {
	created := time.Now()
	key := linodego.SSHKey{
		ID:      123,
		Created: &created,
	}

	rm := &ResourceModel{
		ID:      types.StringUnknown(),
		Created: timetypes.NewRFC3339TimeValue(created.Add(24 * time.Hour)),
	}
	rm.FlattenSSHKey(&key, true)

	assert.True(t, types.StringValue("123").Equal(rm.ID))
	assert.False(t, rm.Created.Equal(timetypes.NewRFC3339TimeValue(created)))
}
