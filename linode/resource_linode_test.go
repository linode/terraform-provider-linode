package linode

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/taoh/linodego"
)

func TestAccLinodeLinodeBasic(t *testing.T) {
	t.Parallel()

	var instance linodego.Linode
	var instanceName = fmt.Sprintf("tf_test_%s", acctest.RandString(10))
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("linode@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeLinodeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeLinodeConfigBasic(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeLinodeExists("linode_linode.foobar", &instance),
					resource.TestCheckResourceAttr("linode_linode.foobar", "name", instanceName),
					resource.TestCheckResourceAttr("linode_linode.foobar", "size", "1024"),
					resource.TestCheckResourceAttr("linode_linode.foobar", "image", "Ubuntu 16.04 LTS"),
					resource.TestCheckResourceAttr("linode_linode.foobar", "region", "Dallas, TX, USA"),
					resource.TestCheckResourceAttr("linode_linode.foobar", "kernel", "Latest 64 bit"),
					resource.TestCheckResourceAttr("linode_linode.foobar", "group", "testing"),
					resource.TestCheckResourceAttr("linode_linode.foobar", "swap_size", "256"),
				),
			},
			resource.TestStep{
				ResourceName:  "linode_linode.foobar",
				ImportState:   true,
				ImportStateId: fmt.Sprintf("%d", instance.LinodeId),
				// ImportStateId: getId(instance),
			},
		},
	})
}

/* @TODO I'm not getting the id of the instance. why not?
func getId(instance linodego.Linode) string {
	fmt.Printf("What did you do? %+v\n", instance)
	return fmt.Sprintf("%d", instance.LinodeId)
}
*/
func TestAccLinodeLinodeUpdate(t *testing.T) {
	t.Parallel()

	var instance linodego.Linode
	var instanceName = fmt.Sprintf("tf_test_%s", acctest.RandString(10))
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("linode@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeLinodeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeLinodeConfigBasic(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeLinodeExists("linode_linode.foobar", &instance),
					resource.TestCheckResourceAttr("linode_linode.foobar", "name", instanceName),
					resource.TestCheckResourceAttr("linode_linode.foobar", "group", "testing"),
				),
			},
			resource.TestStep{
				Config: testAccCheckLinodeLinodeConfigUpdates(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeLinodeExists("linode_linode.foobar", &instance),
					resource.TestCheckResourceAttr("linode_linode.foobar", "name", fmt.Sprintf("%s_renamed", instanceName)),
					resource.TestCheckResourceAttr("linode_linode.foobar", "group", "integration"),
				),
			},
		},
	})
}

func TestAccLinodeLinodeResize(t *testing.T) {
	t.Parallel()

	var instance linodego.Linode
	var instanceName = fmt.Sprintf("tf_test_%s", acctest.RandString(10))
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("linode@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeLinodeDestroy,
		Steps: []resource.TestStep{
			// Start off with a Linode 1024
			resource.TestStep{
				Config: testAccCheckLinodeLinodeConfigUpsizeSmall(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeLinodeExists("linode_linode.foobar", &instance),
					resource.TestCheckResourceAttr("linode_linode.foobar", "size", "1024"),
					resource.TestCheckResourceAttr("linode_linode.foobar", "plan_storage_utilized", "20480"),
				),
			},
			// Bump it to a 2048, but don't expand the disk
			resource.TestStep{
				Config: testAccCheckLinodeLinodeConfigUpsizeBigger(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeLinodeExists("linode_linode.foobar", &instance),
					resource.TestCheckResourceAttr("linode_linode.foobar", "size", "2048"),
					resource.TestCheckResourceAttr("linode_linode.foobar", "plan_storage_utilized", "20480"),
				),
			},
			// Go back down to a 1024
			resource.TestStep{
				Config: testAccCheckLinodeLinodeConfigDownsize(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeLinodeExists("linode_linode.foobar", &instance),
					resource.TestCheckResourceAttr("linode_linode.foobar", "size", "1024"),
				),
			},
		},
	})
}

func TestAccLinodeLinodeExpandDisk(t *testing.T) {
	t.Parallel()

	var instance linodego.Linode
	var instanceName = fmt.Sprintf("tf_test_%s", acctest.RandString(10))
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("linode@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeLinodeDestroy,
		Steps: []resource.TestStep{
			// Start off with a Linode 1024
			resource.TestStep{
				Config: testAccCheckLinodeLinodeConfigUpsizeSmall(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeLinodeExists("linode_linode.foobar", &instance),
					resource.TestCheckResourceAttr("linode_linode.foobar", "size", "1024"),
					resource.TestCheckResourceAttr("linode_linode.foobar", "plan_storage_utilized", "20480"),
				),
			},
			// Bump it to a 2048, and expand the disk
			resource.TestStep{
				Config: testAccCheckLinodeLinodeConfigUpsizeExpandDisk(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeLinodeExists("linode_linode.foobar", &instance),
					resource.TestCheckResourceAttr("linode_linode.foobar", "size", "2048"),
					resource.TestCheckResourceAttr("linode_linode.foobar", "plan_storage_utilized", "20480"),
				),
			},
		},
	})
}

