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
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/monitorlogsdestination/tmpl"
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
