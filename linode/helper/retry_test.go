package helper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func TestRunWithStatusRetries_basic(t *testing.T) {
	numAttempts := 0

	err := helper.RunWithStatusRetries([]int{502}, 10, time.Microsecond, func() error {
		numAttempts++

		if numAttempts <= 2 {
			return fmt.Errorf("oh no %w", &linodego.Error{Code: 502})
		}

		return nil
	})
	if err != nil {
		t.Fatalf("got unexpected error: %s", err)
	}

	if numAttempts != 3 {
		t.Fatalf("got unexpected number of attempts: %d", numAttempts)
	}
}

func TestRunWithStatusRetries_exceedMaxRetries(t *testing.T) {
	err := helper.RunWithStatusRetries([]int{502}, 3, time.Microsecond, func() error {
		return fmt.Errorf("oh no %w", &linodego.Error{Code: 502})
	})

	if err == nil {
		t.Fatalf("missing error")
	}
}

func TestRunWithStatusRetries_irrelevantStatus(t *testing.T) {
	numAttempts := 0

	err := helper.RunWithStatusRetries([]int{502}, 3, time.Microsecond, func() error {
		numAttempts++
		return fmt.Errorf("oh no %w", &linodego.Error{Code: 503})
	})

	if err == nil {
		t.Fatalf("missing error")
	}

	if numAttempts != 1 {
		t.Fatalf("got unexpected number of attempts: %d", numAttempts)
	}
}

func TestRunWithStatusRetries_irrelevantError(t *testing.T) {
	numAttempts := 0

	err := helper.RunWithStatusRetries([]int{502}, 3, time.Microsecond, func() error {
		numAttempts++
		return fmt.Errorf("oh no")
	})

	if err == nil {
		t.Fatalf("missing error")
	}

	if numAttempts != 1 {
		t.Fatalf("got unexpected number of attempts: %d", numAttempts)
	}
}
