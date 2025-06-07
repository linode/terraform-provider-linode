//go:build integration || firewalltemplates

package firewalltemplates_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/firewalltemplates/tmpl"
)

const testTemplatesDataName = "data.linode_firewall_templates.test"

func TestAccDataSourceFirewalls_basic(t *testing.T) {
	t.Parallel()
	testSlug := "public"

	acceptance.RunTestWithRetries(t, 3, func(t *acceptance.WrappedT) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: tmpl.DataBasic(t),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(testTemplatesDataName, tfjsonpath.New("id"), knownvalue.NotNull()),
						statecheck.ExpectKnownValue(testTemplatesDataName, tfjsonpath.New("firewall_templates"), knownvalue.NotNull()),
					},
				},
				{
					Config: tmpl.DataFilter(t, testSlug),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(testTemplatesDataName, tfjsonpath.New("id"), knownvalue.NotNull()),
						statecheck.ExpectKnownValue(testTemplatesDataName, tfjsonpath.New("firewall_templates").AtSliceIndex(0).AtMapKey("slug"), knownvalue.StringExact(testSlug)),
						statecheck.ExpectKnownValue(testTemplatesDataName, tfjsonpath.New("firewall_templates").AtSliceIndex(0).AtMapKey("inbound"), knownvalue.NotNull()),
						statecheck.ExpectKnownValue(testTemplatesDataName, tfjsonpath.New("firewall_templates").AtSliceIndex(0).AtMapKey("inbound_policy"), knownvalue.NotNull()),
						statecheck.ExpectKnownValue(testTemplatesDataName, tfjsonpath.New("firewall_templates").AtSliceIndex(0).AtMapKey("outbound"), knownvalue.NotNull()),
						statecheck.ExpectKnownValue(testTemplatesDataName, tfjsonpath.New("firewall_templates").AtSliceIndex(0).AtMapKey("outbound_policy"), knownvalue.NotNull()),
					},
				},
			},
		})
	})
}
