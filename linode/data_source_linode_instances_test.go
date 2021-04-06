package linode

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strconv"
	"testing"
)

func TestAccDataSourceLinodeInstances_basic(t *testing.T) {
	t.Parallel()

	resName := "data.linode_instances.foobar"
	instanceName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCheckLinodeInstancesBasic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "instance.#", "1"),
					resource.TestCheckResourceAttr(resName, "instance.0.type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "instance.0.tags.#", "2"),
					resource.TestCheckResourceAttr(resName, "instance.0.image", "linode/ubuntu18.04"),
					resource.TestCheckResourceAttr(resName, "instance.0.region", "us-east"),
					resource.TestCheckResourceAttr(resName, "instance.0.group", "tf_test"),
					resource.TestCheckResourceAttr(resName, "instance.0.swap_size", "256"),
					resource.TestCheckResourceAttr(resName, "instance.0.ipv4.#", "2"),
					resource.TestCheckResourceAttrSet(resName, "instance.0.ipv6"),
					resource.TestCheckResourceAttr(resName, "instance.0.disk.#", "2"),
					resource.TestCheckResourceAttr(resName, "instance.0.config.#", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceLinodeInstances_noFilters(t *testing.T) {
	t.Parallel()

	resName := "data.linode_instances.foobar"
	instanceName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCheckLinodeInstancesNoFilters(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceCheckLinodeAllInstancesNotEmpty(resName),
				),
			},
		},
	})
}

func TestAccDataSourceLinodeInstances_multipleInstances(t *testing.T) {
	resName := "data.linode_instances.foobar"
	instanceName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCheckLinodeInstancesMultipleInstances(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "instance.#", "3"),
				),
			},
		},
	})
}

func testAccDataSourceCheckLinodeAllInstancesNotEmpty(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		linodeCount, err := strconv.Atoi(rs.Primary.Attributes["instance.#"])
		if err != nil {
			return fmt.Errorf("failed to parse: %s", err)
		}

		if linodeCount < 1 {
			return fmt.Errorf("expected at least 1 linode instance")
		}

		return nil
	}
}

func testDataSourceCheckLinodeInstancesBasic(instance string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	group = "tf_test"
	tags = ["cool", "cooler"]
	type = "g6-nanode-1"
	image = "linode/ubuntu18.04"
	region = "us-east"
	root_pass = "terraform-test"
	swap_size = 256
	private_ip = true
}
`, instance) + `
data "linode_instances" "foobar" {
	filter {
		name = "id"
		values = [linode_instance.foobar.id]
	}

	filter {
		name = "label"
		values = [linode_instance.foobar.label, "other-label"]
	}

	filter {
		name = "group"
		values = [linode_instance.foobar.group]
	}

	filter {
		name = "region"
		values = [linode_instance.foobar.region]
	}

	filter {
		name = "tags"
		values = linode_instance.foobar.tags
	}
}
`
}

func testDataSourceCheckLinodeInstancesNoFilters(instance string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	group = "tf_test"
	tags = ["cool", "cooler"]
	type = "g6-nanode-1"
	image = "linode/ubuntu18.04"
	region = "us-east"
	root_pass = "terraform-test"
}
`, instance) + `
data "linode_instances" "foobar" {}
`
}

func testDataSourceCheckLinodeInstancesMultipleInstances(instance string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar-0" {
	label = "%s-0"
	group = "tf_test"
	tags = ["cool", "cooler"]
	type = "g6-nanode-1"
	image = "linode/ubuntu18.04"
	region = "us-east"
	root_pass = "terraform-test"
}

resource "linode_instance" "foobar-1" {
	label = "%s-1"
	group = "tf_test"
	tags = ["cool", "cooler"]
	type = "g6-nanode-1"
	image = "linode/ubuntu18.04"
	region = "us-east"
	root_pass = "terraform-test"
}

resource "linode_instance" "foobar-2" {
	label = "%s-2"
	group = "tf_test"
	tags = ["cool", "cooler"]
	type = "g6-nanode-1"
	image = "linode/ubuntu18.04"
	region = "us-east"
	root_pass = "terraform-test"
}
`, instance, instance, instance) + `
data "linode_instances" "foobar" {
	filter {
		name = "group"
		values = [linode_instance.foobar-0.group]
	}
}
`
}
