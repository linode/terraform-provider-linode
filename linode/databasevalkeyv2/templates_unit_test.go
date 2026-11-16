//go:build unit

package databasevalkeyv2

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func TestBasicTemplateUsesValkeyPlan(t *testing.T) {
	tmpl, err := template.ParseFiles("tmpl/basic.gotf")
	require.NoError(t, err)

	var config bytes.Buffer
	err = tmpl.ExecuteTemplate(&config, "database_valkey_v2_basic", struct {
		Label, Region, EngineID string
	}{
		Label: "test-db", Region: "us-mia", EngineID: "valkey/8.1",
	})
	require.NoError(t, err)

	file, diags := hclsyntax.ParseConfig(config.Bytes(), "basic.tf", hcl.InitialPos)
	require.False(t, diags.HasErrors(), "%s", diags.Error())
	body := file.Body.(*hclsyntax.Body)
	require.Len(t, body.Blocks, 1)
	require.Contains(t, body.Blocks[0].Body.Attributes, "type")
	plan, diags := body.Blocks[0].Body.Attributes["type"].Expr.Value(nil)
	require.False(t, diags.HasErrors(), "%s", diags.Error())
	assert.True(t, plan.RawEquals(cty.StringVal("g7-dedicated-4-2")))
}

func TestDataTemplateUsesComplexResource(t *testing.T) {
	tmpl, err := template.ParseFiles("tmpl/complex.gotf", "tmpl/data.gotf")
	require.NoError(t, err)

	var config bytes.Buffer
	err = tmpl.ExecuteTemplate(&config, "database_valkey_v2_data", struct {
		Label, Region, EngineID, Type, AllowedIP string
		ClusterSize                              int
		Updates                                  struct {
			HourOfDay, DayOfWeek, Duration int
			Frequency                      string
		}
	}{
		Label: "test-db", Region: "us-mia", EngineID: "valkey/8.1", Type: "g7-dedicated-4-2",
		AllowedIP: "10.0.0.3/32", ClusterSize: 1,
		Updates: struct {
			HourOfDay, DayOfWeek, Duration int
			Frequency                      string
		}{HourOfDay: 3, DayOfWeek: 2, Duration: 4, Frequency: "weekly"},
	})
	require.NoError(t, err)

	file, diags := hclsyntax.ParseConfig(config.Bytes(), "data.tf", hcl.InitialPos)
	require.False(t, diags.HasErrors(), "%s", diags.Error())
	body := file.Body.(*hclsyntax.Body)
	require.Len(t, body.Blocks, 2)
	assert.Equal(t, "resource", body.Blocks[0].Type)
	assert.Equal(t, []string{"linode_database_valkey_v2", "foobar"}, body.Blocks[0].Labels)

	for name, want := range map[string]cty.Value{
		"type":         cty.StringVal("g7-dedicated-4-2"),
		"cluster_size": cty.NumberIntVal(1),
		"allow_list":   cty.TupleVal([]cty.Value{cty.StringVal("10.0.0.3/32")}),
		"updates": cty.ObjectVal(map[string]cty.Value{
			"hour_of_day": cty.NumberIntVal(3),
			"day_of_week": cty.NumberIntVal(2),
			"duration":    cty.NumberIntVal(4),
			"frequency":   cty.StringVal("weekly"),
		}),
	} {
		require.Contains(t, body.Blocks[0].Body.Attributes, name)
		got, diags := body.Blocks[0].Body.Attributes[name].Expr.Value(nil)
		require.False(t, diags.HasErrors(), "%s", diags.Error())
		assert.True(t, got.RawEquals(want), "%s: got %s, want %s", name, got.GoString(), want.GoString())
	}

	assert.Equal(t, "data", body.Blocks[1].Type)
	assert.Equal(t, []string{"linode_database_valkey_v2", "foobar"}, body.Blocks[1].Labels)
}
