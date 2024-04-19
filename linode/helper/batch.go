package helper

import (
	"errors"
	"sync"
)

type BatchFunction func() error

// RunBatch is intended to simplify executing functions concurrently.
// This is handy for running certain non-sequential API requests in parallel.
func RunBatch(toExecute ...BatchFunction) error {
	var wg sync.WaitGroup
	wg.Add(len(toExecute))

	errCh := make(chan error)
	doneCh := make(chan bool)

	for _, f := range toExecute {
		// Shadow the function so it can be used in the goroutine
		f := f
		go func() {
			defer wg.Done()
			if err := f(); err != nil {
				errCh <- err
			}
		}()
	}

	// Routine to wait for all functions to complete
	go func() {
		wg.Wait()
		doneCh <- true
		close(doneCh)
		close(errCh)
	}()

	allErrors := []error{}

	for {
		select {
		case <-doneCh:
			return errors.Join(allErrors...)
		case err := <-errCh:
			allErrors = append(allErrors, err)
		}
	}
}
