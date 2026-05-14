//go:build unit

package networkingips

import (
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestIPAddressModelParseTags(t *testing.T) {
	ip := linodego.InstanceIP{
		Address:    "1.2.3.4",
		Gateway:    "1.2.3.1",
		SubnetMask: "255.255.255.0",
		Prefix:     32,
		Type:       "ipv4",
		Public:     true,
		RDNS:       "rdns.example.com",
		LinodeID:   100,
		Region:     "us-mia",
		Reserved:   true,
		Tags:       []string{"alpha", "beta"},
		AssignedEntity: &linodego.ReservedIPAssignedEntity{
			ID:    100,
			Label: "my-linode",
			Type:  "linode",
			URL:   "/v4/linode/instances/100",
		},
	}

	model := &IPAddressModel{}
	diags := model.ParseIP(ip)
	assert.False(t, diags.HasError())

	assert.Equal(t, "1.2.3.4", model.Address.ValueString())
	assert.Equal(t, true, model.Reserved.ValueBool())
	assert.Equal(t, 2, len(model.Tags.Elements()))
	assert.Contains(t, model.Tags.String(), "alpha")
	assert.Contains(t, model.Tags.String(), "beta")

	assert.False(t, model.AssignedEntity.IsNull())
	assert.Contains(t, model.AssignedEntity.String(), "my-linode")
	assert.Contains(t, model.AssignedEntity.String(), "/v4/linode/instances/100")
}

func TestIPAddressModelParseEmptyTags(t *testing.T) {
	ip := linodego.InstanceIP{
		Address:  "10.0.0.1",
		Type:     "ipv4",
		Public:   false,
		Reserved: false,
		Tags:     []string{},
	}

	model := &IPAddressModel{}
	diags := model.ParseIP(ip)
	assert.False(t, diags.HasError())

	assert.Equal(t, 0, len(model.Tags.Elements()))
	assert.True(t, model.AssignedEntity.IsNull())
}

func TestFilterModelParseIPAddresses(t *testing.T) {
	ips := []linodego.InstanceIP{
		{
			Address:  "1.2.3.4",
			Type:     "ipv4",
			Public:   true,
			Reserved: true,
			Tags:     []string{"tag1"},
			Region:   "us-east",
		},
		{
			Address:  "5.6.7.8",
			Type:     "ipv4",
			Public:   true,
			Reserved: false,
			Tags:     []string{},
			Region:   "us-west",
		},
	}

	model := &FilterModel{}
	diags := model.parseIPAddresses(ips)
	assert.False(t, diags.HasError())
	assert.Equal(t, 2, len(model.IPAddresses))

	assert.Equal(t, true, model.IPAddresses[0].Reserved.ValueBool())
	assert.Equal(t, 1, len(model.IPAddresses[0].Tags.Elements()))
	assert.Contains(t, model.IPAddresses[0].Tags.String(), "tag1")

	assert.Equal(t, false, model.IPAddresses[1].Reserved.ValueBool())
	assert.Equal(t, 0, len(model.IPAddresses[1].Tags.Elements()))
}
