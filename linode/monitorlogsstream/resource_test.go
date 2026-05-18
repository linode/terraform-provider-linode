//go:build integration || monitorlogsstream

package monitorlogsstream_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/monitorlogsstream/tmpl"
)

const (
	resourceName = "linode_monitor_logs_stream.foobar"
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

// requireNoExistingStreams skips the test if a stream already exists on the
// account, since only one stream is allowed per account.
func requireNoExistingStreams(t *testing.T) {
	t.Helper()

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
		t.Skipf("existing stream(s) on account (IDs: %v); only one stream allowed per account", ids)
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

func TestAccResourceMonitorLogsStream_basic(t *testing.T) {
	acceptance.LongRunningTest(t)
	requireNoExistingStreams(t)

	label := acctest.RandomWithPrefix("tf-test")
	region := getRegionForStreamTest(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, region),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("label"), knownvalue.StringExact(label)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("type"), knownvalue.StringExact("audit_logs")),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
			{
				Config: tmpl.Updates(t, label, region),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("label"), knownvalue.StringExact(label+"-updated")),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("type"), knownvalue.StringExact("audit_logs")),
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
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
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
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
