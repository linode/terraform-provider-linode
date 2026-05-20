//go:build integration || monitorlogsstreamhistory

package monitorlogsstreamhistory_test

import (
	"testing"
)

func TestAccDataSourceMonitorLogsStreamHistory_basic(t *testing.T) {
	t.Skip("covered by monitorlogsstream.TestAccMonitorLogsStream_lifecycle; run TEST_SUITE=monitorlogsstream")
}
