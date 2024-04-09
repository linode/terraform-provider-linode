//go:build integration || vpc

package vpc_test

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
	"github.com/linode/terraform-provider-linode/v2/linode/vpc/tmpl"
)

var testRegion string

func init() {
	resource.AddTestSweepers("linode_vpc", &resource.Sweeper{
		Name: "linode_vpc",
		F:    sweep,
	})

	var err error

	testRegion, err = acceptance.GetRandomRegionWithCaps([]string{"VPCs"})
	if err != nil {
		log.Fatal(fmt.Errorf("Error getting region: %s", err))
	}
}

func sweep(prefix string) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		log.Fatal(fmt.Errorf("Error getting client: %s", err))
	}

	listOpts := acceptance.SweeperListOptions(prefix, "label")
	vpcs, err := client.ListVPCs(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("Error getting VPCs: %s", err)
	}

	for _, vpc := range vpcs {
		if !acceptance.ShouldSweep(prefix, vpc.Label) {
			continue
		}
		err := client.DeleteVPC(context.Background(), vpc.ID)
		if err != nil {
			return fmt.Errorf("Error destroying %s during sweep: %s", vpc.Label, err)
		}
	}

	return nil
}

func TestAccResourceVPC_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_vpc.foobar"
	vpcLabel := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, vpcLabel, testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkVPCExists,
					resource.TestCheckResourceAttr(resName, "label", vpcLabel),
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttrSet(resName, "description"),
					resource.TestCheckResourceAttrSet(resName, "region"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttrSet(resName, "updated"),
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

func TestAccResourceVPC_update(t *testing.T) {
	t.Parallel()
	resName := "linode_vpc.foobar"
	vpcLabel := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, vpcLabel, testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkVPCExists,
					resource.TestCheckResourceAttr(resName, "label", vpcLabel),
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttrSet(resName, "created"),
				),
			},
			{
				Config: tmpl.Updates(t, vpcLabel, testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkVPCExists,
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s-renamed", vpcLabel)),
					resource.TestCheckResourceAttr(resName, "description", "some description"),
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttrSet(resName, "updated"),
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

func TestAccResourceLinodeVPC_create_InvalidLabel(t *testing.T) {
	t.Parallel()

	vpcLabel := "tf-test_123"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      tmpl.Basic(t, vpcLabel, testRegion),
				ExpectError: regexp.MustCompile("Label must include only ASCII letters, numbers, and dashes"),
			},
		},
	})
}

func TestAccResourceLinodeVPC_update_InvalidLabel(t *testing.T) {
	t.Parallel()
	resName := "linode_vpc.foobar"
	vpcLabel := acctest.RandomWithPrefix("tf-test")

	invalidLabel := "tf-test_123"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, vpcLabel, testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkVPCExists,
					resource.TestCheckResourceAttr(resName, "label", vpcLabel),
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttrSet(resName, "created"),
				),
			},
			{
				Config:      tmpl.Updates(t, invalidLabel, testRegion),
				ExpectError: regexp.MustCompile("Label must include only ASCII letters, numbers, and dashes"),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func checkVPCExists(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_vpc" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		_, err = client.GetVPC(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of VPC %s: %s", rs.Primary.Attributes["label"], err)
		}
	}

	return nil
}

func checkVPCDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_vpc" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}
		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)
		}

		_, err = client.GetVPC(context.Background(), id)

		if err == nil {
			return fmt.Errorf("Linode VPC with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error requesting Linode VPC with id %d", id)
		}
	}

	return nil
}
