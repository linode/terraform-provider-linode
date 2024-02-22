//go:build unit

package rdns

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/helper/customtypes"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestFlattenInstanceIP(t *testing.T) {
	ip := &linodego.InstanceIP{
		Address: "192.168.1.1",
		RDNS:    "linode.example.com",
	}

	rm := &ResourceModel{
		RDNS: types.StringValue("linode2.example.com"),
	}

	rm.FlattenInstanceIP(ip, false)

	assert.Equal(t, customtypes.IPAddrValue("192.168.1.1"), rm.Address)
	assert.Equal(t, types.StringValue("linode.example.com"), rm.RDNS)
}

func TestFlattenInstanceIPPreserveKnown(t *testing.T) {
	ip := &linodego.InstanceIP{
		Address: "192.168.1.1",
	}

	rm := &ResourceModel{
		ID:      types.StringUnknown(),
		Address: customtypes.IPAddrValue("192.168.1.2"),
	}

	rm.FlattenInstanceIP(ip, true)

	assert.True(t, customtypes.IPAddrValue("192.168.1.2").Equal(rm.Address))
	assert.True(t, types.StringValue("192.168.1.1").Equal(rm.ID))
}
