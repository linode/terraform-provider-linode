package ipv6range_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/ipv6range/tmpl"
	"regexp"
	"testing"
)

func TestAccIPv6Range_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_ipv6_range.foobar"
	instLabel := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkIPv6RangeDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, instLabel),
				Check: resource.ComposeTestCheckFunc(
					checkIPv6RangeExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "prefix_length", "64"),
					resource.TestCheckResourceAttr(resName, "is_bgp", "false"),
					resource.TestCheckResourceAttr(resName, "region", "us-southeast"),

					resource.TestCheckResourceAttrSet(resName, "range"),
					resource.TestCheckResourceAttrSet(resName, "linode_id"),
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
}

func TestAccIPv6Range_routeTarget(t *testing.T) {
	t.Parallel()

	resName := "linode_ipv6_range.foobar"
	instLabel := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkIPv6RangeDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.RouteTarget(t, instLabel),
				Check: resource.ComposeTestCheckFunc(
					checkIPv6RangeExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "prefix_length", "64"),
					resource.TestCheckResourceAttr(resName, "is_bgp", "false"),
					resource.TestCheckResourceAttr(resName, "region", "us-southeast"),

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
}

func TestAccIPv6Range_noID(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkIPv6RangeDestroy,
		Steps: []resource.TestStep{
			{
				Config:      tmpl.NoID(t),
				ExpectError: regexp.MustCompile("either linode_id or route_target must be specified"),
			},
		},
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
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_ipv6_range" {
			continue
		}

		_, err := client.GetIPv6Range(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Linode IPv6 range with id %s still exists", rs.Primary.ID)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 && apiErr.Code != 405 {
			return fmt.Errorf("error requesting Linode IPv6 range with id %s: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
