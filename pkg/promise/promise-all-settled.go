package promise

import "sync"

func PromiseAllSettled(promises ...*Promise) []Result {
	length := len(promises)
	responses := make([]Result, length)

	var wg sync.WaitGroup

	for index, promise := range promises {
		wg.Add(1)
		go func(index int, promise *Promise) {
			defer wg.Done()
			res, err := promise.Await()
			responses[index] = Result{res, err}
		}(index, promise)
	}

	wg.Wait()

	return responses
}
