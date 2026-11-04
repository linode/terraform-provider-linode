//go:build integration || nfsspace

package nfsspace_test

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/nfsspace/tmpl"
)

const resourceName = "linode_nfs_space.foobar"

func init() {
	resource.AddTestSweepers("linode_nfs_space", &resource.Sweeper{
		Name: "linode_nfs_space",
		F:    sweep,
	})
}

func sweep(prefix string) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	spaces, err := client.ListNFSSpaces(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("error getting NFS Spaces: %s", err)
	}

	for _, space := range spaces {
		if !acceptance.ShouldSweep(prefix, space.Label) {
			continue
		}
		if err := client.DeleteNFSSpace(context.Background(), space.ID); err != nil {
			return fmt.Errorf("error destroying %v during sweep: %s", space.Label, err)
		}
	}

	return nil
}

func TestAccResourceNFSSpace_basic(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkNFSSpaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("label"), knownvalue.StringExact(label)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("status"), knownvalue.StringExact(string(linodego.NFSSpaceStatusActive))),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("max_capacity_bytes"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("used_capacity_bytes"), knownvalue.NotNull()),
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNFSSpace_update(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")
	newLabel := acctest.RandomWithPrefix("tf-test")
	description := "updated description"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkNFSSpaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("label"), knownvalue.StringExact(label)),
				},
			},
			{
				Config: tmpl.Updates(t, newLabel, description),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("label"), knownvalue.StringExact(newLabel)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("description"), knownvalue.StringExact(description)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("tags"), knownvalue.SetSizeExact(2)),
				},
			},
		},
	})
}

func checkNFSSpaceDestroy(s *terraform.State) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		log.Fatalf("failed to get client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_nfs_space" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error parsing %v to int", rs.Primary.ID)
		}

		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)
		}

		_, err = client.GetNFSSpace(context.Background(), id)
		if err == nil {
			return fmt.Errorf("NFS Space with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("error requesting NFS Space with id %d", id)
		}
	}

	return nil
}
