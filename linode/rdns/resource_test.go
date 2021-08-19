package rdns_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func init() {
	resource.AddTestSweepers("linode_rdns", &resource.Sweeper{
		Name: "linode_rdns",
		F:    sweep,
	})
}

func sweep(prefix string) error {
	client, err := acceptance.GetClientForSweepers()
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
	var linodeLabel = acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkRDNSDestroy,
		Steps: []resource.TestStep{
			{
				Config: resouceConfigBasic(linodeLabel),
				Check: resource.ComposeTestCheckFunc(
					checkRDNSExists,
					resource.TestMatchResourceAttr(resName, "rdns", regexp.MustCompile(`.nip.io$`)),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceRDNS_update(t *testing.T) {
	t.Parallel()

	var label = acctest.RandomWithPrefix("tf_test")
	resName := "linode_rdns.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkRDNSDestroy,
		Steps: []resource.TestStep{
			{
				Config: resouceConfigBasic(label),
				Check: resource.ComposeTestCheckFunc(
					checkRDNSExists,
					resource.TestCheckResourceAttrPair(resName, "address", "linode_instance.foobar", "ip_address"),
					resource.TestMatchResourceAttr(resName, "rdns", regexp.MustCompile(`([0-9]{1,3}\.){4}nip.io$`)),
				),
			},
			{
				Config: configResourceChanged(label),
				Check: resource.ComposeTestCheckFunc(
					checkRDNSExists,
					resource.TestMatchResourceAttr(resName, "rdns", regexp.MustCompile(`([0-9]{1,3}\-){3}[0-9]{1,3}.nip.io$`)),
				),
			},
			{
				Config: configResourceDeleted(label),
			},
			{
				Config: configResourceDeleted(label),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.linode_networking_ip.foobar", "rdns", regexp.MustCompile(`.members.linode.com$`)),
				),
			},
		},
	})
}

func checkRDNSExists(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

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
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
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

			if ip.RDNS[len(ip.RDNS)-len("members.linode.com"):] == "members.linode.com" {
				return nil
			}

			return fmt.Errorf("Linode RDNS with IP %s still exists: %s", id, err)
		}
	}

	return nil
}

func resouceConfigBasic(label string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	group = "tf_test"
	image = "linode/alpine3.12"
	type = "g6-standard-1"
	region = "us-east"
}

resource "linode_rdns" "foobar" {
	address = "${linode_instance.foobar.ip_address}"
	rdns = "${linode_instance.foobar.ip_address}.nip.io"
}`, label)
}

func configResourceChanged(label string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	group = "tf_test"
	image = "linode/alpine3.12"
	type = "g6-standard-1"
	region = "us-east"
}

resource "linode_rdns" "foobar" {
	rdns    = "${replace(linode_instance.foobar.ip_address, ".", "-")}.nip.io"
	address = "${linode_instance.foobar.ip_address}"
}
`, label)
}

func configResourceDeleted(label string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	group = "tf_test"
	image = "linode/alpine3.12"
	type = "g6-standard-1"
	region = "us-east"
}

data "linode_networking_ip" "foobar" {
	address = "${linode_instance.foobar.ip_address}"
}
`, label)
}
