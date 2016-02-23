package worker

import (
	"fmt"
	"sync"
)

type Dispatcher struct {
	workerPool chan chan Job
	maxWorkers int
	jobQueue   chan Job
	workers    []Worker
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	workerPool := make(chan chan Job, maxWorkers)

	return &Dispatcher{
		jobQueue:   make(chan Job, 10),
		maxWorkers: maxWorkers,
		workerPool: workerPool,
		workers:    []Worker{},
	}
}

func (d *Dispatcher) AddJob(wg *sync.WaitGroup, url string) {
	newJob := Job{
		URL:        url,
		ResultChan: make(chan JobResult),
	}

	d.jobQueue <- newJob

	go func(job Job) {
		defer wg.Done()
		result := <-job.ResultChan

		fmt.Printf("%s downloaded %d bytes in %s\n", job.URL, result.Length, result.Time)
	}(newJob)
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(i+1, d.workerPool)
		worker.start()

		d.workers = append(d.workers, worker)
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.jobQueue:
			go func() {
				workerJobQueue := <-d.workerPool
				workerJobQueue <- job
			}()
		}
	}
}

func (d *Dispatcher) Stop() {
	for _, worker := range d.workers {
		worker.Stop()
	}
}
