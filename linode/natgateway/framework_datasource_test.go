//go:build integration || natgateway

package natgateway_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/natgateway/tmpl"
)

const natDataSource = "data.linode_natgateway.foobar"

// TestAccDataSourceNATGateway_basic reads a NAT Gateway attached to a VPC subnet and checks
// that every field is populated.
func TestAccDataSourceNATGateway_basic(t *testing.T) {
	t.Parallel()

	region, err := acceptance.GetRandomRegionWithCaps(
		[]linodego.RegionCapability{linodego.CapabilityVPCs}, "core")
	if err != nil {
		t.Fatal(err)
	}

	// The label must begin with a letter.
	label := "tf-" + acctest.RandString(10)

	var checks []statecheck.StateCheck

	topLevel := []string{
		"id",
		"region",
		"label",
		"default_ports_per_interface",
		"addresses",
		"address_autoscale_max",
		"portset_assignments",
		"portset_capacity",
		"vpc_subnet",
		"created",
		"updated",
	}
	for _, field := range topLevel {
		checks = append(checks, statecheck.ExpectKnownValue(
			natDataSource, tfjsonpath.New(field), knownvalue.NotNull()))
	}

	subnetFields := []string{"id", "type", "url", "label", "vpc_id", "vpc_label"}
	for _, field := range subnetFields {
		checks = append(checks, statecheck.ExpectKnownValue(
			natDataSource, tfjsonpath.New("vpc_subnet").AtMapKey(field), knownvalue.NotNull()))
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:            tmpl.DataBasic(t, label, region),
				ConfigStateChecks: checks,
			},
		},
	})
}
