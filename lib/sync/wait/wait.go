package wait

import (
	"sync"
	"time"
)

// Wait is similar with sync.WaitGroup, but it can wait with timeout
type Wait struct {
	wg sync.WaitGroup
}

// Add adds delta, which may be negative, to the WaitGroup counter.
func (w *Wait) Add(delta int) {
	w.wg.Add(delta)
}

// Done decrements the WaitGroup counter by one
func (w *Wait) Done() {
	w.wg.Done()
}

// Wait blocks until the WaitGroup counter is zero
func (w *Wait) Wait() {
	w.wg.Wait()
}

// WaitWithTimeout blocks until the WaitGroup counter is zero or timeout
// return true if timeout
func (w *Wait) WaitWithTimeout(timeout int) bool {
	c := make(chan bool, 1)
	go func() {
		w.wg.Wait()
		c <- true
		defer close(c)

	}()
	select {
	case <-c:
		return false
	case <-time.After(time.Duration(timeout)):
		return true
	}
}
