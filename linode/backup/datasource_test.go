package backup_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

var testRegion string

func init() {
	region, err := acceptance.GetRandomRegionWithCaps(nil)
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccDataSourceInstanceBackups_basic(t *testing.T) {
	t.Parallel()

	instanceName := acctest.RandomWithPrefix("tf_test")
	snapshotName := acctest.RandomWithPrefix("tf_test_cool")

	resourceName := "data.linode_instance_backups.foobar"

	var instance linodego.Instance
	var snapshot *linodego.InstanceSnapshot

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceInstanceBasic(instanceName, testRegion),
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
				Config: dataSourceConfigBasic(instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "in_progress.0.id"),
					resource.TestCheckResourceAttr(resourceName, "in_progress.0.label", snapshotName),
					resource.TestCheckResourceAttrSet(resourceName, "in_progress.0.status"),
					resource.TestCheckResourceAttrSet(resourceName, "in_progress.0.type"),
					resource.TestCheckResourceAttrSet(resourceName, "in_progress.0.created"),
					resource.TestCheckResourceAttrSet(resourceName, "in_progress.0.available"),
				),
			},
			{
				PreConfig: func() {
					client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
					if _, err := client.WaitForSnapshotStatus(context.Background(), instance.ID, snapshot.ID, linodego.SnapshotSuccessful, 600); err != nil {
						t.Fatal(err)
					}
				},
				Config: dataSourceConfigBasic(instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "current.0.id"),
					resource.TestCheckResourceAttr(resourceName, "current.0.label", snapshotName),
					resource.TestCheckResourceAttr(resourceName, "current.0.status", "successful"),
					resource.TestCheckResourceAttr(resourceName, "current.0.available", "true"),
					resource.TestCheckResourceAttr("linode_instance.foobar", "backups.0.available", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "current.0.type"),
					resource.TestCheckResourceAttrSet(resourceName, "current.0.created"),
					resource.TestCheckResourceAttrSet(resourceName, "current.0.updated"),
					resource.TestCheckResourceAttrSet(resourceName, "current.0.finished"),
				),
			},
		},
	})
}

func resourceInstanceBasic(label, region string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	type = "g6-nanode-1"
	image = "linode/alpine3.15"
	region = "%s"
	root_pass = "terraform-test"
	swap_size = 256
	backups_enabled = true
}`, label, region)
}

func dataSourceConfigBasic(instanceLabel, region string) string {
	return resourceInstanceBasic(instanceLabel, region) + `
data "linode_instance_backups" "foobar" {
	linode_id = linode_instance.foobar.id
}`
}
