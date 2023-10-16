//go:build unit

package rdns

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/helper/customtypes"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseConfiguredAttributes(t *testing.T) {
	ip := &linodego.InstanceIP{
		Address: "192.168.1.1",
		RDNS:    "linode.example.com",
	}

	rm := &ResourceModel{}

	rm.parseConfiguredAttributes(ip)

	assert.Equal(t, customtypes.IPAddrValue("192.168.1.1"), rm.Address)
	assert.Equal(t, types.StringValue("linode.example.com"), rm.RDNS)
}

func TestParseComputedAttributes(t *testing.T) {
	ip := &linodego.InstanceIP{
		Address: "192.168.1.1",
	}

	rm := &ResourceModel{}

	rm.parseComputedAttributes(ip)

	assert.Equal(t, types.StringValue("192.168.1.1"), rm.ID)
}
