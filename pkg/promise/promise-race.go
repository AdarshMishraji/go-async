package promise

func PromiseRace(functions ...func() (any, error)) (any, error) {
	resultCh := make(chan Result)

	for _, function := range functions {
		go func(fn func() (any, error)) {
			res, err := fn()
			select {
			case resultCh <- Result{res, err}:
			default:
			}
		}(function)
	}

	// Receive the first result or error
	res := <-resultCh
	return res.Value, res.Err
}
