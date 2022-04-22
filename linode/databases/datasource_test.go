package databases_test

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/databases/tmpl"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"log"
	"testing"
)

var engineVersion string

func init() {
	client, err := acceptance.GetClientForSweepers()
	if err != nil {
		log.Fatalf("failed to get client: %s", err)
	}

	v, err := helper.ResolveValidDBEngine(context.Background(), *client, "mysql")
	if err != nil {
		log.Fatalf("failde to get db engine version: %s", err)
	}

	engineVersion = v.ID
}

func TestAccDataSourceDatabases_byAttr(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_databases.foobar"
	dbName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.ByLabel(t, engineVersion, dbName, dbName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "databases.0.label", dbName),
					resource.TestCheckResourceAttr(resourceName, "databases.0.cluster_size", "1"),
					resource.TestCheckResourceAttr(resourceName, "databases.0.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "databases.0.engine", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "databases.0.region", "us-southeast"),
					resource.TestCheckResourceAttr(resourceName, "databases.0.type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resourceName, "databases.0.allow_list.#", "0"),

					resource.TestCheckResourceAttrSet(resourceName, "databases.0.created"),
					resource.TestCheckResourceAttrSet(resourceName, "databases.0.host_primary"),
					resource.TestCheckResourceAttrSet(resourceName, "databases.0.host_secondary"),
					resource.TestCheckResourceAttrSet(resourceName, "databases.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "databases.0.instance_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "databases.0.status"),
					resource.TestCheckResourceAttrSet(resourceName, "databases.0.updated"),
					resource.TestCheckResourceAttrSet(resourceName, "databases.0.version"),
				),
			},
			{
				Config: tmpl.ByLabel(t, engineVersion, dbName, "not"+dbName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "databases.#", "0"),
				),
			},
			{
				Config: tmpl.ByEngine(t, engineVersion, dbName, "mysql"),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckResourceAttrGreaterThan(resourceName, "databases.#", 0),
					resource.TestCheckResourceAttr(resourceName, "databases.0.engine", "mysql"),
				),
			},
		},
	})
}
