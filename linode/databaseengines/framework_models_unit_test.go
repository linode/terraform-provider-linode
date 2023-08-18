//go:build unit

package databaseengines

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseEngines(t *testing.T) {
	mockEngines := []linodego.DatabaseEngine{
		{
			ID:      "mysql/8.0.26",
			Engine:  "mysql",
			Version: "8.0.26",
		},
		{
			ID:      "postgresql/13.0.26",
			Engine:  "postgresql",
			Version: "13.0.26",
		},
	}

	model := DatabaseEngineFilterModel{}

	model.parseEngines(mockEngines)

	assert.Len(t, model.Engines, len(mockEngines))

	assert.Equal(t, types.StringValue("mysql/8.0.26"), model.Engines[0].ID)
	assert.Equal(t, types.StringValue("mysql"), model.Engines[0].Engine)
	assert.Equal(t, types.StringValue("8.0.26"), model.Engines[0].Version)

	assert.Equal(t, types.StringValue("postgresql/13.0.26"), model.Engines[1].ID)
	assert.Equal(t, types.StringValue("postgresql"), model.Engines[1].Engine)
	assert.Equal(t, types.StringValue("13.0.26"), model.Engines[1].Version)
}
