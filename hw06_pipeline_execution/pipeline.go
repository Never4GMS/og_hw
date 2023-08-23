package hw06pipelineexecution

import "sync"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil {
		b := make(Bi)
		in = b
		close(b)
	}

	if len(stages) == 0 {
		return in
	}

	input := wrapWithDone(done, in)
	for _, stage := range stages {
		input = wrapWithDone(done, merge(done, stage(input), stage(input)))
	}

	return input
}

func wrapWithDone(done In, in In) In {
	bi := make(Bi)
	go func() {
		defer close(bi)
		for {
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case bi <- val:
				}
			}
		}
	}()
	return bi
}

func merge(done In, outputs ...In) In {
	var wg sync.WaitGroup
	out := make(Bi)
	multiplex := func(c In) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case out <- i:
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