func TestAccLinodeLinodePrivateNetworking(t *testing.T) {
	t.Parallel()

	var instance linodego.Linode
	var instanceName = fmt.Sprintf("tf_test_%s", acctest.RandString(10))
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("linode@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeLinodeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeLinodeConfigPrivateNetworking(instanceName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeLinodeExists("linode_linode.foobar", &instance),
					testAccCheckLinodeLinodeAttributesPrivateNetworking("linode_linode.foobar"),
					resource.TestCheckResourceAttr("linode_linode.foobar", "private_networking", "true"),
				),
			},
		},
	})
}

func testAccCheckLinodeLinodeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*linodego.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_linode" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Failed parsing %v to int", rs.Primary.ID)
		}

		linodes, err := client.Linode.List(id)

		if err != nil {
			return fmt.Errorf("Failed to get Linode list with %d id", id)
		} else if len(linodes.Linodes) > 0 {
			return fmt.Errorf("Linode still exists %+v", linodes)
		} else {
			return nil
		}
	}

	return nil
}

func testAccCheckLinodeLinodeExists(n string, instance *linodego.Linode) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found %+v", rs)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Linode id set")
		}

		client := testAccProvider.Meta().(*linodego.Client)
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			panic(err)
		}

		linodes, err := client.Linode.List(id)
		if err != nil {
			return err
		} else if len(linodes.Linodes) != 1 {
			return fmt.Errorf("Unexpected linode list response for %d: %+v", id, linodes)
		}

		*instance = linodes.Linodes[0]

		return nil
	}
}

func testAccCheckLinodeLinodeAttributesPrivateNetworking(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found %+v", rs)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Linode id set")
		}

		client := testAccProvider.Meta().(*linodego.Client)
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			panic(err)
		}
		_, err = client.Linode.List(id)
		if err != nil {
			return err
		}

		_, privateIP, err := getIps(client, int(id))
		if err != nil {
			return err
		}

		if privateIP == "" {
			return fmt.Errorf("Private Ip is not set")
		}
		return nil
	}
}

func testAccCheckLinodeLinodeConfigBasic(instance string, pubkey string) string {
	return fmt.Sprintf(`
resource "linode_linode" "foobar" {
	name = "%s"
	group = "testing"
	size = 1024
	image = "Ubuntu 16.04 LTS"
	region = "Dallas, TX, USA"
	kernel = "Latest 64 bit"
	root_password = "terraform-test"
	swap_size = 256
	ssh_key = "%s"
}`, instance, pubkey)
}

func testAccCheckLinodeLinodeConfigUpdates(instance string, pubkey string) string {
	return fmt.Sprintf(`
resource "linode_linode" "foobar" {
	name = "%s_renamed"
	group = "integration"
	size = 1024
	image = "Ubuntu 16.04 LTS"
	region = "Dallas, TX, USA"
	kernel = "Latest 64 bit"
	root_password = "terraform-test"
	swap_size = 256
	ssh_key = "%s"
}`, instance, pubkey)
}

func testAccCheckLinodeLinodeConfigUpsizeSmall(instance string, pubkey string) string {
	return fmt.Sprintf(`
resource "linode_linode" "foobar" {
	name = "%s"
	group = "integration"
	size = 1024
	image = "Ubuntu 16.04 LTS"
	region = "Dallas, TX, USA"
	kernel = "Latest 64 bit"
	root_password = "terraform-test"
	swap_size = 256
	ssh_key = "%s"
}`, instance, pubkey)
}

func testAccCheckLinodeLinodeConfigUpsizeBigger(instance string, pubkey string) string {
	return fmt.Sprintf(`
resource "linode_linode" "foobar" {
	name = "%s_upsized"
	group = "integration"
	size = 2048
	image = "Ubuntu 16.04 LTS"
	region = "Dallas, TX, USA"
	kernel = "Latest 64 bit"
	root_password = "terraform-test"
	swap_size = 256
	ssh_key = "%s"
}`, instance, pubkey)
}

func testAccCheckLinodeLinodeConfigDownsize(instance string, pubkey string) string {
	return fmt.Sprintf(`
resource "linode_linode" "foobar" {
	name = "%s_downsized"
	group = "integration"
	size = 1024
	image = "Ubuntu 16.04 LTS"
	region = "Dallas, TX, USA"
	kernel = "Latest 64 bit"
	root_password = "terraform-test"
	swap_size = 256
	ssh_key = "%s"
}`, instance, pubkey)
}

func testAccCheckLinodeLinodeConfigUpsizeExpandDisk(instance string, pubkey string) string {
	return fmt.Sprintf(`
resource "linode_linode" "foobar" {
	name = "%s_expanded"
	group = "integration"
	size = 2048
	disk_expansion = true
	image = "Ubuntu 16.04 LTS"
	region = "Dallas, TX, USA"
	kernel = "Latest 64 bit"
	root_password = "terraform-test"
	swap_size = 256
	ssh_key = "%s"
}`, instance, pubkey)
}

func testAccCheckLinodeLinodeConfigPrivateNetworking(instance string, pubkey string) string {
	return fmt.Sprintf(`
resource "linode_linode" "foobar" {
	name = "%s"
	group = "integration"
	size = 1024
	image = "Ubuntu 16.04 LTS"
	region = "Dallas, TX, USA"
	kernel = "Latest 64 bit"
	root_password = "terraform-test"
	swap_size = 256
	private_networking = true
	ssh_key = "%s"
}`, instance, pubkey)
}
