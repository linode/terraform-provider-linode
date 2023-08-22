//go:build unit

package vlan

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseVLANs(t *testing.T) {
	vlans := []linodego.VLAN{
		{
			Label:   "VLAN 1",
			Linodes: []int{123, 456},
			Region:  "us-east",
			Created: &time.Time{},
		},
		{
			Label:   "VLAN 2",
			Linodes: []int{789},
			Region:  "us-west",
			Created: nil,
		},
	}

	data := &VLANsFilterModel{}
	diags := data.parseVLANs(context.Background(), vlans)

	assert.Empty(t, diags.HasError(), "No errors expected")
	assert.Equal(t, len(vlans), len(data.VLANs), "Number of parsed VLANs should match")
}

func TestParseVLAN(t *testing.T) {
	vlan := linodego.VLAN{
		Label:   "Test VLAN",
		Linodes: []int{123, 456},
		Region:  "us-east",
		Created: &time.Time{},
	}

	data := &VLANModel{}
	diags := data.parseVLAN(context.Background(), vlan)

	assert.Empty(t, diags.HasError(), "No errors expected")
	assert.Equal(t, types.StringValue(vlan.Label), data.Label)
	for _, linodeId := range vlan.Linodes {
		assert.Contains(t, data.Linodes.String(), strconv.Itoa(linodeId))
	}
	assert.Equal(t, types.StringValue(vlan.Region), data.Region)
	assert.NotNil(t, data.Created)
}
