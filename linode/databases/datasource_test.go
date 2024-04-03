//go:build integration || databases

package databases_test

import (
	"context"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/databases/tmpl"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

var (
	testRegion    string
	engineVersion string
)

func init() {
	client, err := acceptance.GetTestClient()
	if err != nil {
		log.Fatalf("failed to get client: %s", err)
	}

	v, err := helper.ResolveValidDBEngine(context.Background(), *client, "mysql")
	if err != nil {
		log.Fatalf("failed to get db engine version: %s", err)
	}

	engineVersion = v.ID

	region, err := acceptance.GetRandomRegionWithCaps([]string{"Managed Databases"})
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccDataSourceDatabases_byAttr(t *testing.T) {
	acceptance.LongRunningTest(t)
	t.Parallel()

	resourceName := "data.linode_databases.foobar"
	dbName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.ByLabel(t, engineVersion, dbName, dbName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "databases.0.label", dbName),
					resource.TestCheckResourceAttr(resourceName, "databases.0.cluster_size", "1"),
					resource.TestCheckResourceAttr(resourceName, "databases.0.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "databases.0.engine", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "databases.0.region", testRegion),
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
				Config: tmpl.ByLabel(t, engineVersion, dbName, "not"+dbName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "databases.#", "0"),
				),
			},
			{
				Config: tmpl.ByEngine(t, engineVersion, dbName, "mysql", testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckResourceAttrGreaterThan(resourceName, "databases.#", 0),
					resource.TestCheckResourceAttr(resourceName, "databases.0.engine", "mysql"),
				),
			},
		},
	})
}
