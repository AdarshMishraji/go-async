package promise

func PromiseRace(promises ...*Promise) (any, error) {
	responseChan := make(chan *Result)

	for _, promise := range promises {
		go func(promise *Promise) {
			select {
			case <-responseChan:
				return
			default:
				res, err := promise.Await()
				select {
				case <-responseChan:
					return
				default:
					responseChan <- &Result{res, err}
					close(responseChan)
				}
			}
		}(promise)
	}

	result := <-responseChan
	return result.Value, result.Err
}
