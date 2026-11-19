package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct{}

func DataBasic(t testing.TB) string {
	return acceptance.ExecuteTemplate(t, "monitor_logs_stream_quotas_data_basic", TemplateData{})
}
