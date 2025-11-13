//go:build unit

package consumerimagesharegrouptoken

import (
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/require"
)

func ParseImageShareGroupToken(t *testing.T) {
	m := linodego.ImageShareGroupToken{
		TokenUUID:              "b1966cda-4083-4414-a140-45b78d48ec27",
		Status:                 "active",
		Label:                  "my-label",
		ValidForShareGroupUUID: "c52b0eda-8f5b-47c0-8bea-3881a272b117",
		ShareGroupUUID:         linodego.Pointer("c52b0eda-8f5b-47c0-8bea-3881a272b117"),
		ShareGroupLabel:        linodego.Pointer("my-sg-label"),
	}

	data := &DataSourceModel{}

	data.ParseImageShareGroupToken(&m)

	require.Equal(t, "b1966cda-4083-4414-a140-45b78d48ec27", data.TokenUUID.ValueString())
	require.Equal(t, "active", data.Status.ValueString())
	require.Equal(t, "my-label", data.Label.ValueString())
	require.Equal(t, "c52b0eda-8f5b-47c0-8bea-3881a272b117", data.ValidForShareGroupUUID.ValueString())
	require.Equal(t, "c52b0eda-8f5b-47c0-8bea-3881a272b117", data.ShareGroupUUID.ValueString())
	require.Equal(t, "my-sg-label", data.ShareGroupLabel.ValueString())
}
