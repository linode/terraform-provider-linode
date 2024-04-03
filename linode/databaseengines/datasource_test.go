//go:build integration || databaseengines

package databaseengines_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/databaseengines/tmpl"
)

func TestAccDataSourceDatabaseEngines_all(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_database_engines.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataAll(t),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckResourceAttrGreaterThan(resourceName, "engines.#", 1),
					resource.TestCheckResourceAttrSet(resourceName, "engines.0.engine"),
					resource.TestCheckResourceAttrSet(resourceName, "engines.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "engines.0.version"),
				),
			},
		},
	})
}

func TestAccDataSourceDatabaseEngines_byEngine(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_database_engines.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataByEngine(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "engines.0.engine", "mysql"),
					resource.TestCheckResourceAttrSet(resourceName, "engines.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "engines.0.version"),
				),
			},
		},
	})
}
