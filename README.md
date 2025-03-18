# Go Worker Pool with Retries

### 🚀 Overview

This Go project implements a concurrent worker pool that processes jobs efficiently using goroutines and channels. It mirrors what would be necessary for real-world scenarios where tasks need to be executed in parallel, such as:

### 🔹 Real-World Use Cases

- Data Processing Pipelines – Handling large-scale batch jobs like log analysis or ETL workflows.

- Web Scraping – Running multiple concurrent scrapers without overloading a single process.

- API Rate Limiting – Managing controlled, retryable requests to external services.

- Task Queues – Offloading CPU-intensive or asynchronous tasks in a distributed system.

- Load Testing – Simulating high-traffic scenarios by distributing requests across multiple workers.

If a job fails, it automatically retries up to a defined limit before marking it as failed. The number of workers, jobs, and retry attempts are fully configurable via CLI arguments, making this project highly flexible and scalable.

### 🔹 Features

✅ Concurrent Processing – Uses goroutines & channels to process jobs efficiently.

✅ Automatic Retries – Jobs retry on failure up to a max retry limit.

✅ Configurable via CLI – Set the number of workers, jobs, and retries dynamically.

✅ Graceful Shutdown – Ensures all jobs complete safely before exiting.

✅ Unit Tests – Includes robust test cases to validate worker behavior and retries.

✅ GitHub Actions CI – Automatically runs tests on every push.



## Setup & Usage

### 1️⃣ Install & Run

```
git clone https://github.com/yourusername/go-load-tester.git
cd go-load-tester
```

### 2️⃣ Run with default settings
```
go run main.go
```

### 2️⃣ Run with Customizable CLI Arguments

```
go run main.go --workers=5 --jobs=20 --retries=3
```

### 3️⃣ Run Tests

```
go test -v
```


### 🏎️ Benchmarking (Challenges & Next Steps)

I attempted to benchmark the worker pool, but faced unique challenges:

- Retries complicate measurement – If jobs fail and retry, total execution time fluctuates.

- Channel closing issues – Workers sometimes attempted to retry jobs after the jobs channel was closed, causing panics.

- Go's b.N scaling – The benchmark runner dynamically sets b.N, but the presence of retries led to inconsistent performance results.


### Next Steps

To properly benchmark, I could:

- Separate success and failure benchmarking – Measure jobs that complete without retries separately.

- Track execution time per worker – Instead of total time, record each worker’s processing speed.

- Log performance data instead of benchmarking – Capture real-world execution stats without go test -bench constraints.



### 🎯 Future Improvements

🔹 Logging – Save job results & errors to a file for debugging.

🔹 Dynamic Worker Scaling – Auto-adjust the number of workers based on job load.

🔹 Retry Backoff – Use an exponential backoff for failed job retries.

🔹 More Robust Benchmarking – Develop a specialized benchmark without retries interfering.



### 🚀 Let’s Go!
