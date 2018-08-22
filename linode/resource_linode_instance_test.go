package linode

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/linode/linodego"
)

func TestAccLinodeInstance_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	var instanceName = acctest.RandomWithPrefix("tf_test")
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
				Config: testAccCheckLinodeInstanceBasic(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "image", "linode/ubuntu18.04"),
					resource.TestCheckResourceAttr(resName, "region", "us-east"),
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

func TestAccLinodeInstance_config(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	var instanceName = acctest.RandomWithPrefix("tf_test")
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
				Config: testAccCheckLinodeInstanceWithConfig(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "region", "us-east"),
					// resource.TestCheckResourceAttr(resName, "kernel", "linode/latest-64bit"),
					resource.TestCheckResourceAttr(resName, "group", "testing"),
					resource.TestCheckResourceAttr(resName, "swap_size", "0"),
					testAccCheckComputeInstanceConfig(&instance, "config", "linode/latest-64bit"),
				),
			},

			resource.TestStep{
				ResourceName: resName,
				ImportState:  true,
			},
		},
	})
}

func TestAccLinodeInstance_multipleConfigs(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	var instanceName = acctest.RandomWithPrefix("tf_test")
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
				Config: testAccCheckLinodeInstanceWithMultipleConfigs(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "region", "us-east"),
					// resource.TestCheckResourceAttr(resName, "kernel", "linode/latest-64bit"),
					resource.TestCheckResourceAttr(resName, "group", "testing"),
					resource.TestCheckResourceAttr(resName, "swap_size", "0"),
					testAccCheckComputeInstanceConfig(&instance, "configa", "linode/latest-64bit"),
					testAccCheckComputeInstanceConfig(&instance, "configb", "linode/latest-32bit"),
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
	var instance linodego.Instance
	var instanceName = acctest.RandomWithPrefix("tf_test")
	resName := "linode_instance.foobar"
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("linode@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Error generating test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeInstanceBasic(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					//resource.TestCheckResourceAttr("linode_instance.foobar", "group", "testing"),
				),
			},
			resource.TestStep{
				Config: testAccCheckLinodeInstanceConfigUpdates(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s_renamed", instanceName)),
					//resource.TestCheckResourceAttr("linode_instance.foobar", "group", "integration"),
				),
			},
		},
	})
}

func TestAccLinodeInstanceResize(t *testing.T) {
	t.Parallel()
	var instance linodego.Instance
	var instanceName = acctest.RandomWithPrefix("tf_test")
	resName := "linode_instance.foobar"
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("linode@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Error generating test SSH key pair: %s", err)
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
					testAccCheckLinodeInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "plan_storage_utilized", "25600"),
					resource.TestCheckResourceAttr(resName, "storage_utilized", "25600"),
					resource.TestCheckResourceAttr(resName, "storage", "25600"),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
				),
			},
			// Bump it to a 2048, but don't expand the disk
			resource.TestStep{
				Config: testAccCheckLinodeInstanceConfigUpsizeBigger(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "type", "g6-standard-1"),
					resource.TestCheckResourceAttr(resName, "plan_storage_utilized", "25600"),
					resource.TestCheckResourceAttr(resName, "storage_utilized", "25600"),
					resource.TestCheckResourceAttr(resName, "storage", "25600"),
				),
			},
			// Go back down to a 1024
			resource.TestStep{
				Config: testAccCheckLinodeInstanceConfigDownsize(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
				),
			},
		},
	})
}

func TestAccLinodeInstanceExpandDisk(t *testing.T) {
	t.Parallel()
	var instance linodego.Instance
	var instanceName = acctest.RandomWithPrefix("tf_test")
	resName := "linode_instance.foobar"
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("linode@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Error generating test SSH key pair: %s", err)
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
					testAccCheckLinodeInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "plan_storage_utilized", "25600"),
				),
			},
			// Bump it to a 2048, and expand the disk
			resource.TestStep{
				Config: testAccCheckLinodeInstanceConfigUpsizeExpandDisk(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "type", "g6-standard-1"),
					resource.TestCheckResourceAttr(resName, "plan_storage_utilized", "25600"),
				),
			},
		},
	})
}

func TestAccLinodeInstancePrivateNetworking(t *testing.T) {
	t.Parallel()
	var instance linodego.Instance
	var instanceName = acctest.RandomWithPrefix("tf_test")
	resName := "linode_instance.foobar"
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("linode@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Error generating test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeInstanceConfigPrivateNetworking(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeInstanceExists(resName, &instance),
					testAccCheckLinodeInstanceAttributesPrivateNetworking("linode_instance.foobar"),
					resource.TestCheckResourceAttr(resName, "private_networking", "true"),
				),
			},
		},
	})
}

func testAccCheckLinodeInstanceExists(name string, instance *linodego.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(linodego.Client)

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)

		found, err := client.GetInstance(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of Instance %s: %s", rs.Primary.Attributes["label"], err)
		}

		*instance = *found

		return nil
	}
}

