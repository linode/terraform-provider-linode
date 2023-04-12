package helper

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/linode/linodego"
)

// RunWithStatusRetries runs a function and retries if the function returns a certain
// HTTP status code. If the number of retries exceeds maxRetries, the error will be raised.
func RunWithStatusRetries(toleratedStatuses []int, maxRetries int, retryDelay time.Duration, toRun func() error) error {
	var currentError error
	currentRetry := 0

	isToleratedError := func(err error) bool {
		var linodeError *linodego.Error

		// Handle unexpected errors
		if !errors.As(err, &linodeError) {
			return false
		}

		// Check if the error has a tolerated status
		for _, status := range toleratedStatuses {
			if linodeError.Code == status {
				return true
			}
		}

		return false
	}

	for currentRetry < maxRetries {
		err := toRun()

		// We can return if there isn't an error
		if err == nil {
			return nil
		}

		// Return the error if it's not a tolerated error
		if !isToleratedError(err) {
			return err
		}

		// Retry on tolerated error
		currentError = err
		currentRetry++
		log.Printf("[WARN] Retrying on tolerated error (%d/%d): %s\n", currentRetry, maxRetries, err)

		// Lazy way to delay the next attempt
		time.Sleep(retryDelay)
	}

	return fmt.Errorf("exceeded maximum number of tolerated errors (%d >= %d): %s",
		currentRetry, maxRetries, currentError)
}
