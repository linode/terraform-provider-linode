package acceptance

import (
	"fmt"

	"github.com/linode/linodego/v2"
)

func checkDatabaseDestroyResult(id int, err error) error {
	if err == nil {
		return fmt.Errorf("should not find database ID %d existing after delete", id)
	}
	if !linodego.IsNotFound(err) {
		return fmt.Errorf("failed to get database ID %d: %w", id, err)
	}
	return nil
}
