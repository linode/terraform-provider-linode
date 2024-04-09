//go:build integration || ipv6range

package ipv6range_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/ipv6range/tmpl"
)

// TODO: don't hardcode this once IPv6 sharing has a proper capability string
const testRegion = "eu-central"

func TestAccIPv6Range_basic(t *testing.T) {
	t.Parallel()

	acceptance.RunTestRetry(t, 3, func(retryT *acceptance.TRetry) {
		resName := "linode_ipv6_range.foobar"
		instLabel := acctest.RandomWithPrefix("tf_test")

		resource.Test(retryT, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
			CheckDestroy:             checkIPv6RangeDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.Basic(t, instLabel, testRegion),
					Check: resource.ComposeTestCheckFunc(
						checkIPv6RangeExists(resName, nil),
						resource.TestCheckResourceAttr(resName, "prefix_length", "64"),
						resource.TestCheckResourceAttr(resName, "is_bgp", "false"),
						resource.TestCheckResourceAttr(resName, "region", testRegion),

						resource.TestCheckResourceAttrSet(resName, "range"),
						resource.TestCheckResourceAttrSet(resName, "linode_id"),
						resource.TestCheckResourceAttrSet(resName, "linodes.0"),
						resource.TestCheckResourceAttrSet(resName, "route_target"),
					),
				},
				{
					ResourceName:            resName,
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"linode_id", "route_target"},
				},
			},
		})
	})
}

func TestAccIPv6Range_routeTarget(t *testing.T) {
	t.Parallel()

	acceptance.RunTestRetry(t, 3, func(retryT *acceptance.TRetry) {
		resName := "linode_ipv6_range.foobar"
		instLabel := acctest.RandomWithPrefix("tf_test")

		resource.Test(retryT, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
			CheckDestroy:             checkIPv6RangeDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.RouteTarget(t, instLabel, testRegion),
					Check: resource.ComposeTestCheckFunc(
						checkIPv6RangeExists(resName, nil),
						resource.TestCheckResourceAttr(resName, "prefix_length", "64"),
						resource.TestCheckResourceAttr(resName, "is_bgp", "false"),
						resource.TestCheckResourceAttr(resName, "region", testRegion),

						resource.TestCheckResourceAttrSet(resName, "range"),
						resource.TestCheckResourceAttrSet(resName, "route_target"),
						resource.TestCheckResourceAttrSet(resName, "linodes.0"),
					),
				},
				{
					ResourceName:            resName,
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"linode_id", "route_target"},
				},
			},
		})
	})
}

func TestAccIPv6Range_noID(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkIPv6RangeDestroy,
		Steps: []resource.TestStep{
			{
				Config:      tmpl.NoID(t),
				ExpectError: regexp.MustCompile("Either linode_id or route_target must be specified."),
			},
		},
	})
}

func TestAccIPv6Range_reassignment(t *testing.T) {
	t.Parallel()

	acceptance.RunTestRetry(t, 3, func(retryT *acceptance.TRetry) {
		resName := "linode_ipv6_range.foobar"
		instance1ResName := "linode_instance.foobar"
		instance2ResName := "linode_instance.foobar2"

		instLabel := acctest.RandomWithPrefix("tf_test")

		var instance1 linodego.Instance
		var instance2 linodego.Instance

		resource.Test(retryT, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
			CheckDestroy:             checkIPv6RangeDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.ReassignmentStep1(t, instLabel, testRegion),
					Check: resource.ComposeTestCheckFunc(
						checkIPv6RangeExists(resName, nil),
						acceptance.CheckInstanceExists(instance1ResName, &instance1),
						acceptance.CheckInstanceExists(instance2ResName, &instance2),

						resource.TestCheckResourceAttr(resName, "prefix_length", "64"),
						resource.TestCheckResourceAttr(resName, "is_bgp", "false"),
						resource.TestCheckResourceAttr(resName, "region", testRegion),

						resource.TestCheckResourceAttrSet(resName, "range"),
						resource.TestCheckResourceAttrSet(resName, "linode_id"),
						resource.TestCheckResourceAttrSet(resName, "linodes.0"),
						resource.TestCheckResourceAttrSet(resName, "route_target"),
					),
				},
				{
					PreConfig: func() {
						validateInstanceIPv6Assignments(t, instance1.ID, instance2.ID)
					},
					Config: tmpl.ReassignmentStep2(t, instLabel, testRegion),
					Check: resource.ComposeTestCheckFunc(
						checkIPv6RangeExists(resName, nil),
						resource.TestCheckResourceAttr(resName, "prefix_length", "64"),
						resource.TestCheckResourceAttr(resName, "is_bgp", "false"),
						resource.TestCheckResourceAttr(resName, "region", testRegion),

						resource.TestCheckResourceAttrSet(resName, "range"),
						resource.TestCheckResourceAttrSet(resName, "linode_id"),
						resource.TestCheckResourceAttrSet(resName, "linodes.0"),
						resource.TestCheckResourceAttrSet(resName, "route_target"),
					),
				},
				{
					Config: tmpl.ReassignmentStep2(t, instLabel, testRegion),
					PreConfig: func() {
						validateInstanceIPv6Assignments(t, instance2.ID, instance1.ID)
					},
				},
			},
		})
	})
}

