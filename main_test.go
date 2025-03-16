package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// Custom worker function that prevents panics when channels are closed
func testWorker(id int, jobs chan int, results chan<- int, errors chan<- error, retries map[int]int, maxRetries int, mutex *sync.Mutex, wg *sync.WaitGroup, done chan struct{}) {
	defer wg.Done() // Ensure we signal WaitGroup when done

	for {
		select {
		case job, open := <-jobs: // Read from jobs channel
			if !open {
				return // Exit if jobs channel is closed
			}

			fmt.Printf("Worker %d processing job %d (Attempt %d)\n", id, job, retries[job]+1)
			time.Sleep(time.Millisecond * 100) // Shorter sleep for faster testing

			mutex.Lock()
			retries[job]++

			if retries[job] < maxRetries { // If not exceeded max retries, retry job
				fmt.Printf("Worker %d: Job %d failed, retrying...\n", id, job)

				// **Check if the jobs channel is closed before retrying**
				select {
				case <-done: // Stop retrying if workers are shutting down
					fmt.Printf("Worker %d: Job %d could not be retried because channel is closed\n", id, job)
				case jobs <- job: // Only retry if channel is open
				}
			} else {
				errors <- fmt.Errorf("worker %d: job %d failed after %d attempts", id, job, maxRetries)
			}
			mutex.Unlock()
		case <-done: // Stop worker if the shutdown signal is received
			return
		}
	}
}

// TestMaxRetries ensures jobs fail after exceeding retry limit
func TestMaxRetries(t *testing.T) {
	jobs := make(chan int, 1)
	results := make(chan int, 1)
	errors := make(chan error, 1)
	done := make(chan struct{}) // Signal to safely close channels

	var wg sync.WaitGroup
	mutex := &sync.Mutex{}
	retries := make(map[int]int)
	maxRetries := 2

	// Start the test worker
	wg.Add(1)
	go testWorker(1, jobs, results, errors, retries, maxRetries, mutex, &wg, done)

	jobs <- 1 // Send a job that will fail

	// **Wait before closing jobs to allow retries**
	time.Sleep(time.Millisecond * 300)
	close(done) // Signal workers to stop processing
	close(jobs) // Close the jobs channel only after workers stop

	wg.Wait()     // Wait for worker to finish processing
	close(errors) // Close errors channel

	// **Ensure an error is received after max retries**
	select {
	case err := <-errors:
		if err == nil {
			t.Error("Expected error, but job succeeded") // Fail test if no error is reported
		}
	case <-time.After(time.Second * 2): // Prevent hanging
		t.Error("Worker did not process job in time")
	}
}