func testAccCheckLinodeInstanceDestroy(s *terraform.State) error {
	client, ok := testAccProvider.Meta().(linodego.Client)
	if !ok {
		return fmt.Errorf("Error getting Linode client")
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_instance" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v as int", rs.Primary.ID)
		}

		if id == 0 {
			return fmt.Errorf("should not have Linode ID 0")
		}

		_, err = client.GetInstance(context.Background(), id)

		if err == nil {
			return fmt.Errorf("should not find Linode ID %d existing after delete", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error getting Linode ID %d: %s", id, err)
		}
	}

	return nil
}

func testAccCheckLinodeInstanceAttributesPrivateNetworking(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("should have found linode_instance resource %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("should have a Linode ID")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("should have an integer Linode ID: %s", err)
		}

		client, ok := testAccProvider.Meta().(linodego.Client)
		if !ok {
			return fmt.Errorf("should have a linodego.Client")
		}

		if err != nil {
			return err
		}

		instanceIPs, err := client.GetInstanceIPAddresses(context.Background(), id)
		if err != nil {
			return err
		}
		if len(instanceIPs.IPv4.Private) == 0 {
			return fmt.Errorf("should have a private ip on Linode ID %d", id)
		}
		return nil
	}
}

func testAccCheckComputeInstanceConfig(instance *linodego.Instance, label string, kernel string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(linodego.Client)

		if instance.ID == 0 {
			return fmt.Errorf("Error fetching configs, Instance ID is 0")
		}

		instanceConfigs, err := client.ListInstanceConfigs(context.Background(), instance.ID, nil)

		if err != nil {
			return fmt.Errorf("Error fetching configs: %s", err)
		}

		if len(instanceConfigs) == 0 {
			return fmt.Errorf("No configs")
		}

		for _, config := range instanceConfigs {
			if config.Label == label && config.Kernel == kernel {
				return nil
			}
		}

		return fmt.Errorf("Config not found: %s", label)
	}
}

func testAccCheckComputeInstanceDisk(instance *linodego.Instance, label string, size int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(linodego.Client)

		instanceDisks, err := client.ListInstanceDisks(context.Background(), instance.ID, nil)

		if err != nil {
			return fmt.Errorf("Error fetching disks: %s", err)
		}

		if len(instanceDisks) == 0 {
			return fmt.Errorf("No disks")
		}

		for _, disk := range instanceDisks {
			if disk.Label == label && disk.Size == size {
				return nil
			}
		}

		return fmt.Errorf("Disk not found: %s", label)
	}
}

func testAccCheckLinodeInstanceBasic(instance string, pubkey string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	type = "g6-nanode-1"
	image = "linode/ubuntu18.04"
	region = "us-east"
	root_pass = "terraform-test"
	swap_size = 256
	authorized_keys = "%s"
	group = "testing"
}`, instance, pubkey)
}

func testAccCheckLinodeInstanceWithConfig(instance string, pubkey string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	type = "g6-nanode-1"
	region = "us-east"
	config {
		label = "config"
		kernel = "linode/latest-64bit"
	}
	group = "testing"
}`, instance)
}

func testAccCheckLinodeInstanceWithMultipleConfigs(instance string, pubkey string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	type = "g6-nanode-1"
	region = "us-east"
	config {
		label = "configa"
		kernel = "linode/latest-64bit"
	}
	config {
		label = "configb"
		kernel = "linode/latest-32bit"
	}
	group = "testing"
}`, instance)
}

func testAccCheckLinodeInstanceWithDisk(instance string, pubkey string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	type = "g6-nanode-1"
	image = "linode/ubuntu18.04"
	region = "us-east"
	config {
		kernel = "linode/latest-64bit"
	}
	root_pass = "terraform-test"
	swap_size = 256
	authorized_keys = "%s"
	group = "testing"
}`, instance, pubkey)
}

func testAccCheckLinodeInstanceWithMultipleDisks(instance string, pubkey string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	type = "g6-nanode-1"
	image = "linode/ubuntu18.04"
	region = "us-east"
	config {
		kernel = "linode/latest-64bit"
	}
	root_pass = "terraform-test"
	swap_size = 256
	authorized_keys = "%s"
	group = "testing"
}`, instance, pubkey)
}

func testAccCheckLinodeInstanceWithDiskAndConfig(instance string, pubkey string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	type = "g6-nanode-1"
	image = "linode/ubuntu18.04"
	region = "us-east"
	config {
		kernel = "linode/latest-64bit"
	}
	root_pass = "terraform-test"
	swap_size = 256
	authorized_keys = "%s"
	group = "testing"
}`, instance, pubkey)
}

func testAccCheckLinodeInstanceWithMultipleDiskAndConfig(instance string, pubkey string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	type = "g6-nanode-1"
	image = "linode/ubuntu18.04"
	region = "us-east"
	config {
		kernel = "linode/latest-64bit"
	}
	root_pass = "terraform-test"
	swap_size = 256
	authorized_keys = "%s"
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
	authorized_keys = "%s"
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
	authorized_keys = "%s"
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
	authorized_keys = "%s"
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
	authorized_keys = "%s"
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
	authorized_keys = "%s"
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
	authorized_keys = "%s"
	group = "testing"
}`, instance, pubkey)
}
