package main

import (
	"fmt"
	"sync"

	"github.com/czertbytes/fullworker/worker"
)

func main() {
	workerQueue()
}

func workerQueue() {

	urls := []string{
		"https://www.google.com",
		"https://www.apple.com",
		"https://www.golang.org",
		"https://www.microsoft.com",
		"https://www.yahoo.com",
		"https://www.amazon.com",
		"https://www.netflix.com",
		"https://www.spotify.com",
		"https://www.dropbox.com",
	}

	jobQueue := make(chan worker.Job, 10)

	dispatcher := worker.NewDispatcher(jobQueue, 5)
	dispatcher.Run()

	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		newJob := worker.Job{
			URL:        url,
			ResultChan: make(chan worker.JobResult),
		}

		jobQueue <- newJob

		go func(job worker.Job) {
			defer wg.Done()

			result := <-job.ResultChan

			fmt.Printf("%s downloaded %d bytes in %s\n", job.URL, result.Length, result.Time)
		}(newJob)
	}

	wg.Wait()
}
