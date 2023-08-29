//go:build unit

package kernel

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseKernel(t *testing.T) {
	testKernel := linodego.LinodeKernel{
		ID:           "linode/latest-64bit",
		Label:        "Latest 64 bit (4.15.7-x86_64-linode102)",
		Version:      "4.15.7",
		Architecture: "x86_64",
		Deprecated:   false,
		KVM:          true,
		XEN:          false,
		PVOPS:        false,
		Built:        &time.Time{},
	}

	var data DataSourceModel
	data.ParseKernel(nil, &testKernel)

	assert.Equal(t, types.StringValue("linode/latest-64bit"), data.ID)
	assert.Equal(t, types.StringValue("x86_64"), data.Architecture)
	assert.Equal(t, types.BoolValue(false), data.Deprecated)
	assert.Equal(t, types.BoolValue(true), data.KVM)
	assert.Equal(t, types.StringValue("Latest 64 bit (4.15.7-x86_64-linode102)"), data.Label)
	assert.Equal(t, types.BoolValue(false), data.PVOPS)
	assert.Equal(t, types.StringValue("4.15.7"), data.Version)
	assert.Equal(t, types.BoolValue(false), data.XEN)
}
