package pipeline

const (
	jobBuffer = 100
)

type process interface {
	Run()
}

type Pipeline struct {
	processes []process
}

func NewPipeline() *Pipeline {
	return &Pipeline{}
}

func (p *Pipeline) AddProcess(proc process) {
	p.processes = append(p.processes, proc)
}

func (p *Pipeline) AddProcesses(procs ...process) {
	for _, proc := range procs {
		p.AddProcess(proc)
	}
}

func (p *Pipeline) Run() {
	for i, proc := range p.processes {
		if i < len(p.processes)-1 {
			go proc.Run()
		} else {
			proc.Run()
		}
	}
}
