package promise

import "sync"

func PromiseAll(promises ...*Promise) ([]any, error) {
	length := len(promises)
	var wg sync.WaitGroup

	responses := make([]any, length)
	errorChan := make(chan error)

	for index, promise := range promises {
		wg.Add(1)
		go func(index int, promise *Promise) {
			defer wg.Done()
			select {
			case <-errorChan:
				return
			default:
				res, err := promise.Await()
				select {
				case <-errorChan:
					return
				default:
					if err != nil {
						errorChan <- err
						close(errorChan)
					} else {
						responses[index] = res
					}
				}
			}
		}(index, promise)
	}

	doneWaitingChan := make(chan bool)
	go func() {
		wg.Wait()
		close(doneWaitingChan)
	}()

	select {
	case err := <-errorChan:
		return nil, err
	case <-doneWaitingChan:
		return responses, nil
	}
}
