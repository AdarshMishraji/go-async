package promise

func PromiseRace(functions ...func() (any, error)) (any, error) {
	resultCh := make(chan Result)

	for _, function := range functions {
		go func(fn func() (any, error)) {
			select {
			case <-resultCh:
				return
			default:
				res, err := fn()
				select {
				case <-resultCh:
					return
				default:
					resultCh <- Result{res, err}
					close(resultCh)
				}
			}
		}(function)
	}

	// Receive the first result or error
	res := <-resultCh
	return res.Value, res.Err
}
