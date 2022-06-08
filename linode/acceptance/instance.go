package acceptance

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func CheckInstanceExists(name string, instance *linodego.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := TestAccProvider.Meta().(*helper.ProviderMeta).Client

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		found, err := client.GetInstance(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of Instance %s: %s", rs.Primary.Attributes["label"], err)
		}

		*instance = *found

		return nil
	}
}

func CheckInstanceDestroy(s *terraform.State) error {
	client := TestAccProvider.Meta().(*helper.ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_instance" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v as int", rs.Primary.ID)
		}

		if id == 0 {
			return fmt.Errorf("should not have Linode ID 0")
		}

		_, err = client.GetInstance(context.Background(), id)

		if err == nil {
			return fmt.Errorf("should not find Linode ID %d existing after delete", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error getting Linode ID %d: %s", id, err)
		}
	}

	return nil
}

func AssertInstanceReboot(t *testing.T, shouldRestart bool, instance *linodego.Instance) func() {
	return func() {
		client := TestAccProvider.Meta().(*helper.ProviderMeta).Client
		eventFilter := fmt.Sprintf(
			`{"entity.type": "linode", "entity.id": %d, "action": "linode_reboot", "created": { "+gte": "%s" }}`,
			instance.ID, instance.Created.Format("2006-01-02T15:04:05"))

		events, err := client.ListEvents(context.Background(), &linodego.ListOptions{Filter: eventFilter})
		if err != nil {
			t.Fail()
		}

		if len(events) == 0 && shouldRestart {
			t.Fatal("expected instance to have been rebooted")
		}

		if len(events) > 0 && !shouldRestart {
			t.Fatal("expected instance to not have been rebooted")
		}
	}
}
