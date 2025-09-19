package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func WaitForCondition(
	ctx context.Context,
	interval time.Duration,
	checkFunc func(ctx context.Context) (bool, error),
) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tflog.Debug(ctx, "Running condition...")
			ok, err := checkFunc(ctx)
			if err != nil {
				return fmt.Errorf("check function had error: %w", err)
			}

			if ok {
				return nil
			}
		case <-ctx.Done():
			return fmt.Errorf("failed to wait for condition: %w", ctx.Err())
		}
	}
}
