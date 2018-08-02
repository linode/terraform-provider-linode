package linode

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/chiefy/linodego"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccLinodeInstanceBasic(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"

	var instanceName = acctest.RandomWithPrefix("tf_test_")
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("linode@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeInstanceConfigBasic(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					//testAccCheckLinodeInstanceExists(&instance),

					testAccCheckLinodeInstanceExists,
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "image", "linode/ubuntu18.04"),
					resource.TestCheckResourceAttr(resName, "region", "us-east"),
					resource.TestCheckResourceAttr(resName, "kernel", "linode/latest-64bit"),
					resource.TestCheckResourceAttr(resName, "group", "testing"),
					resource.TestCheckResourceAttr(resName, "swap_size", "256"),
				),
			},

			resource.TestStep{
				ResourceName: resName,
				ImportState:  true,
			},
		},
	})
}

func TestAccLinodeInstanceUpdate(t *testing.T) {
	t.Parallel()

	var instanceName = acctest.RandomWithPrefix("tf_test_")
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("linode@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeInstanceConfigBasic(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeInstanceExists,
					resource.TestCheckResourceAttr("linode_instance.foobar", "label", instanceName),
					//resource.TestCheckResourceAttr("linode_instance.foobar", "group", "testing"),
				),
			},
			resource.TestStep{
				Config: testAccCheckLinodeInstanceConfigUpdates(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeInstanceExists,
					resource.TestCheckResourceAttr("linode_instance.foobar", "label", fmt.Sprintf("%s_renamed", instanceName)),
					//resource.TestCheckResourceAttr("linode_instance.foobar", "group", "integration"),
				),
			},
		},
	})
}

func TestAccLinodeInstanceResize(t *testing.T) {
	t.Parallel()

	var instanceName = acctest.RandomWithPrefix("tf_test_")
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("linode@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeInstanceDestroy,
		Steps: []resource.TestStep{
			// Start off with a Linode 1024
			resource.TestStep{
				Config: testAccCheckLinodeInstanceConfigUpsizeSmall(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeInstanceExists,
					resource.TestCheckResourceAttr("linode_instance.foobar", "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr("linode_instance.foobar", "plan_storage_utilized", "25600"),
					resource.TestCheckResourceAttr("linode_instance.foobar", "storage_utilized", "25600"),
					resource.TestCheckResourceAttr("linode_instance.foobar", "storage", "25600"),
				),
			},
			// Bump it to a 2048, but don't expand the disk
			resource.TestStep{
				Config: testAccCheckLinodeInstanceConfigUpsizeBigger(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeInstanceExists,
					resource.TestCheckResourceAttr("linode_instance.foobar", "type", "g6-standard-1"),
					resource.TestCheckResourceAttr("linode_instance.foobar", "plan_storage_utilized", "25600"),
					resource.TestCheckResourceAttr("linode_instance.foobar", "storage_utilized", "25600"),
					resource.TestCheckResourceAttr("linode_instance.foobar", "storage", "25600"),
				),
			},
			// Go back down to a 1024
			resource.TestStep{
				Config: testAccCheckLinodeInstanceConfigDownsize(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeInstanceExists,
					resource.TestCheckResourceAttr("linode_instance.foobar", "type", "g6-nanode-1"),
				),
			},
		},
	})
}

func TestAccLinodeInstanceExpandDisk(t *testing.T) {
	t.Parallel()

	var instanceName = acctest.RandomWithPrefix("tf_test_")
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("linode@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeInstanceDestroy,
		Steps: []resource.TestStep{
			// Start off with a Linode 1024
			resource.TestStep{
				Config: testAccCheckLinodeInstanceConfigUpsizeSmall(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeInstanceExists,
					resource.TestCheckResourceAttr("linode_instance.foobar", "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr("linode_instance.foobar", "plan_storage_utilized", "25600"),
				),
			},
			// Bump it to a 2048, and expand the disk
			resource.TestStep{
				Config: testAccCheckLinodeInstanceConfigUpsizeExpandDisk(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeInstanceExists,
					resource.TestCheckResourceAttr("linode_instance.foobar", "type", "g6-standard-1"),
					resource.TestCheckResourceAttr("linode_instance.foobar", "plan_storage_utilized", "25600"),
				),
			},
		},
	})
}

