//go:build integration || lkeversion

package lkeversion_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/lkeversion/tmpl"
)

func TestAccDataSourceLinodeLkeVersion_NoTier(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_lke_version.foobar"

	client, err := acceptance.GetTestClient()
	if err != nil {
		t.Fatal(err)
	}

	// Resolve an LKE version
	versions, err := client.ListLKEVersions(context.Background(), nil)
	if err != nil {
		t.Fatalf("failed to list lke versions: %s", err)
	}

	version := versions[0]

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataNoTier(t, version.ID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", version.ID),
					resource.TestCheckNoResourceAttr(resourceName, "tier"),
				),
			},
		},
	})
}

func TestAccDataSourceLinodeLkeVersion_Tier(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_lke_version.foobar"

	client, err := acceptance.GetTestClient()
	if err != nil {
		t.Fatal(err)
	}

	tier := "enterprise"

	// Resolve an LKE version
	versions, err := client.ListLKETierVersions(context.Background(), tier, nil)
	if err != nil {
		t.Fatalf("failed to list lke versions: %s", err)
	}

	version := versions[0]

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataTier(t, version.ID, tier),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", version.ID),
					resource.TestCheckResourceAttr(resourceName, "tier", tier),
				),
			},
		},
	})
}
