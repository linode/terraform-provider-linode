//go:build unit

package monitoralertdefinition

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestAlertDefinitionModel_FlattenAlertDefinition(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	alertChannels := []linodego.AlertChannelEnvelope{
		{
			ID:    1,
			Label: "Email",
			Type:  "email",
			URL:   "mailto:test@example.com",
		},
	}

	ruleCriteria := linodego.RuleCriteria{
		Rules: []linodego.Rule{
			{
				AggregateFunction: "avg",
				DimensionFilters: []linodego.DimensionFilter{
					{
						DimensionLabel: "host",
						Operator:       "eq",
						Value:          "web01",
					},
				},
				Metric:    "cpu",
				Operator:  "gt",
				Threshold: 80.0,
			},
		},
	}

	triggerConditions := linodego.TriggerConditions{
		CriteriaCondition:       "AND",
		EvaluationPeriodSeconds: 60,
		PollingIntervalSeconds:  30,
		TriggerOccurrences:      2,
	}

	alertDef := &linodego.AlertDefinition{
		ID:                42,
		ServiceType:       "monitoring",
		ChannelIDs:        []int{1},
		Description:       "High CPU usage",
		EntityIDs:         []string{"123", "456"},
		Label:             "CPU Alert",
		Status:            "active",
		WaitFor:           true,
		Severity:          2,
		RuleCriteria:      ruleCriteria,
		TriggerConditions: triggerConditions,
		Type:              "user",
		HasMoreResources:  false,
		AlertChannels:     alertChannels,
		Created:           &now,
		Updated:           &now,
		CreatedBy:         "admin",
		UpdatedBy:         "admin",
		Class:             "system",
	}

	var model AlertDefinitionModel
	diags := model.FlattenAlertDefinition(ctx, alertDef, false)

	assert.False(t, diags.HasError(), "Diagnostics should not have errors")
	assert.Equal(t, types.Int64Value(42), model.ID)
	assert.Equal(t, types.StringValue("monitoring"), model.ServiceType)
	assert.Equal(t, types.StringValue("High CPU usage"), model.Description)
	assert.Equal(t, types.StringValue("CPU Alert"), model.Label)
	assert.Equal(t, types.StringValue("active"), model.Status)
	assert.Equal(t, types.Int64Value(2), model.Severity)
	assert.Equal(t, types.StringValue("user"), model.Type)
	assert.Equal(t, types.BoolValue(false), model.HasMoreResources)
	assert.Equal(t, types.StringValue("admin"), model.CreatedBy)
	assert.Equal(t, types.StringValue("admin"), model.UpdatedBy)
	assert.Equal(t, types.StringValue("system"), model.Class)
	assert.Equal(t, timetypes.NewRFC3339TimePointerValue(&now), model.Created)
	assert.Equal(t, timetypes.NewRFC3339TimePointerValue(&now), model.Updated)
	assert.False(t, model.ChannelIDs.IsNull())
	assert.Equal(t, 1, len(model.ChannelIDs.Elements()))
	assert.False(t, model.AlertChannels.IsNull())
	assert.Equal(t, 1, len(model.AlertChannels.Elements()))
	assert.False(t, model.EntityIDs.IsNull())
	assert.Equal(t, 2, len(model.EntityIDs.Elements()))
	assert.False(t, model.RuleCriteria.IsNull())
	assert.False(t, model.TriggerConditions.IsNull())
}
