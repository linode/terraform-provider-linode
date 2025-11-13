//go:build unit

package producerimagesharegroup

import (
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/require"
)

func TestParseImageShareGroup(t *testing.T) {
	sg := linodego.ProducerImageShareGroup{
		ID:           123,
		UUID:         "b1966cda-4083-4414-a140-45b78d48ec27",
		Label:        "my-label",
		Description:  "My description.",
		IsSuspended:  false,
		ImagesCount:  1,
		MembersCount: 1,
	}

	data := &DataSourceModel{}

	data.ParseImageShareGroup(&sg)

	require.Equal(t, int64(123), data.ID.ValueInt64())
	require.Equal(t, "b1966cda-4083-4414-a140-45b78d48ec27", data.UUID.ValueString())
	require.Equal(t, "my-label", data.Label.ValueString())
	require.Equal(t, "My description.", data.Description.ValueString())
	require.Equal(t, false, data.IsSuspended.ValueBool())
	require.Equal(t, int64(1), data.ImagesCount.ValueInt64())
	require.Equal(t, int64(1), data.MembersCount.ValueInt64())
}
