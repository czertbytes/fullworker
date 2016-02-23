package pipeline

import "fmt"

type sizer struct {
	process
	In  chan Payload
	Out chan JobResult
}

func NewSizer() *sizer {
	return &sizer{
		In:  make(chan Payload, jobBuffer),
		Out: make(chan JobResult, jobBuffer),
	}
}

func (p *sizer) Run() {
	defer close(p.Out)

	for payload := range p.In {
		fmt.Printf("started sizer %s\n", payload.URL)
		p.Out <- JobResult{
			Payload: payload,
			Length:  len(payload.Content),
		}

		fmt.Printf("completed sizer %s\n", payload.URL)
	}
}
