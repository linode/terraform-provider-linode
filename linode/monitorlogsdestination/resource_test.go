//go:build integration || monitorlogsdestination

package monitorlogsdestination_test

import (
	"context"
	"fmt"
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
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
	"github.com/linode/terraform-provider-linode/v4/linode/monitorlogsdestination/tmpl"
)

func init() {
	resource.AddTestSweepers("linode_monitor_logs_destination", &resource.Sweeper{
		Name: "linode_monitor_logs_destination",
		F:    sweep,
	})
}

func sweep(prefix string) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	destinations, err := client.ListLogsDestinations(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("error listing logs destinations: %s", err)
	}

	for _, dest := range destinations {
		if !acceptance.ShouldSweep(prefix, dest.Label) {
			continue
		}
		if err := client.DeleteLogsDestination(context.Background(), dest.ID); err != nil {
			return fmt.Errorf("error destroying logs destination %s during sweep: %s", dest.Label, err)
		}
	}

	return nil
}

func TestAccResourceLogsDestination_basic(t *testing.T) {
	t.Parallel()

	endpoint, err := acceptance.GetRandomObjectStorageEndpoint()
	if err != nil {
		t.Fatal(err)
	}

	testCluster, err := acceptance.GetEndpointCluster(*endpoint)
	if err != nil {
		t.Fatal(err)
	}

	resName := "linode_monitor_logs_destination.foobar"
	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkLogsDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, endpoint.Region, testCluster),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("label"), knownvalue.StringExact(label)),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("type"), knownvalue.StringExact("akamai_object_storage")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("created"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("created_by"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName,
						tfjsonpath.New("akamai_object_storage_details").AtMapKey("bucket_name"),
						knownvalue.StringExact(label+"-bucket"),
					),
				},
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"akamai_object_storage_details.access_key_secret"},
			},
			// Wait for the backend to finish flushing logs and releasing object locks
			// before Terraform continues with bucket teardown
			{
				Config: tmpl.BucketOnly(t, label, endpoint.Region, testCluster),
				Check: func(_ *terraform.State) error {
					time.Sleep(60 * time.Second)
					return nil
				},
			},
		},
	})
}

func TestAccResourceLogsDestination_update(t *testing.T) {
	t.Parallel()

	endpoint, err := acceptance.GetRandomObjectStorageEndpoint()
	if err != nil {
		t.Fatal(err)
	}

	testCluster, err := acceptance.GetEndpointCluster(*endpoint)
	if err != nil {
		t.Fatal(err)
	}

	resName := "linode_monitor_logs_destination.foobar"
	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkLogsDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, endpoint.Region, testCluster),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("label"), knownvalue.StringExact(label)),
				},
			},
			{
				Config: tmpl.Updates(t, label, endpoint.Region, testCluster),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("label"), knownvalue.StringExact(label+"-updated")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("type"), knownvalue.StringExact("akamai_object_storage")),
				},
			},
			// Destroy only the destination first, then sleep to allow the Linode
			// backend to finish flushing in-flight log objects before the test
			// framework deletes the bucket in teardown.
			{
				Config: tmpl.BucketOnly(t, label, endpoint.Region, testCluster),
				Check: func(_ *terraform.State) error {
					time.Sleep(60 * time.Second)
					return nil
				},
			},
		},
	})
}

func TestAccResourceLogsDestination_invalidType(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      tmpl.InvalidType(t, label),
				ExpectError: regexp.MustCompile(`value must be one of`),
			},
		},
	})
}

func checkLogsDestinationDestroy(s *terraform.State) error {
	client := acceptance.TestAccSDKv2Provider.Meta().(*helper.ProviderMeta).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_monitor_logs_destination" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error parsing logs destination ID %v: %s", rs.Primary.ID, err)
		}

		_, err = client.GetLogsDestination(context.Background(), id)
		if err != nil {
			if !linodego.IsNotFound(err) {
				return fmt.Errorf("error getting logs destination with ID %d: %s", id, err)
			}
		} else {
			return fmt.Errorf("logs destination with ID %d still exists", id)
		}
	}

	return nil
}

const (
	trafficPeakEndpoint           = "TRAFFIC_PEAK_ENDPOINT_URL"
	defaultTrafficPeakEndpointURL = "https://example.com"
)

func trafficPeakEndpointURL() string {
	if v := os.Getenv(trafficPeakEndpoint); v != "" {
		return v
	}
	return defaultTrafficPeakEndpointURL
}

// TestAccResourceLogsDestination_trafficPeak covers create, update, and import
// against a single destination to minimize created resources.
func TestAccResourceLogsDestination_trafficPeak(t *testing.T) {
	t.Parallel()

	resName := "linode_monitor_logs_destination.foobar"
	label := acctest.RandomWithPrefix("tf-test")
	endpointURL := trafficPeakEndpointURL()
	username := acctest.RandomWithPrefix("tf-user")
	password := acctest.RandString(24)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkLogsDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.TrafficPeakBasic(t, label, endpointURL, username, password),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("label"), knownvalue.StringExact(label)),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("type"), knownvalue.StringExact("traffic_peak")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("created"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName,
						tfjsonpath.New("traffic_peak_details").AtMapKey("endpoint_url"),
						knownvalue.StringExact(endpointURL)),
					statecheck.ExpectKnownValue(resName,
						tfjsonpath.New("traffic_peak_details").AtMapKey("content_type"),
						knownvalue.StringExact("application/json")),
					statecheck.ExpectKnownValue(resName,
						tfjsonpath.New("traffic_peak_details").AtMapKey("data_compression"),
						knownvalue.StringExact("gzip")),
					statecheck.ExpectKnownValue(resName,
						tfjsonpath.New("traffic_peak_details").AtMapKey("authentication").AtMapKey("type"),
						knownvalue.StringExact("basic")),
					statecheck.ExpectKnownValue(resName,
						tfjsonpath.New("akamai_object_storage_details"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resName,
						tfjsonpath.New("custom_https_details"), knownvalue.Null()),
				},
			},
			{
				// authentication.type omitted: must remain "basic" (optional + computed)
				Config: tmpl.TrafficPeakUpdates(t, label, endpointURL, username, password),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("label"), knownvalue.StringExact(label+"-updated")),
					statecheck.ExpectKnownValue(resName,
						tfjsonpath.New("traffic_peak_details").AtMapKey("content_type"),
						knownvalue.StringExact("application/json; charset=utf-8")),
					statecheck.ExpectKnownValue(resName,
						tfjsonpath.New("traffic_peak_details").AtMapKey("authentication").AtMapKey("type"),
						knownvalue.StringExact("basic")),
				},
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"traffic_peak_details.authentication.username",
					"traffic_peak_details.authentication.password",
				},
			},
		},
	})
}

func TestAccResourceLogsDestination_trafficPeakInvalidAuthType(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      tmpl.TrafficPeakInvalidAuthType(t, acctest.RandomWithPrefix("tf-test")),
				ExpectError: regexp.MustCompile(`value must be one of`),
			},
		},
	})
}

func TestAccResourceLogsDestination_trafficPeakMissingAuth(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      tmpl.TrafficPeakMissingAuth(t, acctest.RandomWithPrefix("tf-test")),
				ExpectError: regexp.MustCompile(`(?s)"authentication" is required`),
			},
		},
	})
}

// content_type is required for traffic_peak destinations.
func TestAccResourceLogsDestination_trafficPeakMissingContentType(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      tmpl.TrafficPeakMissingContentType(t, acctest.RandomWithPrefix("tf-test")),
				ExpectError: regexp.MustCompile(`(?s)"content_type" is required`),
			},
		},
	})
}
