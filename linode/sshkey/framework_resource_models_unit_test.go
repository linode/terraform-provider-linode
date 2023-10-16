//go:build unit

package sshkey

import (
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseConfiguredAttributes(t *testing.T) {
	key := linodego.SSHKey{
		Label:  "Test Key",
		SSHKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
	}

	rm := &ResourceModel{}
	rm.parseConfiguredAttributes(&key)

	assert.Equal(t, types.StringValue(key.Label), rm.Label)
	assert.Equal(t, types.StringValue(key.SSHKey), rm.SSHKey)
}

func TestParseComputedAttributes(t *testing.T) {
	created := time.Now()
	key := linodego.SSHKey{
		ID:      123,
		Created: &created,
	}

	rm := &ResourceModel{}
	rm.parseComputedAttributes(&key)

	assert.Equal(t, types.StringValue(strconv.Itoa(key.ID)), rm.ID)
	assert.Equal(t, types.StringValue(created.Format(time.RFC3339)), rm.Created.StringValue)
}
