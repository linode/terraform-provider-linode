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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeLinodeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeLinodeConfigBasic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeLinodeExists("linode_linode.foobar", &instance),
					resource.TestCheckResourceAttr("linode_linode.foobar", "name", instanceName),
					resource.TestCheckResourceAttr("linode_linode.foobar", "size", "1024"),
					resource.TestCheckResourceAttr("linode_linode.foobar", "image", "Ubuntu 14.04 LTS"),
					resource.TestCheckResourceAttr("linode_linode.foobar", "region", "Dallas, TX, USA"),
					resource.TestCheckResourceAttr("linode_linode.foobar", "kernel", "Latest 64 bit"),
					resource.TestCheckResourceAttr("linode_linode.foobar", "group", "testing"),
					resource.TestCheckResourceAttr("linode_linode.foobar", "swap_size", "256"),
				),
			},
			resource.TestStep{
				ResourceName:  "linode_linode.foobar",
				ImportState:   true,
				ImportStateId: fmt.Sprintf("%s", instance.LinodeId),
			},
		},
	})
}

func TestAccLinodeLinodeUpdate(t *testing.T) {
	t.Parallel()

	var instance linodego.Linode
	var instanceName = fmt.Sprintf("tf_test_%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeLinodeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeLinodeConfigBasic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeLinodeExists("linode_linode.foobar", &instance),
					resource.TestCheckResourceAttr("linode_linode.foobar", "name", fmt.Sprintf(instanceName)),
					resource.TestCheckResourceAttr("linode_linode.foobar", "group", "testing"),
				),
			},
			resource.TestStep{
				Config: testAccCheckLinodeLinodeConfigUpdates(instanceName),
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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeLinodeDestroy,
		Steps: []resource.TestStep{
			// Start off with a Linode 1024
			resource.TestStep{
				Config: testAccCheckLinodeLinodeConfigUpsizeSmall(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeLinodeExists("linode_linode.foobar", &instance),
					resource.TestCheckResourceAttr("linode_linode.foobar", "size", "1024"),
					resource.TestCheckResourceAttr("linode_linode.foobar", "plan_storage_utilized", "20480"),
				),
			},
			// Bump it to a 2048, but don't expand the disk
			resource.TestStep{
				Config: testAccCheckLinodeLinodeConfigUpsizeBigger(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeLinodeExists("linode_linode.foobar", &instance),
					resource.TestCheckResourceAttr("linode_linode.foobar", "size", "2048"),
					resource.TestCheckResourceAttr("linode_linode.foobar", "plan_storage_utilized", "20480"),
				),
			},
			// Go back down to a 1024
			resource.TestStep{
				Config: testAccCheckLinodeLinodeConfigDownsize(instanceName),
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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeLinodeDestroy,
		Steps: []resource.TestStep{
			// Start off with a Linode 1024
			resource.TestStep{
				Config: testAccCheckLinodeLinodeConfigUpsizeSmall(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeLinodeExists("linode_linode.foobar", &instance),
					resource.TestCheckResourceAttr("linode_linode.foobar", "size", "1024"),
					resource.TestCheckResourceAttr("linode_linode.foobar", "plan_storage_utilized", "20480"),
				),
			},
			// Bump it to a 2048, and expand the disk
			resource.TestStep{
				Config: testAccCheckLinodeLinodeConfigUpsizeExpandDisk(instanceName),
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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeLinodeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeLinodeConfigPrivateNetworking(instanceName),
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

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse %s as int", rs.Primary.ID)
		}

		fmt.Println("Going to look for linode %s", id)
		response, err := client.Linode.List(int(id))
		fmt.Println(response)
		if err == nil {
			return fmt.Errorf("Linode still exists %s", err)
		}
	}

	return nil
}

func testAccCheckLinodeLinodeExists(n string, instance *linodego.Linode) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found %s", rs)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Linode id set")
		}

		client := testAccProvider.Meta().(*linodego.Client)
		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			panic(err)
		}

		linodes, err := client.Linode.List(int(id))
		if err != nil {
			return err
		}

		*instance = linodes.Linodes[0]

		return nil
	}
}

func testAccCheckLinodeLinodeAttributesPrivateNetworking(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found %s", rs)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Linode id set")
		}

		client := testAccProvider.Meta().(*linodego.Client)
		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			panic(err)
		}
		_, err = client.Linode.List(int(id))
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