func TestAccLinodeInstancePrivateNetworking(t *testing.T) {
	t.Parallel()

	var instanceName = acctest.RandomWithPrefix("tf_test_")
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("linode@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeInstanceConfigPrivateNetworking(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeInstanceExists,
					testAccCheckLinodeInstanceAttributesPrivateNetworking("linode_instance.foobar"),
					resource.TestCheckResourceAttr("linode_instance.foobar", "private_networking", "true"),
				),
			},
		},
	})
}

func testAccCheckLinodeInstanceExists(s *terraform.State) error {
	client := testAccProvider.Meta().(linodego.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_instance" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)

		_, err = client.GetInstance(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of Instance %s: %s", rs.Primary.Attributes["label"], err)
		}
	}

	return nil
}

func testAccCheckLinodeInstanceDestroy(s *terraform.State) error {
	client, ok := testAccProvider.Meta().(linodego.Client)
	if !ok {
		return fmt.Errorf("Failed to get Linode client")
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_instance" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Failed parsing %v to int", rs.Primary.ID)
		}
		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)

		}

		_, err = client.GetInstance(context.Background(), id)

		if err == nil {
			return fmt.Errorf("Linode with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Failed to request Linode with id %d", id)
		}
	}

	return nil
}

func testAccCheckLinodeInstanceAttributesPrivateNetworking(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found %+v", rs)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Linode id set")
		}

		client := testAccProvider.Meta().(linodego.Client)
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			panic(err)
		}
		instanceIPs, err := client.GetInstanceIPAddresses(context.Background(), id)
		if err != nil {
			return err
		}
		if len(instanceIPs.IPv4.Private) == 0 {
			return fmt.Errorf("Private Ip is not set")
		}
		return nil
	}
}

func testAccCheckLinodeInstanceConfigBasic(instance string, pubkey string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	type = "g6-nanode-1"
	image = "linode/ubuntu18.04"
	region = "us-east"
	kernel = "linode/latest-64bit"
	root_password = "terraform-test"
	swap_size = 256
	ssh_key = "%s"
	group = "testing"
}`, instance, pubkey)
}

func testAccCheckLinodeInstanceConfigUpdates(instance string, pubkey string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s_renamed"
	type = "g6-nanode-1"
	image = "linode/ubuntu18.04"
	region = "us-east"
	kernel = "linode/latest-64bit"
	root_password = "terraform-test"
	swap_size = 256
	ssh_key = "%s"
	group = "testing"
}`, instance, pubkey)
}

func testAccCheckLinodeInstanceConfigUpsizeSmall(instance string, pubkey string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	type = "g6-nanode-1"
	image = "linode/ubuntu18.04"
	region = "us-east"
	kernel = "linode/latest-64bit"
	root_password = "terraform-test"
	swap_size = 256
	ssh_key = "%s"
	group = "testing"
}`, instance, pubkey)
}

func testAccCheckLinodeInstanceConfigUpsizeBigger(instance string, pubkey string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s_upsized"
	type = "g6-standard-1"
	image = "linode/ubuntu18.04"
	region = "us-east"
	kernel = "linode/latest-64bit"
	root_password = "terraform-test"
	swap_size = 256
	ssh_key = "%s"
	group = "testing"
}`, instance, pubkey)
}

func testAccCheckLinodeInstanceConfigDownsize(instance string, pubkey string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s_downsized"
	type = "g6-nanode-1"
	image = "linode/ubuntu18.04"
	region = "us-east"
	kernel = "linode/latest-64bit"
	root_password = "terraform-test"
	swap_size = 256
	ssh_key = "%s"
	group = "testing"
}`, instance, pubkey)
}

func testAccCheckLinodeInstanceConfigUpsizeExpandDisk(instance string, pubkey string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s_expanded"
	type = "g6-standard-1"
	disk_expansion = true
	image = "linode/ubuntu18.04"
	region = "us-east"
	kernel = "linode/latest-64bit"
	root_password = "terraform-test"
	swap_size = 256
	ssh_key = "%s"
	group = "testing"
}`, instance, pubkey)
}

func testAccCheckLinodeInstanceConfigPrivateNetworking(instance string, pubkey string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	type = "g6-nanode-1"
	image = "linode/ubuntu18.04"
	region = "us-east"
	kernel = "linode/latest-64bit"
	root_password = "terraform-test"
	swap_size = 256
	private_networking = true
	ssh_key = "%s"
	group = "testing"
}`, instance, pubkey)
}
