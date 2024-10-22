package acceptance

import (
	"fmt"
	"sync/atomic"
	"testing"
)

// RunTestWithRetries attempts to retry the given test if an intermittent error occurs.
// This function wraps the given testing.T and handles errors accordingly.
// This should only be used for flapping API tests.
func RunTestWithRetries(t testing.TB, maxAttempts int, f func(t *WrappedT)) {
	for i := 0; i < maxAttempts; i++ {
		wrappedT := &WrappedT{
			TB:     t,
			failed: atomic.Bool{},
		}

		closurePanic := false

		// Run the retryable test closure,
		// capturing any test failures
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("panic: %v", r)
					closurePanic = true
				}
			}()

			f(wrappedT)
		}()

		if !closurePanic && !t.Failed() {
			return
		}

		t.Logf("Retrying %s due to failure. (Attempt %d)", t.Name(), i+1)
	}

	t.Fatalf("Test failed after %d attempts", maxAttempts)
}

// WrappedT implements testing.T with additional retry logic
type WrappedT struct {
	testing.TB

	failed atomic.Bool
}

func (t *WrappedT) Failed() bool {
	return t.failed.Load()
}

func (t *WrappedT) Fail() {
	t.failed.Store(true)
}

func (t *WrappedT) FailNow() {
	t.Fail()
	panic(nil)
}

func (t *WrappedT) Error(args ...any) {
	t.TB.Log("ERROR:", fmt.Sprint(args...))
	t.Fail()
}

func (t *WrappedT) Errorf(format string, args ...any) {
	t.TB.Log("ERROR:", fmt.Sprintf(format, args...))
	t.Fail()
}

func (t *WrappedT) Fatal(args ...any) {
	t.TB.Log("FATAL:", fmt.Sprint(args...))
	t.FailNow()
}

func (t *WrappedT) Fatalf(format string, args ...any) {
	t.TB.Log("FATAL:", fmt.Sprintf(format, args...))
	t.FailNow()
}

func (t *WrappedT) Parallel() {
	t.TB.Fatal("Parallel() cannot be called on instances of WrappedT")
}
