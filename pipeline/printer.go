package pipeline

import "fmt"

type printer struct {
	process
	In chan JobResult
}

func NewPrinter() *printer {
	return &printer{
		In: make(chan JobResult, jobBuffer),
	}
}

func (p *printer) Run() {
	for s := range p.In {
		fmt.Printf("%s downloaded %d bytes in %s\n", s.URL, s.Length, s.Time)
	}
}
