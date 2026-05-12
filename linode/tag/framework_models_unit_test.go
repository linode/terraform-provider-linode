//go:build unit

package tag

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlattenTaggedObjects_DataSource(t *testing.T) {
	objects := linodego.TaggedObjectList{
		{
			Type: "reserved_ipv4_address",
			Data: linodego.InstanceIP{Address: "198.51.100.5"},
		},
		{
			Type: "linode",
			Data: linodego.Instance{ID: 999},
		},
	}

	model := &DataSourceModel{}
	diags := diag.Diagnostics{}
	model.FlattenTaggedObjects(context.Background(), objects, &diags)
	assert.False(t, diags.HasError())

	var items []TaggedObjectModel
	diags.Append(model.Objects.ElementsAs(context.Background(), &items, false)...)
	require.False(t, diags.HasError())
	require.Len(t, items, 2)

	assert.Equal(t, "reserved_ipv4_address", items[0].Type.ValueString())
	assert.Equal(t, "198.51.100.5", items[0].ID.ValueString())
	assert.Equal(t, "linode", items[1].Type.ValueString())
	assert.Equal(t, "999", items[1].ID.ValueString())
}
