//go:build integration || nfsquota

package nfsquota_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
	"github.com/linode/terraform-provider-linode/v4/linode/nfsquota/tmpl"
)

const resourceName = "linode_nfs_quota.foobar"

func TestAccResourceNFSQuota_basic(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")
	fs := acceptance.GetNFSTestFilesystem(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkNFSQuotaDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, fs.SpaceID, fs.ID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("space_id"), knownvalue.Int64Exact(int64(fs.SpaceID))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("filesystem_id"), knownvalue.Int64Exact(int64(fs.ID))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("path"), knownvalue.StringRegexp(regexp.MustCompile(`^/[^/]+/`+label+`$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("is_root"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("max_capacity_bytes"), knownvalue.Int64Exact(2147483648)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("max_file_count"), knownvalue.Int64Exact(2000)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("status"), knownvalue.NotNull()),
				},
			},
			{
				ResourceName:        resourceName,
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: fmt.Sprintf("%d,%d,", fs.SpaceID, fs.ID),
			},
		},
	})
}

func TestAccResourceNFSQuota_update(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")
	fs := acceptance.GetNFSTestFilesystem(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkNFSQuotaDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, fs.SpaceID, fs.ID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("is_root"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("max_capacity_bytes"), knownvalue.Int64Exact(2147483648)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("max_file_count"), knownvalue.Int64Exact(2000)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("user_group_config"), knownvalue.Null()),
				},
			},
			{
				Config: tmpl.Updates(t, label, fs.SpaceID, fs.ID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("max_capacity_bytes"), knownvalue.Int64Exact(4294967296)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("max_file_count"), knownvalue.Int64Exact(4000)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("user_group_config").AtMapKey("default_user_limit").AtMapKey("max_capacity_bytes"),
						knownvalue.Int64Exact(1073741824),
					),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("user_group_config").AtMapKey("default_user_limit").AtMapKey("max_file_count"),
						knownvalue.Int64Exact(1000),
					),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("user_group_config").AtMapKey("default_group_limit").AtMapKey("max_capacity_bytes"),
						knownvalue.Int64Exact(2147483648),
					),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("user_group_config").AtMapKey("default_group_limit").AtMapKey("max_file_count"),
						knownvalue.Int64Exact(2000),
					),
				},
			},
		},
	})
}

func checkNFSQuotaDestroy(s *terraform.State) error {
	client := acceptance.TestAccSDKv2Provider.Meta().(*helper.ProviderMeta).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_nfs_quota" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error parsing NFS Quota ID %v: %s", rs.Primary.ID, err)
		}

		spaceID, err := strconv.Atoi(rs.Primary.Attributes["space_id"])
		if err != nil {
			return fmt.Errorf("error parsing space_id %v: %s", rs.Primary.Attributes["space_id"], err)
		}

		filesystemID, err := strconv.Atoi(rs.Primary.Attributes["filesystem_id"])
		if err != nil {
			return fmt.Errorf("error parsing filesystem_id %v: %s", rs.Primary.Attributes["filesystem_id"], err)
		}

		_, err = client.GetNFSQuota(context.Background(), spaceID, filesystemID, id)
		if err == nil {
			return fmt.Errorf("NFS Quota with ID %d still exists", id)
		}
		if !linodego.IsNotFound(err) {
			return fmt.Errorf("error getting NFS Quota with ID %d: %s", id, err)
		}
	}

	return nil
}
