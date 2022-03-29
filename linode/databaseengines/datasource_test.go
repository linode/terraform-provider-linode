package databaseengines_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/databaseengines/tmpl"
)

func TestAccDataSourceDatabaseEngines_all(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_database_engines.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.TestAccProviders,
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
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.TestAccProviders,
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