func testAccCheckLinodeLinodeConfigBasic(instance string) string {
	return fmt.Sprintf(`
resource "linode_linode" "foobar" {
	name = "%s"
	group = "testing"
	size = 1024
	image = "Ubuntu 14.04 LTS"
	region = "Dallas, TX, USA"
	kernel = "Latest 64 bit"
	root_password = "terraform-test"
	swap_size = 256
	ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCxtdizvJzTT38y2oXuoLUXbLUf9V0Jy9KsM0bgIvjUCSEbuLWCXKnWqgBmkv7iTKGZg3fx6JA10hiufdGHD7at5YaRUitGP2mvC2I68AYNZmLCGXh0hYMrrUB01OEXHaYhpSmXIBc9zUdTreL5CvYe3PAYzuBA0/lGFTnNsHosSd+suA4xfJWMr/Fr4/uxrpcy8N8BE16pm4kci5tcMh6rGUGtDEj6aE9k8OI4SRmSZJsNElsu/Z/K4zqCpkW/U06vOnRrE98j3NE07nxVOTqdAMZqopFiMP0MXWvd6XyS2/uKU+COLLc0+hVsgj+dVMTWfy8wZ58OJDsIKk/cI/7yF+GZz89Js+qYx7u9mNhpEgD4UrcRHpitlRgVhA8p6R4oBqb0m/rpKBd2BAFdcty3GIP9CWsARtsCbN6YDLJ1JN3xI34jSGC1ROktVHg27bEEiT5A75w3WJl96BlSo5zJsIZDTWlaqnr26YxNHba4ILdVLKigQtQpf8WFsnB9YzmDdb9K3w9szf5lAkb/SFXw+e+yPS9habkpOncL0oCsgag5wUGCEmZ7wpiY8QgARhuwsQUkxv1aUi/Nn7b7sAkKSkxtBI3LBXZ+vcUxZTH0ut4pe9rbrEed3ktAOF5FafjA1VtarPqqZ+g46xVO9llgpXcl3rVglFtXzTcUy09hGw== btobolaski@Brendans-MacBook-Pro.local"
}`, instance)
}

func testAccCheckLinodeLinodeConfigUpdates(instance string) string {
	return fmt.Sprintf(`
resource "linode_linode" "foobar" {
	name = "%s_renamed"
	group = "integration"
	size = 1024
	image = "Ubuntu 14.04 LTS"
	region = "Dallas, TX, USA"
	kernel = "Latest 64 bit"
	root_password = "terraform-test"
	swap_size = 256
	ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCxtdizvJzTT38y2oXuoLUXbLUf9V0Jy9KsM0bgIvjUCSEbuLWCXKnWqgBmkv7iTKGZg3fx6JA10hiufdGHD7at5YaRUitGP2mvC2I68AYNZmLCGXh0hYMrrUB01OEXHaYhpSmXIBc9zUdTreL5CvYe3PAYzuBA0/lGFTnNsHosSd+suA4xfJWMr/Fr4/uxrpcy8N8BE16pm4kci5tcMh6rGUGtDEj6aE9k8OI4SRmSZJsNElsu/Z/K4zqCpkW/U06vOnRrE98j3NE07nxVOTqdAMZqopFiMP0MXWvd6XyS2/uKU+COLLc0+hVsgj+dVMTWfy8wZ58OJDsIKk/cI/7yF+GZz89Js+qYx7u9mNhpEgD4UrcRHpitlRgVhA8p6R4oBqb0m/rpKBd2BAFdcty3GIP9CWsARtsCbN6YDLJ1JN3xI34jSGC1ROktVHg27bEEiT5A75w3WJl96BlSo5zJsIZDTWlaqnr26YxNHba4ILdVLKigQtQpf8WFsnB9YzmDdb9K3w9szf5lAkb/SFXw+e+yPS9habkpOncL0oCsgag5wUGCEmZ7wpiY8QgARhuwsQUkxv1aUi/Nn7b7sAkKSkxtBI3LBXZ+vcUxZTH0ut4pe9rbrEed3ktAOF5FafjA1VtarPqqZ+g46xVO9llgpXcl3rVglFtXzTcUy09hGw== btobolaski@Brendans-MacBook-Pro.local"
}`, instance)
}

