package promise

import "sync"

type Result struct {
	Value any
	Err   error
}

func PromiseAllSettled(functions ...func() (any, error)) ([]Result, error) {
	noOfFunctions := len(functions)
	var wg sync.WaitGroup
	responses := make([]Result, noOfFunctions)

	for index, function := range functions {
		wg.Add(1)

		go func(index int, function func() (any, error)) {
			defer wg.Done()
			res, err := function()
			responses[index] = Result{res, err}
		}(index, function)
	}

	wg.Wait()

	return responses, nil
}
