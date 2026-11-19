package helper

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego/v2"
)

const (
	// NFSAccessPolicyTimeout bounds how long to wait for an NFS access policy
	// to finish a pending operation.
	NFSAccessPolicyTimeout = 15 * time.Minute

	nfsAccessPolicyPollInterval = 5 * time.Second
)

// WaitForNFSAccessPolicySettled polls an NFS access policy until it leaves a
// transitional status (creating or updating). Unlike linodego's WaitFor helpers,
// it does not require the status to be exactly "active", so it cannot hang on
// statuses the API reports for unconfigured or disabled policies.
func WaitForNFSAccessPolicySettled[T any](
	ctx context.Context,
	timeout time.Duration,
	get func(ctx context.Context) (*T, linodego.NFSAccessPolicyStatus, error),
) (*T, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(nfsAccessPolicyPollInterval)
	defer ticker.Stop()

	var lastStatus linodego.NFSAccessPolicyStatus

	for {
		policy, status, err := get(ctx)
		if err != nil {
			return nil, err
		}

		lastStatus = status

		switch status {
		case linodego.NFSAccessPolicyStatusCreating, linodego.NFSAccessPolicyStatusUpdating:
			tflog.Debug(ctx, "Waiting for NFS access policy to settle", map[string]any{
				"status": status,
			})
		case linodego.NFSAccessPolicyStatusError:
			return policy, fmt.Errorf("access policy reached status %q", status)
		default:
			return policy, nil
		}

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf(
				"timed out waiting for access policy to leave status %q: %w", lastStatus, ctx.Err(),
			)
		case <-ticker.C:
		}
	}
}

// isNFSAccessPolicyNotReadyError reports whether err is the 400 the API returns
// when the parent space or filesystem is not yet ready for an access policy
// update, even though its status already reports active.
func isNFSAccessPolicyNotReadyError(err error) bool {
	return linodego.ErrHasStatus(err, http.StatusBadRequest) &&
		strings.Contains(strings.ToLower(err.Error()), "not ready")
}

// RetryNFSAccessPolicyUpdate calls update, retrying while the API reports that
// the resource is not yet ready for an access policy update.
func RetryNFSAccessPolicyUpdate(
	ctx context.Context,
	timeout time.Duration,
	update func(ctx context.Context) error,
) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(nfsAccessPolicyPollInterval)
	defer ticker.Stop()

	for {
		err := update(ctx)
		if err == nil || !isNFSAccessPolicyNotReadyError(err) {
			return err
		}

		tflog.Debug(ctx, "NFS access policy is not ready for update; retrying", map[string]any{
			"error": err.Error(),
		})

		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out retrying access policy update: %w", err)
		case <-ticker.C:
		}
	}
}