func testAccCheckLinodeLinodeConfigUpsizeSmall(instance string) string {
	return fmt.Sprintf(`
resource "linode_linode" "foobar" {
	name = "%s"
	group = "integration"
	size = 1024
	image = "Ubuntu 14.04 LTS"
	region = "Dallas, TX, USA"
	kernel = "Latest 64 bit"
	root_password = "terraform-test"
	swap_size = 256
	ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCxtdizvJzTT38y2oXuoLUXbLUf9V0Jy9KsM0bgIvjUCSEbuLWCXKnWqgBmkv7iTKGZg3fx6JA10hiufdGHD7at5YaRUitGP2mvC2I68AYNZmLCGXh0hYMrrUB01OEXHaYhpSmXIBc9zUdTreL5CvYe3PAYzuBA0/lGFTnNsHosSd+suA4xfJWMr/Fr4/uxrpcy8N8BE16pm4kci5tcMh6rGUGtDEj6aE9k8OI4SRmSZJsNElsu/Z/K4zqCpkW/U06vOnRrE98j3NE07nxVOTqdAMZqopFiMP0MXWvd6XyS2/uKU+COLLc0+hVsgj+dVMTWfy8wZ58OJDsIKk/cI/7yF+GZz89Js+qYx7u9mNhpEgD4UrcRHpitlRgVhA8p6R4oBqb0m/rpKBd2BAFdcty3GIP9CWsARtsCbN6YDLJ1JN3xI34jSGC1ROktVHg27bEEiT5A75w3WJl96BlSo5zJsIZDTWlaqnr26YxNHba4ILdVLKigQtQpf8WFsnB9YzmDdb9K3w9szf5lAkb/SFXw+e+yPS9habkpOncL0oCsgag5wUGCEmZ7wpiY8QgARhuwsQUkxv1aUi/Nn7b7sAkKSkxtBI3LBXZ+vcUxZTH0ut4pe9rbrEed3ktAOF5FafjA1VtarPqqZ+g46xVO9llgpXcl3rVglFtXzTcUy09hGw== btobolaski@Brendans-MacBook-Pro.local"
}`, instance)
}

func testAccCheckLinodeLinodeConfigUpsizeBigger(instance string) string {
	return fmt.Sprintf(`
resource "linode_linode" "foobar" {
	name = "%s_upsized"
	group = "integration"
	size = 2048
	image = "Ubuntu 14.04 LTS"
	region = "Dallas, TX, USA"
	kernel = "Latest 64 bit"
	root_password = "terraform-test"
	swap_size = 256
	ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCxtdizvJzTT38y2oXuoLUXbLUf9V0Jy9KsM0bgIvjUCSEbuLWCXKnWqgBmkv7iTKGZg3fx6JA10hiufdGHD7at5YaRUitGP2mvC2I68AYNZmLCGXh0hYMrrUB01OEXHaYhpSmXIBc9zUdTreL5CvYe3PAYzuBA0/lGFTnNsHosSd+suA4xfJWMr/Fr4/uxrpcy8N8BE16pm4kci5tcMh6rGUGtDEj6aE9k8OI4SRmSZJsNElsu/Z/K4zqCpkW/U06vOnRrE98j3NE07nxVOTqdAMZqopFiMP0MXWvd6XyS2/uKU+COLLc0+hVsgj+dVMTWfy8wZ58OJDsIKk/cI/7yF+GZz89Js+qYx7u9mNhpEgD4UrcRHpitlRgVhA8p6R4oBqb0m/rpKBd2BAFdcty3GIP9CWsARtsCbN6YDLJ1JN3xI34jSGC1ROktVHg27bEEiT5A75w3WJl96BlSo5zJsIZDTWlaqnr26YxNHba4ILdVLKigQtQpf8WFsnB9YzmDdb9K3w9szf5lAkb/SFXw+e+yPS9habkpOncL0oCsgag5wUGCEmZ7wpiY8QgARhuwsQUkxv1aUi/Nn7b7sAkKSkxtBI3LBXZ+vcUxZTH0ut4pe9rbrEed3ktAOF5FafjA1VtarPqqZ+g46xVO9llgpXcl3rVglFtXzTcUy09hGw== btobolaski@Brendans-MacBook-Pro.local"
}`, instance)
}

func testAccCheckLinodeLinodeConfigDownsize(instance string) string {
	return fmt.Sprintf(`
resource "linode_linode" "foobar" {
	name = "%s_downsized"
	group = "integration"
	size = 1024
	image = "Ubuntu 14.04 LTS"
	region = "Dallas, TX, USA"
	kernel = "Latest 64 bit"
	root_password = "terraform-test"
	swap_size = 256
	ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCxtdizvJzTT38y2oXuoLUXbLUf9V0Jy9KsM0bgIvjUCSEbuLWCXKnWqgBmkv7iTKGZg3fx6JA10hiufdGHD7at5YaRUitGP2mvC2I68AYNZmLCGXh0hYMrrUB01OEXHaYhpSmXIBc9zUdTreL5CvYe3PAYzuBA0/lGFTnNsHosSd+suA4xfJWMr/Fr4/uxrpcy8N8BE16pm4kci5tcMh6rGUGtDEj6aE9k8OI4SRmSZJsNElsu/Z/K4zqCpkW/U06vOnRrE98j3NE07nxVOTqdAMZqopFiMP0MXWvd6XyS2/uKU+COLLc0+hVsgj+dVMTWfy8wZ58OJDsIKk/cI/7yF+GZz89Js+qYx7u9mNhpEgD4UrcRHpitlRgVhA8p6R4oBqb0m/rpKBd2BAFdcty3GIP9CWsARtsCbN6YDLJ1JN3xI34jSGC1ROktVHg27bEEiT5A75w3WJl96BlSo5zJsIZDTWlaqnr26YxNHba4ILdVLKigQtQpf8WFsnB9YzmDdb9K3w9szf5lAkb/SFXw+e+yPS9habkpOncL0oCsgag5wUGCEmZ7wpiY8QgARhuwsQUkxv1aUi/Nn7b7sAkKSkxtBI3LBXZ+vcUxZTH0ut4pe9rbrEed3ktAOF5FafjA1VtarPqqZ+g46xVO9llgpXcl3rVglFtXzTcUy09hGw== btobolaski@Brendans-MacBook-Pro.local"
}`, instance)
}

