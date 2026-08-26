//go:build integration || vpcsubnet

package vpcsubnet_test

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

// waitForVPCSubnetNodebalancer polls the Linode API until the VPC subnet reports
// at least wantCount nodebalancers, or timeout is reached.
func waitForVPCSubnetNodebalancer(resName string, wantCount int, timeout time.Duration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resName)
		}

		subnetID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error parsing subnet ID %q: %w", rs.Primary.ID, err)
		}

		vpcID, err := strconv.Atoi(rs.Primary.Attributes["vpc_id"])
		if err != nil {
			return fmt.Errorf("error parsing vpc_id: %w", err)
		}

		client := acceptance.TestAccSDKv2Provider.Meta().(*helper.ProviderMeta).Client

		deadline := time.Now().Add(timeout)
		for time.Now().Before(deadline) {
			subnet, err := client.GetVPCSubnet(context.Background(), vpcID, subnetID)
			if err != nil {
				return fmt.Errorf("error fetching VPC subnet %d: %w", subnetID, err)
			}
			if len(subnet.Nodebalancers) >= wantCount {
				return nil
			}
			time.Sleep(5 * time.Second)
		}

		return fmt.Errorf("timed out waiting for VPC subnet %d to have %d nodebalancer(s)", subnetID, wantCount)
	}
}
