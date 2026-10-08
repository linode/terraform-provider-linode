//go:build integration || monitorlogsdestination

package monitorlogsdestination_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/monitorlogsdestination/tmpl"
)

func TestAccDataSourceLogsDestination_basic(t *testing.T) {
	t.Parallel()

	endpoint, err := acceptance.GetRandomObjectStorageEndpoint()
	if err != nil {
		t.Fatal(err)
	}

	testCluster, err := acceptance.GetEndpointCluster(*endpoint)
	if err != nil {
		t.Fatal(err)
	}

	dataName := "data.linode_monitor_logs_destination.foobar"
	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkLogsDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, label, endpoint.Region, testCluster),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dataName, tfjsonpath.New("label"), knownvalue.StringExact(label)),
					statecheck.ExpectKnownValue(dataName, tfjsonpath.New("type"), knownvalue.StringExact("akamai_object_storage")),
					statecheck.ExpectKnownValue(dataName, tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataName, tfjsonpath.New("created"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataName, tfjsonpath.New("updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataName, tfjsonpath.New("created_by"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataName,
						tfjsonpath.New("details").AtMapKey("bucket_name"),
						knownvalue.StringExact(label+"-bucket"),
					),
				},
			},
			{
				Config:      tmpl.DataNotFound(t),
				ExpectError: regexp.MustCompile(`\[404\]`),
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

// TestAccDataSourceLogsDestination_trafficPeak verifies the data source flattens a
// traffic_peak destination. Both steps reuse the same label so only one
// destination is ever created.
func TestAccDataSourceLogsDestination_trafficPeak(t *testing.T) {
	t.Parallel()

	dataName := "data.linode_monitor_logs_destination.foobar"
	label := acctest.RandomWithPrefix("tf-test")
	endpointURL := trafficPeakEndpointURL()
	username := acctest.RandomWithPrefix("tf-user")
	password := acctest.RandString(24)

	trafficPeakChecks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(dataName, tfjsonpath.New("label"), knownvalue.StringExact(label)),
		statecheck.ExpectKnownValue(dataName, tfjsonpath.New("type"), knownvalue.StringExact("traffic_peak")),
		statecheck.ExpectKnownValue(dataName, tfjsonpath.New("status"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(dataName, tfjsonpath.New("created"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(dataName, tfjsonpath.New("updated"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(dataName, tfjsonpath.New("created_by"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(dataName,
			tfjsonpath.New("details").AtMapKey("endpoint_url"), knownvalue.StringExact(endpointURL)),
		statecheck.ExpectKnownValue(dataName,
			tfjsonpath.New("details").AtMapKey("content_type"), knownvalue.StringExact("application/json")),
		statecheck.ExpectKnownValue(dataName,
			tfjsonpath.New("details").AtMapKey("data_compression"), knownvalue.StringExact("gzip")),
		// The API always defaults authentication.type to "basic".
		statecheck.ExpectKnownValue(dataName,
			tfjsonpath.New("details").AtMapKey("authentication_type"), knownvalue.StringExact("basic")),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkLogsDestinationDestroy,
		Steps: []resource.TestStep{
			{
				// authentication.type explicitly set to the only accepted value.
				Config:            tmpl.TrafficPeakDataBasic(t, label, endpointURL, username, password),
				ConfigStateChecks: trafficPeakChecks,
			},
			{
				// authentication.type omitted: the API still returns "basic".
				Config:            tmpl.TrafficPeakDataNoAuthType(t, label, endpointURL, username, password),
				ConfigStateChecks: trafficPeakChecks,
			},
		},
	})
}
