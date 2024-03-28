package helper

import (
	"fmt"
	"sync"
)

type BatchFunction func() error

// RunBatch is intended to simplify executing functions concurrently.
// This is handy for running certain non-sequential API requests in parallel.
func RunBatch(toExecute ...BatchFunction) error {
	var wg sync.WaitGroup

	errCh := make(chan error)
	defer close(errCh)

	doneCh := make(chan bool)
	defer close(doneCh)

	for _, f := range toExecute {
		// Shadow the function so it can be used in the goroutine
		f := f
		wg.Add(1)
		go func() {
			if err := f(); err != nil {
				errCh <- err
			}
			wg.Done()
		}()
	}

	// Routine to wait for all functions to complete
	go func() {
		wg.Wait()
		doneCh <- true
	}()

	select {
	case <-doneCh:
		return nil
	case err := <-errCh:
		return fmt.Errorf("encountered error when running batch function: %w", err)
	}
}
