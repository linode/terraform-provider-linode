//go:build unit

package instancereservedipassignment

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestFlattenInstanceIPWithTags(t *testing.T) {
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
		Tags:       []string{"prod", "web"},
		AssignedEntity: &linodego.ReservedIPAssignedEntity{
			ID:    100,
			Label: "my-linode",
			Type:  "linode",
			URL:   "/v4/linode/instances/100",
		},
	}

	model := &InstanceIPModel{}
	diags := model.flattenInstanceIP(context.Background(), ip, false)
	assert.False(t, diags.HasError(), "unexpected diags: %v", diags)

	assert.Equal(t, "1.2.3.4", model.Address.ValueString())
	assert.Equal(t, true, model.Reserved.ValueBool())
	assert.Equal(t, 2, len(model.Tags.Elements()))
	assert.Contains(t, model.Tags.String(), "prod")
	assert.Contains(t, model.Tags.String(), "web")

	assert.False(t, model.AssignedEntity.IsNull())
	assert.Contains(t, model.AssignedEntity.String(), "my-linode")
	assert.Contains(t, model.AssignedEntity.String(), "/v4/linode/instances/100")
}

func TestFlattenInstanceIPEmptyTags(t *testing.T) {
	ip := linodego.InstanceIP{
		Address:  "10.0.0.1",
		Type:     "ipv4",
		Public:   false,
		Reserved: false,
		Tags:     []string{},
	}

	model := &InstanceIPModel{}
	diags := model.flattenInstanceIP(context.Background(), ip, false)
	assert.False(t, diags.HasError())

	assert.Equal(t, false, model.Reserved.ValueBool())
	assert.Equal(t, 0, len(model.Tags.Elements()))
	assert.True(t, model.AssignedEntity.IsNull())
}

func TestCopyFromPreservesTags(t *testing.T) {
	ip := linodego.InstanceIP{
		Address:  "1.2.3.4",
		Gateway:  "1.2.3.1",
		Type:     "ipv4",
		Public:   true,
		Reserved: true,
		Tags:     []string{"tag1"},
	}

	source := &InstanceIPModel{}
	diags := source.flattenInstanceIP(context.Background(), ip, false)
	assert.False(t, diags.HasError())

	dest := &InstanceIPModel{}
	diags = dest.CopyFrom(context.Background(), *source, false)
	assert.False(t, diags.HasError())

	assert.Equal(t, true, dest.Reserved.ValueBool())
	assert.Equal(t, 1, len(dest.Tags.Elements()))
	assert.Contains(t, dest.Tags.String(), "tag1")
}
