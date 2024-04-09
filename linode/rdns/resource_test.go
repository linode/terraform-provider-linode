//go:build integration || rdns

package rdns_test

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/rdns/tmpl"
)

var testRegion string

func init() {
	resource.AddTestSweepers("linode_rdns", &resource.Sweeper{
		Name: "linode_rdns",
		F:    sweep,
	})

	region, err := acceptance.GetRandomRegionWithCaps([]string{"linodes"})
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func sweep(prefix string) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	ips, err := client.ListIPAddresses(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("Error getting IPAddresses: %s", err)
	}
	updateOpts := linodego.IPAddressUpdateOptions{RDNS: nil}
	for _, ip := range ips {
		if !acceptance.ShouldSweep(prefix, ip.RDNS) {
			continue
		}
		_, err := client.UpdateIPAddress(context.Background(), ip.Address, updateOpts)
		if err != nil {
			return fmt.Errorf("Error clearing RDNS %s during sweep: %s", ip.RDNS, err)
		}
	}

	return nil
}

func TestAccResourceRDNS_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_rdns.foobar"
	linodeLabel := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkRDNSDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, linodeLabel, testRegion, false),
				Check: resource.ComposeTestCheckFunc(
					checkRDNSExists,
					resource.TestMatchResourceAttr(resName, "rdns", regexp.MustCompile(`.nip.io$`)),
				),
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"wait_for_available"},
			},
		},
	})
}

func TestAccResourceRDNS_update(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf_test")
	resName := "linode_rdns.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkRDNSDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, testRegion, false),
				Check: resource.ComposeTestCheckFunc(
					checkRDNSExists,
					resource.TestCheckResourceAttrPair(resName, "address", "linode_instance.foobar", "ip_address"),
					resource.TestMatchResourceAttr(resName, "rdns", regexp.MustCompile(`([0-9]{1,3}\.){4}nip.io$`)),
				),
			},
			{
				Config: tmpl.Changed(t, label, testRegion, false),
				Check: resource.ComposeTestCheckFunc(
					checkRDNSExists,
					resource.TestMatchResourceAttr(resName, "rdns", regexp.MustCompile(`([0-9]{1,3}\-){3}[0-9]{1,3}.nip.io$`)),
				),
			},
			{
				Config: tmpl.Deleted(t, label, testRegion),
			},
			{
				Config: tmpl.Deleted(t, label, testRegion),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.linode_networking_ip.foobar", "rdns", regexp.MustCompile(`.ip.linodeusercontent.com$`)),
				),
			},
		},
	})
}

// This test case simply ensures a
func TestAccResourceRDNS_waitForAvailable(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf_test")
	resName := "linode_rdns.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkRDNSDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, testRegion, true),
				Check: resource.ComposeTestCheckFunc(
					checkRDNSExists,
					resource.TestCheckResourceAttrPair(resName, "address", "linode_instance.foobar", "ip_address"),
					resource.TestMatchResourceAttr(resName, "rdns", regexp.MustCompile(`([0-9]{1,3}\.){4}nip.io$`)),
				),
			},
			{
				Config: tmpl.Changed(t, label, testRegion, true),
				Check: resource.ComposeTestCheckFunc(
					checkRDNSExists,
					resource.TestMatchResourceAttr(resName, "rdns", regexp.MustCompile(`([0-9]{1,3}\-){3}[0-9]{1,3}.nip.io$`)),
				),
			},
			{
				Config: tmpl.Deleted(t, label, testRegion),
			},
			{
				Config: tmpl.Deleted(t, label, testRegion),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.linode_networking_ip.foobar", "rdns", regexp.MustCompile(`.ip.linodeusercontent.com$`)),
				),
			},
		},
	})
}

func TestAccResourceRDNS_waitForAvailableWithTimeout(t *testing.T) {
	t.Parallel()

	resName := "linode_rdns.foobar"
	linodeLabel := acctest.RandomWithPrefix("tf_test")

	createTimeout := "15m"
	updateTimeout := "15m"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkRDNSDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.WithTimeout(t, linodeLabel, testRegion, createTimeout, updateTimeout),
				Check: resource.ComposeTestCheckFunc(
					checkRDNSExists,
					resource.TestMatchResourceAttr(resName, "rdns", regexp.MustCompile(`.nip.io$`)),
				),
			},
			{
				Config: tmpl.WithTimeoutUpdated(t, linodeLabel, testRegion, createTimeout, updateTimeout),
				Check: resource.ComposeTestCheckFunc(
					checkRDNSExists,
					resource.TestMatchResourceAttr(resName, "rdns", regexp.MustCompile(`([0-9]{1,3}\-){3}[0-9]{1,3}.nip.io$`)),
				),
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"wait_for_available"},
			},
		},
	})
}

func checkRDNSExists(s *terraform.State) error {
	client := acceptance.TestAccFrameworkProvider.Meta.Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_rdns" {
			continue
		}

		_, err := client.GetIPAddress(context.Background(), rs.Primary.Attributes["address"])
		if err != nil {
			return fmt.Errorf("Error retrieving state of RDNS %s: %s", rs.Primary.Attributes["rdns"], err)
		}
	}

	return nil
}

func checkRDNSDestroy(s *terraform.State) error {
	client := acceptance.TestAccFrameworkProvider.Meta.Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_rdns" {
			continue
		}

		id := rs.Primary.ID
		ip, err := client.GetIPAddress(context.Background(), id)
		if err != nil {
			if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code == 404 {
				return nil
			}

			if ip.RDNS[len(ip.RDNS)-len("ip.linodeusercontent.com"):] == "ip.linodeusercontent.com" {
				return nil
			}

			return fmt.Errorf("Linode RDNS with IP %s still exists: %s", id, err)
		}
	}

	return nil
}
