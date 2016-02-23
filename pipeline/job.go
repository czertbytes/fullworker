package pipeline

import "time"

type Payload struct {
	URL     string
	Content string
	Time    time.Duration
}

type JobResult struct {
	Payload
	Length int
}
