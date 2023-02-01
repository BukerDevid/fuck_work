package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(chan interface{})

	for _, stage := range stages {
		in = stage(in)
	}

	go func(in In, out chan interface{}) {
		defer close(out)

		for {
			select {
			case v, ok := <-in:
				if ok {
					out <- v
					continue
				}
				return
			case <-done:
				return
			}
		}
	}(in, out)

	return out
}
