package acceptance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
)

var _ statecheck.StateCheck = customStateCheck{}

type stateCheckFunc func(
	context.Context,
	statecheck.CheckStateRequest,
	*statecheck.CheckStateResponse,
)

type customStateCheck struct {
	inner stateCheckFunc
}

func (c customStateCheck) CheckState(
	ctx context.Context,
	req statecheck.CheckStateRequest,
	resp *statecheck.CheckStateResponse,
) {
	c.inner(ctx, req, resp)
}

func CustomStateCheck(inner stateCheckFunc) statecheck.StateCheck {
	return customStateCheck{inner}
}
