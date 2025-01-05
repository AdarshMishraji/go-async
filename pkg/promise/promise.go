package promise

import "fmt"

type Promise struct {
	executer    func(resolve func(any), reject func(error))
	result      any
	err         error
	doneChannel chan bool
}

type PromiseInterface interface {
	Then(onFulfilled func(any) any) *Promise
	Catch(onRejected func(error) error) *Promise
	Finally(onFinally func() (any, error))
}

func (promise *Promise) resolve(result any) {
	select {
	case <-promise.doneChannel:
		return
	default:
	}

	switch response := result.(type) {
	case *Promise: // if result is promise then chained Then will get the response or error in Catch
		res, err := response.Await()
		if err != nil {
			promise.reject(err)
			return
		}

		promise.result = res
	default:
		promise.result = response
	}
	close(promise.doneChannel) // await will now unblocked
}

func (promise *Promise) reject(err error) {
	select {
	case <-promise.doneChannel:
		return
	default:
	}

	promise.err = err
	close(promise.doneChannel) // await will now unblocked
}

func (promise *Promise) handlePanic() {
	r := recover()
	if r != nil {
		switch err := r.(type) {
		case error:
			promise.reject(fmt.Errorf("panic recovery with error: %s", err.Error()))
		default:
			promise.reject(fmt.Errorf("panic recovery with unknown error: %s", fmt.Sprint(err)))
		}
	}
}

func New(executer func(resolve func(any), reject func(error))) *Promise {
	promise := &Promise{
		executer:    executer,
		result:      nil,
		err:         nil,
		doneChannel: make(chan bool), // as it is an unbuffered channel, Await function will wait the promise to either resolve or reject
	}

	go func() {
		defer promise.handlePanic()
		promise.executer(promise.resolve, promise.reject)
	}()

	return promise
}

// if Await is not called, waiting of the promise will not happen
func (promise *Promise) Await() (any, error) {
	<-promise.doneChannel // waiting for promise to either resolve or reject
	return promise.result, promise.err
}

func (promise *Promise) Then(onFulfilled func(data any) any) *Promise {
	return New(func(resolve func(any), reject func(error)) {
		res, err := promise.Await()
		if err != nil {
			reject(err)
		} else {
			resolve(onFulfilled(res))
		}
	})
}

func (promise *Promise) Catch(onRejected func(error) error) *Promise {
	return New(func(resolve func(any), reject func(error)) {
		res, err := promise.Await()
		if err != nil {
			reject(onRejected(err))
		} else {
			resolve(res)
		}
	})
}

func (promise *Promise) Finally(onFinally func()) (any, error) {
	promise.Await()
	onFinally()
	return promise.result, promise.err
}
