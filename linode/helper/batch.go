package helper

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type BatchFunction func(ctx context.Context) error

// RunBatch is intended to simplify executing functions concurrently.
// This is handy for running certain non-sequential API requests in parallel.
// NOTE: This should NOT be used until linodego has been confirmed to be thread-safe.
func RunBatch(ctx context.Context, toExecute ...BatchFunction) error {
	eg, ctx := errgroup.WithContext(ctx)

	for _, f := range toExecute {
		// Shadow the function so it can be used in the goroutine
		f := f
		eg.Go(func() error { return f(ctx) })
	}

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
