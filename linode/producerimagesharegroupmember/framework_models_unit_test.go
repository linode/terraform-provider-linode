//go:build unit

package producerimagesharegroupmember

import (
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/require"
)

func TestParseImageShareGroupMember(t *testing.T) {
	m := linodego.ImageShareGroupMember{
		TokenUUID: "b1966cda-4083-4414-a140-45b78d48ec27",
		Status:    "active",
		Label:     "my-label",
	}

	data := &DataSourceModel{}

	data.ParseImageShareGroupMember(&m)

	require.Equal(t, "b1966cda-4083-4414-a140-45b78d48ec27", data.TokenUUID.ValueString())
	require.Equal(t, "active", data.Status.ValueString())
	require.Equal(t, "my-label", data.Label.ValueString())
}