func testAccCheckLinodeLinodeConfigUpsizeExpandDisk(instance string) string {
	return fmt.Sprintf(`
resource "linode_linode" "foobar" {
	name = "%s_expanded"
	group = "integration"
	size = 2048
	disk_expansion = true
	image = "Ubuntu 14.04 LTS"
	region = "Dallas, TX, USA"
	kernel = "Latest 64 bit"
	root_password = "terraform-test"
	swap_size = 256
	ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCxtdizvJzTT38y2oXuoLUXbLUf9V0Jy9KsM0bgIvjUCSEbuLWCXKnWqgBmkv7iTKGZg3fx6JA10hiufdGHD7at5YaRUitGP2mvC2I68AYNZmLCGXh0hYMrrUB01OEXHaYhpSmXIBc9zUdTreL5CvYe3PAYzuBA0/lGFTnNsHosSd+suA4xfJWMr/Fr4/uxrpcy8N8BE16pm4kci5tcMh6rGUGtDEj6aE9k8OI4SRmSZJsNElsu/Z/K4zqCpkW/U06vOnRrE98j3NE07nxVOTqdAMZqopFiMP0MXWvd6XyS2/uKU+COLLc0+hVsgj+dVMTWfy8wZ58OJDsIKk/cI/7yF+GZz89Js+qYx7u9mNhpEgD4UrcRHpitlRgVhA8p6R4oBqb0m/rpKBd2BAFdcty3GIP9CWsARtsCbN6YDLJ1JN3xI34jSGC1ROktVHg27bEEiT5A75w3WJl96BlSo5zJsIZDTWlaqnr26YxNHba4ILdVLKigQtQpf8WFsnB9YzmDdb9K3w9szf5lAkb/SFXw+e+yPS9habkpOncL0oCsgag5wUGCEmZ7wpiY8QgARhuwsQUkxv1aUi/Nn7b7sAkKSkxtBI3LBXZ+vcUxZTH0ut4pe9rbrEed3ktAOF5FafjA1VtarPqqZ+g46xVO9llgpXcl3rVglFtXzTcUy09hGw== btobolaski@Brendans-MacBook-Pro.local"
}`, instance)
}

func testAccCheckLinodeLinodeConfigPrivateNetworking(instance string) string {
	return fmt.Sprintf(`
resource "linode_linode" "foobar" {
	name = "%s"
	group = "integration"
	size = 1024
	image = "Ubuntu 14.04 LTS"
	region = "Dallas, TX, USA"
	kernel = "Latest 64 bit"
	root_password = "terraform-test"
	swap_size = 256
	private_networking = true
	ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCxtdizvJzTT38y2oXuoLUXbLUf9V0Jy9KsM0bgIvjUCSEbuLWCXKnWqgBmkv7iTKGZg3fx6JA10hiufdGHD7at5YaRUitGP2mvC2I68AYNZmLCGXh0hYMrrUB01OEXHaYhpSmXIBc9zUdTreL5CvYe3PAYzuBA0/lGFTnNsHosSd+suA4xfJWMr/Fr4/uxrpcy8N8BE16pm4kci5tcMh6rGUGtDEj6aE9k8OI4SRmSZJsNElsu/Z/K4zqCpkW/U06vOnRrE98j3NE07nxVOTqdAMZqopFiMP0MXWvd6XyS2/uKU+COLLc0+hVsgj+dVMTWfy8wZ58OJDsIKk/cI/7yF+GZz89Js+qYx7u9mNhpEgD4UrcRHpitlRgVhA8p6R4oBqb0m/rpKBd2BAFdcty3GIP9CWsARtsCbN6YDLJ1JN3xI34jSGC1ROktVHg27bEEiT5A75w3WJl96BlSo5zJsIZDTWlaqnr26YxNHba4ILdVLKigQtQpf8WFsnB9YzmDdb9K3w9szf5lAkb/SFXw+e+yPS9habkpOncL0oCsgag5wUGCEmZ7wpiY8QgARhuwsQUkxv1aUi/Nn7b7sAkKSkxtBI3LBXZ+vcUxZTH0ut4pe9rbrEed3ktAOF5FafjA1VtarPqqZ+g46xVO9llgpXcl3rVglFtXzTcUy09hGw== btobolaski@Brendans-MacBook-Pro.local"
}`, instance)
}
