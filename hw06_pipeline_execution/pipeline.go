package hw06pipelineexecution

import "sync"

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
		input = merge(done,
			stage(input),
			stage(input),
			stage(input),
		)
	}

	return input
}

func merge(done In, outputs ...In) In {
	var wg sync.WaitGroup
	out := make(Bi)
	multiplex := func(c In) {
		defer wg.Done()
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
	}

	wg.Add(len(outputs))
	for _, c := range outputs {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
