package rdns

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
)

// updateIPAddress wraps the client.UpdateIPAddress(...) retry logic depending on the 'wait_for_available' field.
func updateIPAddress(
	ctx context.Context,
	client *linodego.Client,
	address string,
	opts linodego.IPAddressUpdateOptions,
	waitForAvailable bool,
) (*linodego.InstanceIP, error) {
	if waitForAvailable {
		return updateIPAddressWithRetries(ctx, client, address, opts, time.Second*5)
	}

	tflog.Debug(ctx, "client.UpdateIPAddress(...)", map[string]any{
		"options": opts,
	})

	return client.UpdateIPAddress(ctx, address, opts)
}

func updateIPAddressWithRetries(ctx context.Context, client *linodego.Client, address string,
	updateOpts linodego.IPAddressUpdateOptions, retryDuration time.Duration,
) (*linodego.InstanceIP, error) {
	tflog.Debug(ctx, "Attempting to update IP address RDNS with retries")

	ticker := time.NewTicker(retryDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tflog.Debug(ctx, "client.UpdateIPAddress(...)", map[string]any{
				"options": updateOpts,
			})
			result, err := client.UpdateIPAddress(ctx, address, updateOpts)
			if err != nil {
				if lerr, ok := err.(*linodego.Error); ok && lerr.Code != 400 &&
					!strings.Contains(lerr.Error(), "unable to perform a lookup") {
					return nil, fmt.Errorf("failed to update ip address: %s", err)
				}

				tflog.Debug(ctx, "IP is not yet ready for assignment")
				continue
			}

			return result, nil

		case <-ctx.Done():
			// The timeout for this context will implicitly be handled by Terraform
			return nil, fmt.Errorf("failed to update ip address: %s", ctx.Err())
		}
	}
}
