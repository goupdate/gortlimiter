package gortlimiter

import (
	"sync"
	"testing"
	"time"
)

// TestLimiter ensures that the Limiter correctly handles concurrency.
// It tests whether the limiter allows only a specified number of goroutines
// to run in parallel.
func TestLimiter(t *testing.T) {
	// Define the maximum number of concurrent goroutines.
	const maxConcurrency = 5

	// Initialize the limiter with the defined maximum concurrency.
	l := New(maxConcurrency)

	// Use a WaitGroup to wait for all goroutines to finish.
	var wg sync.WaitGroup

	// Define a function that simulates a task.
	// Each task prints its ID and then sleeps for a time proportional to its ID.
	// This variability in task duration helps test the limiter's behavior under different load conditions.
	walk := func(i int) {
		defer wg.Done() // Ensure that the WaitGroup counter is decremented when the task finishes.

		c := l.Get()   // Acquire a slot from the limiter.
		defer c.End() // Release the slot when the task is done.

		// Simulate a task that takes varying amounts of time to complete.
		time.Sleep(time.Millisecond * 500 * time.Duration(i))
	}

	// Start multiple goroutines more than the limit set.
	for i := 0; i < 10; i++ {
		wg.Add(1) // Increment the WaitGroup counter before starting a goroutine.
		go walk(i)
	}

	// Wait for all goroutines to finish.
	wg.Wait()

	// Optionally, add a brief sleep at the end to ensure all resources are cleaned up.
	// This is not strictly necessary but can be useful in more complex tests.
	time.Sleep(time.Second)
}
