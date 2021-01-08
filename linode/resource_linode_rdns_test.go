package linode

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
)

func init() {
	resource.AddTestSweepers("linode_rdns", &resource.Sweeper{
		Name: "linode_rdns",
		F:    testSweepLinodeRDNS,
	})
}

func testSweepLinodeRDNS(prefix string) error {
	client, err := getClientForSweepers()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	ips, err := client.ListIPAddresses(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("Error getting IPAddresses: %s", err)
	}
	updateOpts := linodego.IPAddressUpdateOptions{RDNS: nil}
	for _, ip := range ips {
		if !shouldSweepAcceptanceTestResource(prefix, ip.RDNS) {
			continue
		}
		_, err := client.UpdateIPAddress(context.Background(), ip.Address, updateOpts)

		if err != nil {
			return fmt.Errorf("Error clearing RDNS %s during sweep: %s", ip.RDNS, err)
		}
	}

	return nil
}

func TestAccLinodeRDNS_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_rdns.foobar"
	var linodeLabel = acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeRDNSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeRDNSBasic(linodeLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeRDNSExists,
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

func TestAccLinodeRDNS_update(t *testing.T) {
	t.Parallel()

	var label = acctest.RandomWithPrefix("tf_test")
	resName := "linode_rdns.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeRDNSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeRDNSBasic(label),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeRDNSExists,
					resource.TestCheckResourceAttrPair(resName, "address", "linode_instance.foobar", "ip_address"),
					resource.TestMatchResourceAttr(resName, "rdns", regexp.MustCompile(`([0-9]{1,3}\.){4}nip.io$`)),
				),
			},
			{
				Config: testAccCheckLinodeRDNSChanged(label),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeRDNSExists,
					resource.TestMatchResourceAttr(resName, "rdns", regexp.MustCompile(`([0-9]{1,3}\-){3}[0-9]{1,3}.nip.io$`)),
				),
			},
			{
				Config: testAccCheckLinodeRDNSDeleted(label),
			},
			{
				Config: testAccCheckLinodeRDNSDeleted(label),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.linode_networking_ip.foobar", "rdns", regexp.MustCompile(`.members.linode.com$`)),
				),
			},
		},
	})
}

func testAccCheckLinodeRDNSExists(s *terraform.State) error {
	client := testAccProvider.Meta().(linodego.Client)

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

func testAccCheckLinodeRDNSDestroy(s *terraform.State) error {
	client, ok := testAccProvider.Meta().(linodego.Client)
	if !ok {
		return fmt.Errorf("Error getting Linode client")
	}
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

func testAccCheckLinodeRDNSBasic(label string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	group = "tf_test"
	image = "linode/containerlinux"
	type = "g6-standard-1"
	region = "us-east"
}

resource "linode_rdns" "foobar" {
	address = "${linode_instance.foobar.ip_address}"
	rdns = "${linode_instance.foobar.ip_address}.nip.io"
}`, label)
}

func testAccCheckLinodeRDNSChanged(label string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	group = "tf_test"
	image = "linode/containerlinux"
	type = "g6-standard-1"
	region = "us-east"
}

resource "linode_rdns" "foobar" {
	rdns    = "${replace(linode_instance.foobar.ip_address, ".", "-")}.nip.io"
	address = "${linode_instance.foobar.ip_address}"
}
`, label)
}

func testAccCheckLinodeRDNSDeleted(label string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	group = "tf_test"
	image = "linode/containerlinux"
	type = "g6-standard-1"
	region = "us-east"
}

data "linode_networking_ip" "foobar" {
	address = "${linode_instance.foobar.ip_address}"
}
`, label)
}