func TestAccIPv6Range_raceCondition(t *testing.T) {
	t.Parallel()

	// Occasionally IPv6 range deletions take a bit to replicate
	acceptance.RunTestRetry(t, 3, func(retryT *acceptance.TRetry) {
		instLabel := acctest.RandomWithPrefix("tf_test")

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
			CheckDestroy:             checkIPv6RangeDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.RaceCondition(t, instLabel, testRegion),
					Check:  checkIPv6RangeNoDuplicates,
				},
			},
		})
	})
}

func checkIPv6RangeExists(name string, ipv6Range *linodego.IPv6Range) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		found, err := client.GetIPv6Range(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to retrieve state of ipv6 range %s: %s", rs.Primary.Attributes["range"], err)
		}

		if ipv6Range != nil {
			*ipv6Range = *found
		}

		return nil
	}
}

func checkIPv6RangeDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

	// We should retry here as there is sometimes a delay between deletion request and
	// range deletion. This should significantly reduce the number of intermittent cleanup
	// failures we get.
	err := resource.RetryContext(context.Background(), 30*time.Second, func() *resource.RetryError {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "linode_ipv6_range" {
				continue
			}

			_, err := client.GetIPv6Range(context.Background(), rs.Primary.ID)
			if err == nil {
				return resource.RetryableError(fmt.Errorf("ipv6 range still exists: %s", err))
			}

			if apiErr, ok := err.(*linodego.Error); ok &&
				// Intermittent error codes
				apiErr.Code != 403 && apiErr.Code != 404 && apiErr.Code != 405 {
				return resource.NonRetryableError(
					fmt.Errorf("error requesting ipv6 range with id %s: %s", rs.Primary.ID, err))
			}
		}

		return nil
	})

	return err
}

func checkIPv6RangeNoDuplicates(s *terraform.State) error {
	existingRanges := make(map[string]bool)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_ipv6_range" {
			continue
		}

		if _, ok := existingRanges[rs.Primary.ID]; ok {
			return fmt.Errorf("duplicate range found: %s", rs.Primary.ID)
		}

		existingRanges[rs.Primary.ID] = true
	}

	return nil
}

func validateInstanceIPv6Assignments(t *testing.T, assignedID, unassignedID int) {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

	assignedNetworking, err := client.GetInstanceIPAddresses(context.Background(), assignedID)
	if err != nil {
		t.Fatal(err)
	}

	unassignedNetworking, err := client.GetInstanceIPAddresses(context.Background(), unassignedID)
	if err != nil {
		t.Fatal(err)
	}

	if len(unassignedNetworking.IPv6.Global) > 0 {
		t.Fatalf("expected instance to have no attached ipv6 ranged, got %d", len(unassignedNetworking.IPv6.Global))
	}

	if len(assignedNetworking.IPv6.Global) < 1 {
		t.Fatalf("expected instance to have one attached ipv6 ranged, got %d", len(assignedNetworking.IPv6.Global))
	}
}
