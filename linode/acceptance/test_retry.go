package acceptance

import (
	"fmt"
	"runtime"
	"testing"
)

// RunTestRetry attempts to retry the given test if an intermittent error occurs.
// This function wraps the given testing.T and handles errors accordingly.
// This should only be used for flapping API tests.
func RunTestRetry(t *testing.T, maxAttempts int, f func(t *TRetry)) {
	for i := 0; i < maxAttempts; i++ {
		newT := NewTRetry(t)
		t.Cleanup(newT.Close)

		go func() {
			f(newT)

			newT.SuccessChannel <- true
		}()

		select {
		case err := <-newT.ErrorChannel:
			t.Logf("Retrying on test failure: %s\n", err)
		case <-newT.SuccessChannel:
			return
		}
	}

	t.Fatalf("Test failed after %d attempts", maxAttempts)
}

// TRetry implements testing.T with additional retry logic
type TRetry struct {
	t *testing.T

	ErrorChannel   chan error
	SuccessChannel chan bool
}

func NewTRetry(t *testing.T) *TRetry {
	return &TRetry{
		t:              t,
		ErrorChannel:   make(chan error, 0),
		SuccessChannel: make(chan bool, 0),
	}
}

func (t *TRetry) Close() {
	close(t.ErrorChannel)
	close(t.SuccessChannel)
}

func (t *TRetry) Cleanup(f func()) {
	t.t.Cleanup(f)
}

func (t *TRetry) Error(args ...any) {
	t.ErrorChannel <- fmt.Errorf(fmt.Sprint(args...))
}

func (t *TRetry) Errorf(format string, args ...any) {
	t.ErrorChannel <- fmt.Errorf(format, args...)
}

func (t *TRetry) Fail() {
	runtime.Goexit()
}

func (t *TRetry) Failed() bool {
	// We wrap this logic using channels
	return false
}

func (t *TRetry) Fatal(args ...any) {
	t.ErrorChannel <- fmt.Errorf(fmt.Sprint(fmt.Sprint(args...)))
	t.Fail()
}

func (t *TRetry) Fatalf(format string, args ...any) {
	t.ErrorChannel <- fmt.Errorf(format, args...)
	t.Fail()
}

func (t *TRetry) Helper() {
	t.t.Helper()
}

func (t *TRetry) Log(args ...any) {
	t.t.Log(args...)
}

func (t *TRetry) Logf(format string, args ...any) {
	t.t.Logf(format, args...)
}

func (t *TRetry) Name() string {
	return t.t.Name()
}

func (t *TRetry) Setenv(key, value string) {
	t.t.Setenv(key, value)
}

func (t *TRetry) Skip(args ...any) {
	t.t.Skip(args...)
}

func (t *TRetry) SkipNow() {
	t.t.SkipNow()
}

func (t *TRetry) Skipf(format string, args ...any) {
	t.t.Skipf(format, args...)
}

func (t *TRetry) Skipped() bool {
	return t.t.Skipped()
}

func (t *TRetry) TempDir() string {
	return t.t.TempDir()
}

func (t *TRetry) FailNow() {
	t.Fail()
}

func (t *TRetry) Parallel() {
	t.t.Parallel()
}

//lint:ignore U1000 Ignore unused function
func (t *TRetry) private() {}
