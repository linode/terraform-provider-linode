package vlan_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/vlan/tmpl"
)

func preConfigVLANPoll(t *testing.T, vlanName string) func() {
	return func() {
		client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
		if _, err := waitForVLANWithLabel(client, vlanName, 30); err != nil {
			t.Fatal(err)
		}
	}
}

func TestAccDataSourceVLANs_basic(t *testing.T) {
	t.Parallel()

	instanceName := acctest.RandomWithPrefix("tf_test")
	vlanName := "tf-test"
	resourceName := "data.linode_vlans.foolan"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, instanceName, vlanName),
			},
			{
				PreConfig: preConfigVLANPoll(t, vlanName),
				Config:    tmpl.DataBasic(t, instanceName, vlanName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "vlans.0.label", vlanName),
					resource.TestCheckResourceAttr(resourceName, "vlans.0.region", "us-southeast"),
					resource.TestCheckResourceAttrSet(resourceName, "vlans.0.created"),
					resource.TestCheckResourceAttrSet(resourceName, "vlans.0.linodes.#"),
				),
			},
		},
	})
}

func TestAccDataSourceVLANs_regex(t *testing.T) {
	t.Parallel()

	instanceName := acctest.RandomWithPrefix("tf_test")
	vlanName := "tf-test"
	resourceName := "data.linode_vlans.foolan"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataRegex(t, instanceName, vlanName),
			},
			{
				PreConfig: preConfigVLANPoll(t, vlanName),
				Config:    tmpl.DataRegex(t, instanceName, vlanName),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckResourceAttrGreaterThan(resourceName, "vlans.#", 0),
					resource.TestCheckResourceAttr(resourceName, "vlans.0.label", vlanName+"-new"),
					resource.TestCheckResourceAttr(resourceName, "vlans.0.region", "us-southeast"),
					resource.TestCheckResourceAttrSet(resourceName, "vlans.0.created"),
					resource.TestCheckResourceAttrSet(resourceName, "vlans.0.linodes.#"),
				),
			},
		},
	})
}

// This test is necessary to test a race-condition introduced by provisioning
// multiple concurrent VLAN instances.
// This test is opt-in as it has the potential to spawn a large number of VLANs,
// which cannot be deleted without admin intervention.
func TestAccDataSourceVLANs_ensureNoDuplicates(t *testing.T) {
	acceptance.OptInTest(t)

	t.Parallel()

	instanceName := acctest.RandomWithPrefix("tf_test")
	vlanName := acctest.RandomWithPrefix("tf-test")
	resourceName := "data.linode_vlans.foolan"

	createValidateSteps := func(i int) []resource.TestStep {
		vlanName := fmt.Sprintf("%s-%d", vlanName, i)

		return []resource.TestStep{
			{
				Config: tmpl.DataCheckDuplicate(t, instanceName, vlanName),
			},
			{
				PreConfig: preConfigVLANPoll(t, vlanName),
				Config:    tmpl.DataCheckDuplicate(t, instanceName, vlanName),
				Check: resource.ComposeTestCheckFunc(
					// Ensure only one VLAN is created
					resource.TestCheckResourceAttr(resourceName, "vlans.#", "1"),
				),
			},
		}
	}

	var steps []resource.TestStep

	// Run this test multiple times to test on updates
	for i := 0; i < 3; i++ {
		steps = append(steps, createValidateSteps(i)...)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps:     steps,
	})
}

// waitForVLANWithLabel polls for a VLAN with the given label to exist
// This is necessary in this context because it is not guaranteed that a VLAN will
// be created immediately after the instance is created.
func waitForVLANWithLabel(client linodego.Client, label string, timeoutSeconds int) (*linodego.VLAN, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeoutSeconds))
	defer cancel()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			vlans, err := client.ListVLANs(ctx,
				&linodego.ListOptions{Filter: fmt.Sprintf("{\"label\": \"%s\"}", label)})
			if err != nil {
				return nil, err
			}

			if len(vlans) > 0 {
				return &vlans[0], err
			}

		case <-ctx.Done():
			return nil, fmt.Errorf("Error waiting for VLAN %s: %s", label, ctx.Err())
		}
	}
}
