package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	StreamID string
}

func DataBasic(t testing.TB, streamID string) string {
	return acceptance.ExecuteTemplate(t, "monitor_logs_stream_history_data_basic", TemplateData{
		StreamID: streamID,
	})
}
