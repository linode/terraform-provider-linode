//go:build unit

package stackscript

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestFlattenStackScript(t *testing.T) {
	createdTime := time.Date(2023, time.August, 17, 12, 0, 0, 0, time.UTC)
	updatedTime := time.Date(2023, time.August, 17, 14, 0, 0, 0, time.UTC)

	stackscriptData := &linodego.Stackscript{
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
	}

	model := &StackScriptModel{}
	diagnostics := model.FlattenStackScript(stackscriptData, false)

	assert.False(t, diagnostics.HasError(), "No errors should occur during parsing")

	assert.Equal(t, types.StringValue("a-stackscript"), model.Label)
	assert.Equal(t, types.StringValue("\"#!/bin/bash\"\n"), model.Script)
	assert.Equal(t, types.StringValue("This StackScript installs and configures MySQL\n"), model.Description)
	assert.Equal(t, types.StringValue("Set up MySQL"), model.RevNote)
	assert.Equal(t, types.BoolValue(true), model.IsPublic)

	for _, image := range stackscriptData.Images {
		assert.Contains(t, model.Images.String(), image)
	}
}

func TestFlattenStackScriptPreservingKnown(t *testing.T) {
	createdTime := time.Date(2023, time.August, 17, 12, 0, 0, 0, time.UTC)
	updatedTime := time.Date(2023, time.August, 17, 14, 0, 0, 0, time.UTC)

	stackscriptData := &linodego.Stackscript{
		ID:                10079,
		Username:          "myuser",
		Label:             "a-stackscript",
		Description:       "This StackScript installs and configures MySQL\n",
		RevNote:           "Set up MySQL",
		DeploymentsTotal:  12,
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
	}

	model := &StackScriptModel{
		RevNote:          types.StringValue("Set up PostgreSQL"),
		DeploymentsActive: types.Int64Unknown(),
	}

	diags := model.FlattenStackScript(stackscriptData, true)
	assert.False(t, diags.HasError())

	assert.True(t, model.DeploymentsActive.Equal(types.Int64Value(int64(stackscriptData.DeploymentsActive))))
	assert.False(t, model.RevNote.Equal(types.StringValue(stackscriptData.RevNote)))
}
