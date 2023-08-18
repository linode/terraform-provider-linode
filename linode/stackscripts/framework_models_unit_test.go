package stackscripts

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseStackscripts(t *testing.T) {
	createdTime := time.Date(2023, time.August, 17, 12, 0, 0, 0, time.UTC)
	updatedTime := time.Date(2023, time.August, 17, 14, 0, 0, 0, time.UTC)

	stackscriptsData := []linodego.Stackscript{
		{
			ID:                10079,
			Username:          "myuser",
			Label:             "a-stackscript",
			Description:       "This StackScript installs and configures MySQL\n",
			RevNote:           "Set up MySQL",
			IsPublic:          true,
			Images:            []string{"linode/debian9", "linode/debian8"},
			DeploymentsActive: 1,
			Mine:              true,
			Created:           &createdTime,
			Updated:           &updatedTime,
			Script:            "\"#!/bin/bash\"\n",
			UserDefinedFields: &[]linodego.StackscriptUDF{
				{
					Default: "",
					Example: "hunter2",
					Label:   "Enter the password",
					ManyOf:  "avalue,anothervalue,thirdvalue",
					Name:    "DB_PASSWORD",
					OneOf:   "avalue,anothervalue,thirdvalue",
				},
			},
			UserGravatarID: "a445b305abda30ebc766bc7fda037c37",
		},
	}

	model := &StackscriptFilterModel{}
	diagnostics := model.parseStackscripts(context.Background(), stackscriptsData)

	assert.False(t, diagnostics.HasError(), "No errors should occur during parsing")

	for i, stackscript := range stackscriptsData {
		// Non computed attrs assertions
		assert.Contains(t, model.Stackscripts[i].ID.String(), strconv.Itoa(stackscript.ID))
		assert.Equal(t, types.StringValue(stackscript.Description), model.Stackscripts[i].Description)
		assert.Equal(t, types.StringValue(stackscript.Script), model.Stackscripts[i].Script)
		assert.Equal(t, types.StringValue(stackscript.RevNote), model.Stackscripts[i].RevNote)
		assert.Equal(t, types.BoolValue(stackscript.IsPublic), model.Stackscripts[i].IsPublic)

		for _, image := range stackscript.Images {
			assert.Contains(t, model.Stackscripts[i].Images.String(), image)
		}

		// Computed attr assertions
		assert.Equal(t, types.Int64Value(int64(stackscript.DeploymentsActive)), model.Stackscripts[i].DeploymentsActive)
		assert.Equal(t, types.Int64Value(int64(stackscript.DeploymentsTotal)), model.Stackscripts[i].DeploymentsTotal)
		assert.Equal(t, types.StringValue(stackscript.UserGravatarID), model.Stackscripts[i].UserGravatarID)
		assert.Equal(t, types.StringValue(stackscript.Username), model.Stackscripts[i].Username)
	}
}
