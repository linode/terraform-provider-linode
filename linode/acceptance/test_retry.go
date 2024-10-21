package acceptance

import (
	"testing"
)

// RunTestWithRetries attempts to retry the given test if an intermittent error occurs.
// This function wraps the given testing.T and handles errors accordingly.
// This should only be used for flapping API tests.
func RunTestWithRetries(t *testing.T, maxAttempts int, f func(*testing.T)) {
	for i := 0; i < maxAttempts; i++ {
		if t.Run(t.Name(), f) {
			return
		}

		t.Logf("Retrying %s due to failure. (Attempt %d)", t.Name(), i+1)
	}

	t.Fatalf("Test failed after %d attempts", maxAttempts)
}
