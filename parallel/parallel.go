package parallel

import (
	"sync"
)

// Runs a function over an array of inputs with a diffined set of parrallel go routines, and returns an array with the results
// max_parallel: the maximun number of jobs that will run in parallel
// input: array of input values
// job: a function that will be run wil the input values
func Run[v any, w any](max_parallel uint, input []v, job func(v) (w, error)) ([]w, []error) {
	to_val_sum := make(chan w)
	to_err_sum := make(chan error)

	to_val_return := make(chan []w)
	to_err_return := make(chan []error)

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
	count := uint(0)
	for _, i := range input {
		wg.Add(1)
		count++

		//Run job with the input function
		go func(j v) {
			defer wg.Done()

			value, err := job(j)

			if err != nil {
				to_err_sum <- err
				return
			}

			to_val_sum <- value
		}(i)

		//If it exceed the max parallel value, then wait and reset count
		if count >= max_parallel {
			wg.Wait()
			count = 0
		}

	}

	//Wait if there are jobs still running
	if count < max_parallel {
		wg.Wait()
	}

	//Close channels to stop summary go routines
	close(to_val_sum)
	close(to_err_sum)

	return <-to_val_return, <-to_err_return
}
