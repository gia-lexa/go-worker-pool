package main

import (
	"fmt"
	"sync"
	"time"
)

// worker function processes jobs from the jobs channel
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done() // Ensure the WaitGroup counter decreases when the function exits

	for job := range jobs { // Continuously receive jobs from the jobs channel
		fmt.Printf("Worker %d processing job %d\n", id, job)
		time.Sleep(time.Second) // Simulate some work by sleeping for 1 second
		results <- job * 2      // Send the processed result to the results channel
	}
}

func main() {
	numWorkers := 3 // Number of worker goroutines
	numJobs := 5    // Number of jobs to process

	jobs := make(chan int, numJobs)    // Buffered channel to send jobs to workers
	results := make(chan int, numJobs) // Buffered channel to receive results from workers
	var wg sync.WaitGroup              // WaitGroup to ensure all workers finish before main exits

	// Start worker goroutines
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)                        // Increment WaitGroup counter
		go worker(i, jobs, results, &wg) // Start worker goroutine
	}

	// Send jobs to the jobs channel
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // Close the jobs channel to indicate no more jobs will be sent

	// Wait for all worker goroutines to finish
	wg.Wait()
	close(results) // Close results channel after all workers finish processing

	// Retrieve and print results
	for res := range results {
		fmt.Println("Result:", res)
	}
}
