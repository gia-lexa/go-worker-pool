package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// worker function that processes jobs and handles potential errors
func worker(id int, jobs <-chan int, results chan<- int, errors chan<- error, wg *sync.WaitGroup) {
	defer wg.Done() // Ensure the WaitGroup counter decreases when the function exits

	for job := range jobs { // Continuously receive jobs from the jobs channel
		fmt.Printf("Worker %d processing job %d\n", id, job)
		time.Sleep(time.Second) // Simulate processing time

		// Simulate a random failure in 30% of the jobs
		if rand.Float32() < 0.3 {
			errors <- fmt.Errorf("worker %d: failed to process job %d", id, job)
			continue // Skip to the next job instead of sending a result
		}

		results <- job * 2 // Send successful job result to the results channel
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Initialize random seed for failure simulation

	numWorkers := 3 // Number of worker goroutines
	numJobs := 10   // Number of jobs to process

	jobs := make(chan int, numJobs)     // Buffered channel to send jobs to workers
	results := make(chan int, numJobs)  // Buffered channel to receive successful results
	errors := make(chan error, numJobs) // Buffered channel for error messages
	var wg sync.WaitGroup               // WaitGroup to track when all workers are done

	// Start worker goroutines
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)                                // Increment WaitGroup counter for each worker
		go worker(i, jobs, results, errors, &wg) // Start worker goroutine
	}

	// Send jobs to the jobs channel
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // Close the jobs channel after all jobs are sent

	// Wait for workers to finish processing
	wg.Wait()
	close(results) // Close results channel after all workers finish
	close(errors)  // Close errors channel after all workers finish

	// Retrieve and print results
	for res := range results {
		fmt.Println("Result:", res)
	}

	// Retrieve and print errors
	for err := range errors {
		fmt.Println("Error:", err)
	}
}
