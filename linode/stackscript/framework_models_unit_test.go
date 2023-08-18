package stackscript

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseNonComputedAttributes(t *testing.T) {
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
	diagnostics := model.ParseNonComputedAttributes(context.Background(), stackscriptData)

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

func TestParseComputedAttributes(t *testing.T) {
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

	model := &StackScriptModel{}
	diagnostics := model.ParseComputedAttributes(context.Background(), stackscriptData)

	assert.False(t, diagnostics.HasError(), "No errors should occur during parsing")

	assert.Equal(t, types.StringValue("10079"), model.ID)
	assert.Equal(t, types.Int64Value(1), model.DeploymentsActive)
	assert.Equal(t, types.StringValue("a445b305abda30ebc766bc7fda037c37"), model.UserGravatarID)
	assert.Equal(t, types.Int64Value(12), model.DeploymentsTotal)
	assert.Equal(t, types.StringValue("myuser"), model.Username)

	assert.NotNil(t, model.Created)
	assert.NotNil(t, model.Updated)

	udfs := model.UserDefinedFields
	assert.Contains(t, udfs.String(), "Enter the password")
	assert.Contains(t, udfs.String(), "DB_PASSWORD")
	assert.Contains(t, udfs.String(), "avalue,anothervalue,thirdvalue")
}
