//go:build integration || monitorlogsstreams

package monitorlogsstreams_test

import (
	"testing"
)

func TestAccDataSourceMonitorLogsStreams_basic(t *testing.T) {
	t.Skip("covered by monitorlogsstream.TestAccMonitorLogsStream_lifecycle; run TEST_SUITE=monitorlogsstream")
}
