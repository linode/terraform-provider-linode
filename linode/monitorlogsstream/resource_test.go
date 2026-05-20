//go:build integration || monitorlogsstream

package monitorlogsstream_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/monitorlogsstream/tmpl"
)

const (
	resourceName          = "linode_monitor_logs_stream.foobar"
	dataSourceName        = "data.linode_monitor_logs_stream.by_id"
	dataSourcesListName   = "data.linode_monitor_logs_streams.filtered"
	dataSourceHistoryName = "data.linode_monitor_logs_stream_history.hist"

	// oneStreamLimitEnvVar enables the one-stream-per-account pre-check when true (default).
	oneStreamLimitEnvVar = "ONE_STREAM_LIMIT"

	// Streams can take up to 30 minutes to provision; poll until editable/deletable.
	logStreamPollInterval     = 30 * time.Second
	logStreamProvisionTimeout = 45 * time.Minute
)

func init() {
	resource.AddTestSweepers("linode_monitor_logs_stream", &resource.Sweeper{
		Name: "linode_monitor_logs_stream",
		F:    sweep,
	})
}

func sweep(prefix string) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		log.Fatal(fmt.Errorf("Error getting client: %s", err))
	}

	streams, err := client.ListLogStreams(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("Error listing logs streams: %s", err)
	}

	for _, stream := range streams {
		if !acceptance.ShouldSweep(prefix, stream.Label) {
			continue
		}
		err := client.DeleteLogStream(context.Background(), stream.ID)
		if err != nil {
			return fmt.Errorf("Error destroying stream %d during sweep: %s", stream.ID, err)
		}
	}

	return nil
}

// oneStreamLimitEnabled reports whether the one-stream-per-account pre-check is
// active. Defaults to true when ONE_STREAM_LIMIT is unset.
func oneStreamLimitEnabled() bool {
	v, ok := os.LookupEnv(oneStreamLimitEnvVar)
	if !ok || v == "" {
		return true
	}

	enabled, err := strconv.ParseBool(v)
	if err != nil {
		return true
	}

	return enabled
}

// requireNoExistingStreams skips the test if a stream already exists on the
// account when ONE_STREAM_LIMIT is true (default). Set ONE_STREAM_LIMIT=false to
// disable this check.
func requireNoExistingStreams(t *testing.T) {
	t.Helper()

	if !oneStreamLimitEnabled() {
		return
	}

	client, err := acceptance.GetTestClient()
	if err != nil {
		t.Fatalf("Error getting client: %s", err)
	}

	existing, err := client.ListLogStreams(context.Background(), nil)
	if err != nil {
		t.Fatalf("Error listing log streams: %s", err)
	}

	if len(existing) > 0 {
		ids := make([]int, len(existing))
		for i, s := range existing {
			ids[i] = s.ID
		}
		t.Skipf(
			"existing stream(s) on account (IDs: %v); only one stream allowed per account (set %s=false to skip this check)",
			ids, oneStreamLimitEnvVar,
		)
	}
}

func getRegionForStreamTest(t *testing.T) string {
	t.Helper()

	region, err := acceptance.GetRandomRegionWithCaps(
		[]string{linodego.CapabilityObjectStorage}, "core")
	if err != nil {
		t.Skipf("could not get region with Object Storage: %s", err)
	}

	return region
}

func streamIsReady(status linodego.StreamStatus) bool {
	switch status {
	case linodego.StreamStatusActive, linodego.StreamStatusInactive:
		return true
	case linodego.StreamStatusProvisioning, linodego.StreamStatusDeactivating:
		return false
	default:
		return false
	}
}

func logStreamIDFromState(s *terraform.State, name string) (int, error) {
	rs, ok := s.RootModule().Resources[name]
	if !ok {
		return 0, fmt.Errorf("resource not found: %s", name)
	}

	idStr, ok := rs.Primary.Attributes["id"]
	if !ok {
		idStr = rs.Primary.ID
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid stream id %q: %w", idStr, err)
	}

	return id, nil
}

// waitForLogStreamReady polls until the stream is no longer provisioning or deactivating.
func waitForLogStreamReady(ctx context.Context, client *linodego.Client, id int, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(logStreamPollInterval)
	defer ticker.Stop()

	var lastStatus linodego.StreamStatus

	for {
		stream, err := client.GetLogStream(ctx, id)
		if err != nil {
			return fmt.Errorf("get log stream %d: %w", id, err)
		}

		lastStatus = stream.Status
		if streamIsReady(stream.Status) {
			return nil
		}

		log.Printf("waiting for log stream %d to finish provisioning (status=%s)", id, stream.Status)

		select {
		case <-ctx.Done():
			return fmt.Errorf(
				"timed out waiting for log stream %d to be ready (last status=%s): %w",
				id, lastStatus, ctx.Err(),
			)
		case <-ticker.C:
		}
	}
}

