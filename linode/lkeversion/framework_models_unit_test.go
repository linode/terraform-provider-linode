//go:build unit

package lkeversion

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseLKEVersion(t *testing.T) {
	mockLKEVersion := &linodego.LKEVersion{
		ID: "1.31",
	}

	data := &DataSourceModel{}

	diags := data.ParseLKEVersion(mockLKEVersion)
	assert.False(t, diags.HasError(), "Unexpected error")

	assert.Equal(t, types.StringValue("1.31"), data.ID)
}

func TestParseLKETierVersion(t *testing.T) {
	mockLKETierVersion := &linodego.LKETierVersion{
		ID:   "1.31",
		Tier: "enterprise",
	}

	data := &DataSourceModel{}

	diags := data.ParseLKETierVersion(mockLKETierVersion)
	assert.False(t, diags.HasError(), "Unexpected error")

	assert.Equal(t, types.StringValue("1.31"), data.ID)
	assert.Equal(t, types.StringValue("enterprise"), data.Tier)
}
