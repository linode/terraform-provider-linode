//go:build integration || iamuser

package iamuser_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/iamuser/tmpl"
)

var testRegion string

func init() {
	resource.AddTestSweepers("linode_user", &resource.Sweeper{
		Name: "linode_user",
		F:    userSweep,
	})

	resource.AddTestSweepers("linode_volume", &resource.Sweeper{
		Name: "linode_volume",
		F:    volumeSweep,
	})

	region, err := acceptance.GetRandomRegionWithCaps([]string{linodego.CapabilityBlockStorage}, "core")
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func userSweep(prefix string) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	listOpts := acceptance.SweeperListOptions(prefix, "username")
	users, err := client.ListUsers(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("Error getting users: %s", err)
	}
	for _, user := range users {
		if !acceptance.ShouldSweep(prefix, user.Username) {
			continue
		}
		err := client.DeleteUser(context.Background(), user.Username)
		if err != nil {
			return fmt.Errorf("Error destroying %s during sweep: %s", user.Username, err)
		}
	}

	return nil
}

func volumeSweep(prefix string) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	listOpts := acceptance.SweeperListOptions(prefix, "label")
	volumes, err := client.ListVolumes(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("Error getting volumes: %s", err)
	}
	for _, volume := range volumes {
		if !acceptance.ShouldSweep(prefix, volume.Label) {
			continue
		}
		err := client.DeleteVolume(context.Background(), volume.ID)
		if err != nil {
			return fmt.Errorf("Error destroying %s during sweep: %s", volume.Label, err)
		}
	}

	return nil
}

func TestAccResourceIAMUser_Update(t *testing.T) {
	t.Parallel()

	resName := "linode_iam_user.test_iam_user"
	volumeName := acctest.RandomWithPrefix("tf_test")
	username := acctest.RandomWithPrefix("tf_test")
	email := username + "@example.com"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Update(t, volumeName, testRegion, username, email, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "account_access.0", "account_event_viewer"),
					resource.TestCheckResourceAttr(resName, "entity_access.0.type", "volume"),
					resource.TestCheckResourceAttr(resName, "entity_access.0.roles.0", "volume_admin"),
				),
			},
		},
	})
}