func testAccCheckLogStreamReady(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		id, err := logStreamIDFromState(s, name)
		if err != nil {
			return err
		}

		client, err := acceptance.GetTestClient()
		if err != nil {
			return fmt.Errorf("error getting client: %w", err)
		}

		return waitForLogStreamReady(context.Background(), client, id, logStreamProvisionTimeout)
	}
}

// lifecycleStateChecks returns state checks for the resource and related data sources.
func lifecycleStateChecks(label string) []statecheck.StateCheck {
	return []statecheck.StateCheck{
		statecheck.ExpectKnownValue(resourceName,
			tfjsonpath.New("label"), knownvalue.StringExact(label)),
		statecheck.ExpectKnownValue(resourceName,
			tfjsonpath.New("type"), knownvalue.StringExact("audit_logs")),
		statecheck.ExpectKnownValue(resourceName,
			tfjsonpath.New("status"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(resourceName,
			tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(dataSourceName,
			tfjsonpath.New("label"), knownvalue.StringExact(label)),
		statecheck.ExpectKnownValue(dataSourceName,
			tfjsonpath.New("type"), knownvalue.StringExact("audit_logs")),
		statecheck.ExpectKnownValue(dataSourceName,
			tfjsonpath.New("status"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(dataSourcesListName,
			tfjsonpath.New("streams"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(dataSourceHistoryName,
			tfjsonpath.New("stream_id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(dataSourceHistoryName,
			tfjsonpath.New("streams"), knownvalue.NotNull()),
	}
}

// TestAccMonitorLogsStream_lifecycle exercises resource CRUD and all monitor logs
// stream data sources against a single stream (one stream per account).
func TestAccMonitorLogsStream_lifecycle(t *testing.T) {
	acceptance.LongRunningTest(t)
	requireNoExistingStreams(t)

	label := acctest.RandomWithPrefix("tf-test")
	updatedLabel := label + "-updated"
	region := getRegionForStreamTest(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config:            tmpl.Lifecycle(t, label, region),
				ConfigStateChecks: lifecycleStateChecks(label),
				Check:             testAccCheckLogStreamReady(resourceName),
			},
			{
				// Only update the stream label; keep bucket/key/destination labels stable
				// to avoid regenerating temp object storage keys (sensitive drift).
				Config:            tmpl.Updates(t, label, region),
				ConfigStateChecks: lifecycleStateChecks(updatedLabel),
				Check:             testAccCheckLogStreamReady(resourceName),
			},
			{
				RefreshState: true,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status", "version"},
				Check:                   testAccCheckLogStreamReady(resourceName),
			},
		},
	})
}

func TestAccResourceMonitorLogsStream_lkeAuditLogs(t *testing.T) {
	acceptance.LongRunningTest(t)
	requireNoExistingStreams(t)

	label := acctest.RandomWithPrefix("tf-test")
	region := getRegionForStreamTest(t)
	clusterID := os.Getenv("LINODE_TEST_LKE_CLUSTER_ID")
	if clusterID == "" {
		t.Skip("skipping: LINODE_TEST_LKE_CLUSTER_ID must be set to an existing LKE cluster ID")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.LKEAuditLogs(t, label, clusterID, region),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("label"), knownvalue.StringExact(label)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("type"), knownvalue.StringExact("lke_audit_logs")),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("status"), knownvalue.NotNull()),
				},
				Check: testAccCheckLogStreamReady(resourceName),
			},
			{
				RefreshState: true,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status", "version"},
				Check:                   testAccCheckLogStreamReady(resourceName),
			},
		},
	})
}

func TestAccResourceMonitorLogsStream_invalidType(t *testing.T) {
	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      tmpl.InvalidType(t, label),
				ExpectError: regexp.MustCompile(`(?i)invalid.*type|type.*invalid`),
			},
		},
	})
}

func TestAccResourceMonitorLogsStream_invalidDestination(t *testing.T) {
	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      tmpl.InvalidDestination(t, label),
				ExpectError: regexp.MustCompile(`(?i)400|not found|invalid|destination`),
			},
		},
	})
}
