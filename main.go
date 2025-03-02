package main

import (
	"flag"      // Package for parsing CLI arguments
	"fmt"       // Package for formatted I/O
	"math/rand" // Package for generating random numbers
	"sync"      // Package for concurrency utilities (WaitGroup, Mutex)
	"time"      // Package for time handling
)

// worker function that processes jobs and retries on failure
func worker(id int, jobs chan int, results chan<- int, errors chan<- error, retries map[int]int, maxRetries int, mutex *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done() // Ensure the WaitGroup counter decreases when the function exits

	for job := range jobs { // Continuously receive jobs from the jobs channel
		fmt.Printf("Worker %d processing job %d (Attempt %d)\n", id, job, retries[job]+1)
		time.Sleep(time.Second) // Simulate processing time

		// Simulate a random failure in 30% of cases
		if rand.Float32() < 0.3 {
			mutex.Lock()                   // Lock before modifying the retries map to avoid race conditions
			retries[job]++                 // Increment the retry counter for this job
			if retries[job] < maxRetries { // Check if job can be retried
				fmt.Printf("Worker %d: Job %d failed, retrying...\n", id, job)
				jobs <- job // Resend job for another attempt
			} else {
				errors <- fmt.Errorf("worker %d: job %d failed after %d attempts", id, job, maxRetries) // Log final failure
			}
			mutex.Unlock() // Unlock to allow other goroutines to modify retries
			continue       // Skip sending a result since the job failed
		}

		results <- job * 2 // Send successful result to the results channel
	}
}

func main() {
	// Define CLI arguments using Go's flag package
	numWorkers := flag.Int("workers", 3, "Number of worker goroutines") // Number of worker goroutines
	numJobs := flag.Int("jobs", 10, "Total number of jobs to process")  // Total jobs to be processed
	maxRetries := flag.Int("retries", 3, "Maximum retries per job")     // Max retry attempts per job before failing

	flag.Parse() // Parse CLI arguments from user input

	rand.Seed(time.Now().UnixNano()) // Initialize random seed for failure simulation

	// Create necessary channels and sync structures
	jobs := make(chan int, *numJobs)     // Buffered channel to send jobs to workers
	results := make(chan int, *numJobs)  // Buffered channel to receive successful results
	errors := make(chan error, *numJobs) // Buffered channel for error messages
	var wg sync.WaitGroup                // WaitGroup to track when all workers are done
	mutex := &sync.Mutex{}               // Mutex to prevent concurrent writes to retries map

	retries := make(map[int]int) // Map to track how many times each job has been retried

	// Start worker goroutines
	for i := 1; i <= *numWorkers; i++ {
		wg.Add(1)                                                             // Increment WaitGroup counter for each worker
		go worker(i, jobs, results, errors, retries, *maxRetries, mutex, &wg) // Start worker goroutine
	}

	// Send jobs to the jobs channel
	for j := 1; j <= *numJobs; j++ {
		jobs <- j // Send job ID into the channel for processing
	}
	close(jobs) // Close the jobs channel after all jobs are sent to prevent deadlocks

	// Wait for workers to finish processing
	wg.Wait()      // Block main function until all workers have finished their tasks
	close(results) // Close results channel after all workers finish
	close(errors)  // Close errors channel after all workers finish

	// Retrieve and print successful results
	for res := range results {
		fmt.Println("Result:", res) // Print each successful result
	}

	// Retrieve and print errors
	for err := range errors {
		fmt.Println("Error:", err) // Print each error message
	}
}
