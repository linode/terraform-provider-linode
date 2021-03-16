package linode

import "sync"

// waitGroupCh creates a new readonly struct channel that is signaled when
// the underlying sync.WaitGroup channel reaches 0.
func waitGroupCh(wg *sync.WaitGroup) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()
	return done
}
