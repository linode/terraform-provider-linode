package backup_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func TestAccDataSourceInstanceBackups_basic(t *testing.T) {
	t.Parallel()

	instanceName := acctest.RandomWithPrefix("tf_test")
	snapshotName := acctest.RandomWithPrefix("tf_test_cool")

	resourceName := "data.linode_instance_backups.foobar"

	var instance linodego.Instance
	var snapshot *linodego.InstanceSnapshot

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceInstanceBasic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists("linode_instance.foobar", &instance),
				),
			},
			{
				PreConfig: func() {
					client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
					newSnapshot, err := client.CreateInstanceSnapshot(context.Background(), instance.ID, snapshotName)
					if err != nil {
						t.Fatal(err)
					}

					snapshot = newSnapshot
				},
				Config: dataSourceConfigBasic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "in_progress.0.id"),
					resource.TestCheckResourceAttr(resourceName, "in_progress.0.label", snapshotName),
					resource.TestCheckResourceAttrSet(resourceName, "in_progress.0.status"),
					resource.TestCheckResourceAttrSet(resourceName, "in_progress.0.type"),
					resource.TestCheckResourceAttrSet(resourceName, "in_progress.0.created"),
				),
			},
			{
				PreConfig: func() {
					client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
					if _, err := client.WaitForSnapshotStatus(context.Background(), instance.ID, snapshot.ID, linodego.SnapshotSuccessful, 600); err != nil {
						t.Fatal(err)
					}
				},
				Config: dataSourceConfigBasic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "current.0.id"),
					resource.TestCheckResourceAttr(resourceName, "current.0.label", snapshotName),
					resource.TestCheckResourceAttr(resourceName, "current.0.status", "successful"),
					resource.TestCheckResourceAttrSet(resourceName, "current.0.type"),
					resource.TestCheckResourceAttrSet(resourceName, "current.0.created"),
					resource.TestCheckResourceAttrSet(resourceName, "current.0.updated"),
					resource.TestCheckResourceAttrSet(resourceName, "current.0.finished"),
				),
			},
		},
	})
}

func resourceInstanceBasic(label string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	type = "g6-nanode-1"
	image = "linode/alpine3.13"
	region = "us-east"
	root_pass = "terraform-test"
	swap_size = 256
	backups_enabled = true
}`, label)
}

func dataSourceConfigBasic(instanceLabel string) string {
	return resourceInstanceBasic(instanceLabel) + `
data "linode_instance_backups" "foobar" {
	linode_id = linode_instance.foobar.id
}`
}
