//go:build integration || monitorlogsstream

package monitorlogsstream_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/monitorlogsstream/tmpl"
)

func TestAccDataSourceMonitorLogsStream_basic(t *testing.T) {
	acceptance.LongRunningTest(t)

	streamIDStr := getTestStreamID(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, streamIDStr),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.linode_monitor_logs_stream.foobar",
						tfjsonpath.New("id"), knownvalue.StringExact(streamIDStr)),
					statecheck.ExpectKnownValue("data.linode_monitor_logs_stream.foobar",
						tfjsonpath.New("label"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.linode_monitor_logs_stream.foobar",
						tfjsonpath.New("type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.linode_monitor_logs_stream.foobar",
						tfjsonpath.New("status"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func getTestStreamID(t testing.TB) string {
	t.Helper()

	streamID := os.Getenv("LINODE_TEST_MONITOR_LOGS_STREAM_ID")
	if streamID == "" {
		t.Skip("skipping: LINODE_TEST_MONITOR_LOGS_STREAM_ID must be set to an existing stream ID")
	}

	return streamID
}

func TestAccDataSourceMonitorLogsStream_notFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      tmpl.DataNotFound(t),
				ExpectError: regexp.MustCompile(`(?i)404|not found`),
			},
		},
	})
}
