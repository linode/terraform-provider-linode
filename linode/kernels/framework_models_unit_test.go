//go:build unit

package kernels

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseKernels(t *testing.T) {
	kernels := []linodego.LinodeKernel{
		{
			ID:           "linode/latest-64bit",
			Label:        "Latest 64 bit (4.15.7-x86_64-linode102)",
			Version:      "4.15.7",
			Architecture: "x86_64",
			Deprecated:   false,
			KVM:          true,
			XEN:          false,
			PVOPS:        false,
			Built:        &time.Time{},
		},
	}

	var filterModel KernelFilterModel
	filterModel.parseKernels(context.Background(), kernels)

	assert.Len(t, filterModel.Kernels, len(kernels))

	assert.Equal(t, types.StringValue("linode/latest-64bit"), filterModel.Kernels[0].ID)
	assert.Equal(t, types.StringValue("x86_64"), filterModel.Kernels[0].Architecture)
}
