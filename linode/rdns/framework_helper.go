package rdns

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/linode/linodego"
)

// updateIPAddress wraps the client.UpdateIPAddress(...) retry logic depending on the 'wait_for_available' field.
func updateIPAddress(
	ctx context.Context,
	client *linodego.Client,
	address string,
	desiredRDNS *string,
	waitForAvailable bool,
) (*linodego.InstanceIP, error) {
	updateOpts := linodego.IPAddressUpdateOptions{
		RDNS: desiredRDNS,
	}

	if waitForAvailable {
		return updateIPAddressWithRetries(ctx, client, address, updateOpts, time.Second*5)
	}

	return client.UpdateIPAddress(ctx, address, updateOpts)
}

func updateIPAddressWithRetries(ctx context.Context, client *linodego.Client, address string,
	updateOpts linodego.IPAddressUpdateOptions, retryDuration time.Duration,
) (*linodego.InstanceIP, error) {
	ticker := time.NewTicker(retryDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			result, err := client.UpdateIPAddress(ctx, address, updateOpts)
			if err != nil {
				if lerr, ok := err.(*linodego.Error); ok && lerr.Code != 400 &&
					!strings.Contains(lerr.Error(), "unable to perform a lookup") {
					return nil, fmt.Errorf("failed to update ip address: %s", err)
				}

				continue
			}

			return result, nil

		case <-ctx.Done():
			// The timeout for this context will implicitly be handled by Terraform
			return nil, fmt.Errorf("failed to update ip address: %s", ctx.Err())
		}
	}
}
