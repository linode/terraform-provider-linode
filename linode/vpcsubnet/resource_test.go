//go:build integration || vpcsubnet

package vpcsubnet_test

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/vpcsubnet/tmpl"
)

var testRegion string

func init() {
	r, err := acceptance.GetRandomRegionWithCaps([]string{"VPCs"})
	if err != nil {
		log.Fatal(fmt.Errorf("Error getting region: %s", err))
	}

	testRegion = r
}

func TestAccResourceVPCSubnet_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_vpc_subnet.foobar"
	subnetLabel := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkVPCSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, subnetLabel, "172.16.0.0/24", testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkVPCSubnetExists,
					resource.TestCheckResourceAttr(resName, "label", subnetLabel),
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttrSet(resName, "created"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: resourceImportStateID,
			},
		},
	})
}

func TestAccResourceVPCSubnet_update(t *testing.T) {
	t.Parallel()
	resName := "linode_vpc_subnet.foobar"
	subnetLabel := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkVPCSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, subnetLabel, "192.168.0.0/26", testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkVPCSubnetExists,
					resource.TestCheckResourceAttr(resName, "label", subnetLabel),
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttrSet(resName, "created"),
				),
			},
			{
				Config: tmpl.Updates(t, subnetLabel, "192.168.0.0/26", testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkVPCSubnetExists,
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s-renamed", subnetLabel)),
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttrSet(resName, "updated"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: resourceImportStateID,
			},
		},
	})
}

func TestAccResourceVPCSubnet_create_InvalidLabel_basic(t *testing.T) {
	t.Parallel()

	subnetLabel := acctest.RandomWithPrefix("tf-test") + "__"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkVPCSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config:      tmpl.Basic(t, subnetLabel, "172.16.0.0/24", testRegion),
				ExpectError: regexp.MustCompile("Label must include only ASCII letters, numbers, and dashes"),
			},
		},
	})
}

func TestAccResourceVPCSubnet_update_invalidLabel(t *testing.T) {
	t.Parallel()
	resName := "linode_vpc_subnet.foobar"
	subnetLabel := acctest.RandomWithPrefix("tf-test")

	invalidLabel := "invalid_test_label"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkVPCSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, subnetLabel, "192.168.0.0/26", testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkVPCSubnetExists,
					resource.TestCheckResourceAttr(resName, "label", subnetLabel),
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttrSet(resName, "created"),
				),
			},
			{
				Config:      tmpl.Updates(t, invalidLabel, "192.168.0.0/26", testRegion),
				ExpectError: regexp.MustCompile("Label must include only ASCII letters, numbers, and dashes"),
			},
		},
	})
}

func TestAccResourceVPCSubnet_attached(t *testing.T) {
	t.Parallel()

	resName := "linode_vpc_subnet.foobar"
	subnetLabel := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkVPCSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Attached(t, subnetLabel, "172.16.0.0/24", testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkVPCSubnetExists,
					resource.TestCheckResourceAttr(resName, "label", subnetLabel),
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttrSet(resName, "created"),
				),
			},
			{
				// Refresh the configuration so the `linodes` field is updated
				Config: tmpl.Attached(t, subnetLabel, "172.16.0.0/24", testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "linodes.#", "1"),
					resource.TestCheckResourceAttrSet(resName, "linodes.0.id"),
					resource.TestCheckResourceAttr(resName, "linodes.0.interfaces.#", "1"),
					resource.TestCheckResourceAttrSet(resName, "linodes.0.interfaces.0.id"),
					resource.TestCheckResourceAttr(resName, "linodes.0.interfaces.0.active", "false"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: resourceImportStateID,
			},
		},
	})
}

func checkVPCSubnetExists(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_vpc_subnet" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		vpcID, err := strconv.Atoi(rs.Primary.Attributes["vpc_id"])
		if err != nil {
			return fmt.Errorf("failed to parse vpc_id: %w", err)
		}

		_, err = client.GetVPCSubnet(context.Background(), vpcID, id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of VPC subnet %s: %s", rs.Primary.Attributes["label"], err)
		}
	}

	return nil
}

func checkVPCSubnetDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_vpc_subnet" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}
		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)
		}

		vpcID, err := strconv.Atoi(rs.Primary.Attributes["vpc_id"])
		if err != nil {
			return fmt.Errorf("failed to parse vpc_id: %w", err)
		}

		_, err = client.GetVPCSubnet(context.Background(), vpcID, id)

		if err == nil {
			return fmt.Errorf("Linode VPC subnet with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error requesting Linode VPC subnet with id %d", id)
		}
	}

	return nil
}

func resourceImportStateID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_vpc_subnet" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return "", fmt.Errorf("Error parsing ID %v to int", rs.Primary.ID)
		}
		vpcID, err := strconv.Atoi(rs.Primary.Attributes["vpc_id"])
		if err != nil {
			return "", fmt.Errorf("Error parsing vpc_id %v to int", rs.Primary.Attributes["vpc_id"])
		}
		return fmt.Sprintf("%d,%d", vpcID, id), nil
	}

	return "", fmt.Errorf("Error finding linode_vpc_subnet")
}
