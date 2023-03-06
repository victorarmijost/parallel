package fp

import "sync"

// Runs a function over an slice of inputs with a defined set of parrallel go routines, and returns an slice with the results
// max_parallel: the maximun number of jobs that will run in parallel
// input: slice of input values
// job: a function that will be run wil the input values
func Parallel[v any, w any](max_parallel uint, input []v, job func(v) (w, error)) ([]w, []error) {
	to_val_sum := make(chan w)
	to_err_sum := make(chan error)

	to_val_return := make(chan []w)
	to_err_return := make(chan []error)

	inChan := make(chan v)
	doneChan := make(chan struct{}, max_parallel)

	//Stream inputs
	go func() {
		for _, in := range input {
			inChan <- in
			doneChan <- struct{}{}
		}

		close(inChan)
	}()

	//Summarize values
	go func() {
		agg := []w{}

		for result := range to_val_sum {
			agg = append(agg, result)
		}

		to_val_return <- agg
	}()

	//Summarize errors
	go func() {
		agg := []error{}

		for err := range to_err_sum {
			agg = append(agg, err)
		}

		to_err_return <- agg
	}()

	//Loop inputs an run function
	wg := sync.WaitGroup{}
	for in := range inChan {
		wg.Add(1)
		//Run job with the input function
		go func(j v) {
			defer func() { <-doneChan; wg.Done() }()
			value, err := job(j)

			if err != nil {
				to_err_sum <- err
				return
			}

			to_val_sum <- value
		}(in)
	}

	wg.Wait()

	//Close channels to stop summary go routines
	close(to_val_sum)
	close(to_err_sum)

	return <-to_val_return, <-to_err_return
}
