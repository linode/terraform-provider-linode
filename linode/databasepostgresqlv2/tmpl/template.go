package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateDataUpdates struct {
	HourOfDay, DayOfWeek, Duration int
	Frequency                      string
}

type TemplateData struct {
	Label       string
	Region      string
	EngineID    string
	Type        string
	AllowedIP   string
	ClusterSize int
	Suspended   bool

	Updates TemplateDataUpdates
}

func Basic(t testing.TB, label, region, engine, nodeType string) string {
	return acceptance.ExecuteTemplate(
		t,
		"database_postgresql_v2_basic",
		TemplateData{
			Label:    label,
			Region:   region,
			EngineID: engine,
			Type:     nodeType,
		},
	)
}

func Complex(
	t testing.TB,
	data TemplateData,
) string {
	return acceptance.ExecuteTemplate(
		t,
		"database_postgresql_v2_complex",
		data,
	)
}

func Fork(t testing.TB, label, region, engine, nodeType string) string {
	return acceptance.ExecuteTemplate(
		t,
		"database_postgresql_v2_fork",
		TemplateData{
			Label:    label,
			Region:   region,
			EngineID: engine,
			Type:     nodeType,
		},
	)
}

func Suspension(
	t testing.TB,
	data TemplateData,
) string {
	return acceptance.ExecuteTemplate(
		t,
		"database_postgresql_v2_suspension",
		data,
	)
}

func Data(
	t testing.TB,
	data TemplateData,
) string {
	return acceptance.ExecuteTemplate(
		t,
		"database_postgresql_v2_data",
		data,
	)
}
