//go:build integration || monitorlogsstreamhistory

package monitorlogsstreamhistory_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/monitorlogsstreamhistory/tmpl"
)

func TestAccDataSourceMonitorLogsStreamHistory_basic(t *testing.T) {
	acceptance.LongRunningTest(t)

	streamID := os.Getenv("LINODE_TEST_MONITOR_LOGS_STREAM_ID")
	if streamID == "" {
		t.Skip("skipping: LINODE_TEST_MONITOR_LOGS_STREAM_ID must be set to an existing stream ID")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, streamID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.linode_monitor_logs_stream_history.foobar",
						tfjsonpath.New("stream_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.linode_monitor_logs_stream_history.foobar",
						tfjsonpath.New("streams"), knownvalue.NotNull()),
				},
			},
		},
	})
}
