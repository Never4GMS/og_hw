package hw06pipelineexecution

import (
	"runtime"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil {
		temp := make(chan interface{})
		in = temp
		close(temp)
	}

	if len(stages) == 0 {
		return in
	}

	runners := max(runtime.NumCPU(), 5)
	inputs := make([]Bi, runners)
	outputs := make([]In, runners)
	pipeCompleted := make(chan In, runners)

	for i := 0; i < runners; i++ {
		inputs[i] = make(Bi)
		outputs[i] = createRunner(inputs[i], stages)
	}

	go scheduleRunners(done, in, inputs, outputs, pipeCompleted, runners)

	return merge(done, pipeCompleted)
}

func createRunner(in In, stages []Stage) In {
	for _, stage := range stages {
		in = stage(in)
	}
	return in
}

func scheduleRunners(done In, in In, inputs []Bi, outputs []In, pipeCompleted chan<- In, runners int) {
	defer closeAll(inputs)
	defer close(pipeCompleted)
	roundRobbin := 0
	for {
		select {
		case <-done:
			return
		case i, ok := <-in:
			if !ok {
				return
			}
			inputs[roundRobbin] <- i
			pipeCompleted <- outputs[roundRobbin]
			roundRobbin = (roundRobbin + 1) % runners
		}
	}
}

func closeAll(inputs []Bi) {
	for _, input := range inputs {
		close(input)
	}
}

func merge(done In, pipeCompleted <-chan In) In {
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case pipe, ok := <-pipeCompleted:
				if !ok {
					return
				}

				select {
				case <-done:
					return
				case v, ok := <-pipe:
					if !ok {
						return
					}
					out <- v
				}
			}
		}
	}()

	return out
}
