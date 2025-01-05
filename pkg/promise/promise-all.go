package promise

import (
	"sync"
)

// // Using Buffered error chan
// func PromiseAll(functions ...func() (any, error)) ([]any, error) {
// 	noOfFunctions := len(functions)
// 	var wg sync.WaitGroup
// 	errChan := make(chan error, 1)
// 	responses := make([]any, noOfFunctions)

// 	for index, function := range functions {
// 		wg.Add(1)

// 		go func(index int, function func() (any, error)) {
// 			defer wg.Done()
// 			select {
// 			// Before we do any work, see if someone else has already closed errCh.
// 			// If it's closed or there's already an error in it, exit immediately.
// 			case <-errChan:
// 				// errCh was either closed or already has an error, so stop now.
// 				return
// 			default:
// 				res, err := function()
// 				if err != nil {
// 					// Attempt to store the error if nobody else has done so yet.
// 					select {
// 					case errChan <- err:
// 						// Close 'done' to signal early exit to the other goroutines.
// 						close(errChan)
// 					default:
// 						// errCh already has an error or is closed. Nothing to do.
// 					}
// 				} else {
// 					select {
// 					case <-errChan:
// 						return
// 					default:
// 						responses[index] = res
// 					}
// 				}
// 			}
// 		}(index, function)
// 	}

// 	// We'll wait for all goroutines in a separate goroutine,
// 	// then signal on another channel that we're done.
// 	doneWaiting := make(chan bool)
// 	go func() {
// 		wg.Wait()
// 		close(doneWaiting)
// 	}()

// 	// Wait for either:
// 	// - An error from errChan
// 	// - All goroutines to complete successfully
// 	select {
// 	case err := <-errChan:
// 		// If we get an actual non-nil error, return it.
// 		// Note: Once errCh is closed, reading from it again
// 		// returns zero-value (nil) + ok=false. So the first real
// 		// error is the only one we care about.
// 		return nil, err
// 	case <-doneWaiting:
// 		// No error was sent, success!
// 		return responses, nil
// 	}
// }

// Using Unbuffered error chan

func PromiseAll(functions ...func() (any, error)) ([]any, error) {
	noOfFunctions := len(functions)
	var wg sync.WaitGroup
	errChan := make(chan error)
	responses := make([]any, noOfFunctions)

	for index, function := range functions {
		wg.Add(1)

		go func(index int, function func() (any, error)) {
			defer wg.Done()
			select {
			// Before we do any work, see if someone else has already closed errCh.
			// If it's closed or there's already an error in it, exit immediately.
			case <-errChan:
				// errCh was either closed or already has an error, so stop now.
				return
			default:
				res, err := function()
				if err != nil {
					// Send the error (will block until the main goroutine is reading from errCh).
					errChan <- err
					// After sending the error, close the channel to signal all other tasks.
					close(errChan)
				} else {
					responses[index] = res
				}
			}
		}(index, function)
	}

	// We'll wait for all goroutines in a separate goroutine,
	// then signal on another channel that we're done.
	doneWaiting := make(chan bool)
	go func() {
		wg.Wait()
		close(doneWaiting)
	}()

	// Wait for either:
	// - An error from errChan
	// - All goroutines to complete successfully
	select {
	case err := <-errChan:
		// If we get an actual non-nil error, return it.
		// Note: Once errCh is closed, reading from it again
		// returns zero-value (nil) + ok=false. So the first real
		// error is the only one we care about.
		return nil, err
	case <-doneWaiting:
		// No error was sent, success!
		return responses, nil
	}
}
