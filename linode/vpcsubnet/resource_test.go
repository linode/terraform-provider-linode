//go:build integration

package vpcsubnet_test

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/vpcsubnet/tmpl"
)

var vpcID int

func init() {
	resource.AddTestSweepers("linode_vpc_subnet", &resource.Sweeper{
		Name: "linode_vpc_subnet",
		F:    sweep,
	})

	client, err := acceptance.GetTestClient()
	if err != nil {
		log.Fatal(fmt.Errorf("Error getting client: %s", err))
	}

	testRegion, err := acceptance.GetRandomRegionWithCaps([]string{"VPCs"})
	if err != nil {
		log.Fatal(fmt.Errorf("Error getting region: %s", err))
	}

	vpc, err := client.CreateVPC(context.Background(), linodego.VPCCreateOptions{
		Label:  acctest.RandomWithPrefix("tf-test"),
		Region: testRegion,
	})
	if err != nil {
		log.Fatal(err)
	}

	vpcID = vpc.ID
}

func sweep(prefix string) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		log.Fatal(fmt.Errorf("Error getting client: %s", err))
	}

	listOpts := acceptance.SweeperListOptions(prefix, "label")
	subnets, err := client.ListVPCSubnet(context.Background(), vpcID, listOpts)
	if err != nil {
		return fmt.Errorf("Error getting VPC subnets: %s", err)
	}
	for _, subnet := range subnets {
		if !acceptance.ShouldSweep(prefix, subnet.Label) {
			continue
		}
		err := client.DeleteVPCSubnet(context.Background(), vpcID, subnet.ID)
		if err != nil {
			return fmt.Errorf("Error destroying %s during sweep: %s", subnet.Label, err)
		}
	}

	return nil
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
				Config: tmpl.Basic(t, vpcID, subnetLabel, "172.16.0.0/24"),
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
				Config: tmpl.Basic(t, vpcID, subnetLabel, "192.168.0.0/26"),
				Check: resource.ComposeTestCheckFunc(
					checkVPCSubnetExists,
					resource.TestCheckResourceAttr(resName, "label", subnetLabel),
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttrSet(resName, "created"),
				),
			},
			{
				Config: tmpl.Updates(t, vpcID, subnetLabel, "192.168.0.0/26"),
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
