//go:build integration || monitorlogsdestination

package monitorlogsdestination_test

import (
	"context"
	"fmt"
	"log"
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
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/monitorlogsdestination/tmpl"
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
		log.Fatal(fmt.Errorf("error getting client: %s", err))
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

	region, err := acceptance.GetRandomRegionWithCaps([]string{linodego.CapabilityObjectStorage}, "core")
	if err != nil {
		t.Skipf("could not get region with Object Storage: %s", err)
	}

	resName := "linode_monitor_logs_destination.foobar"
	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkLogsDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, region),
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
			{
				// Short sleep before the automatic post-test destroy to allow any
				// pending bucket activity to settle, preventing a "bucket not
				// empty" failure during obj_bucket_force_delete cleanup.
				PreConfig: func() { time.Sleep(15 * time.Second) },
				Config:    tmpl.Basic(t, label, region),
			},
		},
	})
}

func TestAccResourceLogsDestination_update(t *testing.T) {
	t.Parallel()

	region, err := acceptance.GetRandomRegionWithCaps([]string{linodego.CapabilityObjectStorage}, "core")
	if err != nil {
		t.Skipf("could not get region with Object Storage: %s", err)
	}

	resName := "linode_monitor_logs_destination.foobar"
	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkLogsDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, region),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("label"), knownvalue.StringExact(label)),
				},
			},
			{
				Config: tmpl.Updates(t, label, region),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("label"), knownvalue.StringExact(label+"-updated")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("type"), knownvalue.StringExact("akamai_object_storage")),
				},
			},
			{
				// Short sleep before the automatic post-test destroy to allow any
				// pending bucket activity to settle, preventing a "bucket not
				// empty" failure during obj_bucket_force_delete cleanup.
				PreConfig: func() { time.Sleep(15 * time.Second) },
				Config:    tmpl.Updates(t, label, region),
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
		if err == nil {
			return fmt.Errorf("logs destination with ID %d still exists", id)
		}
	}

	return nil
}
