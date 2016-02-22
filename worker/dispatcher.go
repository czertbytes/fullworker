package worker

type Dispatcher struct {
	workerPool chan chan Job
	maxWorkers int
	jobQueue   chan Job
	workers    []Worker
}

func NewDispatcher(jobQueue chan Job, maxWorkers int) *Dispatcher {
	workerPool := make(chan chan Job, maxWorkers)

	return &Dispatcher{
		jobQueue:   jobQueue,
		maxWorkers: maxWorkers,
		workerPool: workerPool,
		workers:    []Worker{},
	}
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
