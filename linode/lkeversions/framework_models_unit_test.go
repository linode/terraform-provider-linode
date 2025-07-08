//go:build unit

package lkeversions

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseLKEVersions(t *testing.T) {
	mockVersions := []linodego.LKEVersion{
		{
			ID: "1.30",
		},
		{
			ID: "1.31",
		},
		{
			ID: "1.32",
		},
	}

	data := &DataSourceModel{}

	diags := data.parseLKEVersions(mockVersions)
	assert.False(t, diags.HasError(), "Unexpected error")

	for i, mockVersion := range mockVersions {
		assert.Equal(t, data.Versions[i].ID, types.StringValue(mockVersion.ID), "ID doesn't match")
	}
}

func TestParseLKETierVersions(t *testing.T) {
	mockVersions := []linodego.LKETierVersion{
		{
			ID:   "1.30",
			Tier: "enterprise",
		},
		{
			ID:   "1.31",
			Tier: "enterprise",
		},
		{
			ID:   "1.32",
			Tier: "standard",
		},
	}

	data := &DataSourceModel{}

	diags := data.parseLKETierVersions(mockVersions)
	assert.False(t, diags.HasError(), "Unexpected error")

	for i, mockVersion := range mockVersions {
		assert.Equal(t, data.Versions[i].ID, types.StringValue(mockVersion.ID), "ID doesn't match")
		assert.Equal(
			t,
			data.Versions[i].Tier,
			types.StringValue(string(mockVersion.Tier)),
			"Tier doesn't match",
		)
	}
}
