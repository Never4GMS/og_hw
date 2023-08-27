package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil || len(stages) == 0 {
		output := make(Bi)
		close(output)
		return output
	}

	input := in
	for _, stage := range stages {
		out := make(Bi)
		go func(c In) {
			defer close(out)
			for {
				select {
				case <-done:
					return
				case val, ok := <-c:
					if !ok {
						return
					}
					select {
					case <-done:
						return
					case out <- val:
					}
				}
			}
		}(input)
		input = stage(out)
	}

	return input
}
