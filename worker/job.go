package worker

import "time"

type Job struct {
	URL        string
	ResultChan chan JobResult
}

type JobResult struct {
	Time   time.Duration
	Length int
}
