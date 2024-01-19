// Package limiter provides a simple way to limit the number of goroutines
// running concurrently. This is useful for controlling resource usage
// in applications that perform a large number of concurrent operations.
package gortlimiter

// Limiter controls the number of goroutines that can run concurrently.
// It uses a channel to limit concurrency.
type Limiter struct {
	slots chan struct{} // Channel used as a semaphore to limit concurrency.
	count int           // The maximum number of concurrent goroutines.
}

// New creates a new Limiter that allows up to `count` goroutines to run concurrently.
func New(count int) *Limiter {
	// Initialize a buffered channel with capacity equal to count.
	// This channel will be used as a semaphore to control concurrency.
	slots := make(chan struct{}, count)

	// Pre-fill the channel with tokens to indicate available slots.
	for i := 0; i < count; i++ {
		slots <- struct{}{}
	}

	return &Limiter{
		slots: slots,
		count: count,
	}
}

// Child represents a single goroutine that has been granted permission
// to run by the Limiter.
type Child struct {
	parent *Limiter // Reference to the parent Limiter.
}

// Get acquires a slot from the Limiter. It blocks until a slot becomes available.
// It returns a Child, which should call End() when the goroutine's work is done.
func (l *Limiter) Get() *Child {
	<-l.slots // Wait for a slot to become available.
	return &Child{parent: l}
}

// End releases the slot back to the Limiter, allowing another goroutine to run.
func (c *Child) End() {
	c.parent.slots <- struct{}{} // Release the slot.
}
