package worker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Worker struct {
	id         int
	jobQueue   chan Job
	workerPool chan chan Job
	quitChan   chan bool
}

func NewWorker(id int, workerPool chan chan Job) Worker {
	return Worker{
		id:         id,
		jobQueue:   make(chan Job),
		workerPool: workerPool,
		quitChan:   make(chan bool),
	}
}

func (w Worker) start() {
	go func() {
		for {
			w.workerPool <- w.jobQueue

			select {
			case job := <-w.jobQueue:
				start := time.Now()
				fmt.Printf("worker%d: started %s\n", w.id, job.URL)

				resp, err := http.Get(job.URL)
				if err != nil {
					return
				}
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)

				fmt.Printf("worker%d: completed %s\n", w.id, job.URL)

				job.ResultChan <- JobResult{
					Time:   time.Now().Sub(start),
					Length: len(body),
				}
			case <-w.quitChan:
				fmt.Printf("worker%d stopping\n", w.id)
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.quitChan <- true
	}()
}
