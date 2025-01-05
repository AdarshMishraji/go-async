package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/AdarshMishrji/go-async/pkg/promise"
)

var (
	promise1 = promise.New(func(resolve func(any), reject func(error)) {
		time.Sleep(1 * time.Second)
		println("[INTERNAL] resolved promise1")
		resolve(1)
	})
	promise2 = promise.New(func(resolve func(any), reject func(error)) {
		time.Sleep(2 * time.Second)
		println("[INTERNAL] resolved promise2")
		reject(errors.New("some error"))
	})
	promise3 = promise.New(func(resolve func(any), reject func(error)) {
		time.Sleep(3 * time.Second)
		println("[INTERNAL] resolved promise3")
		resolve(3)
	})
)

func main() {
	res, err := promise.New(func(resolve func(any), reject func(error)) {
		resolve(promise1)
	}).Then(func(data any) any {
		fmt.Printf("Data: %+v\n", data)
		return promise2
	}).Finally(func() {
		println("Finally")
	})
	fmt.Printf("Res: %+v\n", res)
	fmt.Printf("Err: %+v\n", err)

	println("*************** ALL ********************")
	resAll, err := promise.PromiseAll(promise1, promise2, promise3)
	if err != nil {
		println(err.Error())
	} else {
		for _, val := range resAll {
			fmt.Printf("Struct: %+v\n", val) // Prints field names and values
		}
	}

	println("*************** ALL Settled ********************")
	resAllSettled := promise.PromiseAllSettled(promise1, promise2, promise3)
	for _, val := range resAllSettled {
		fmt.Printf("Struct: %+v\n", val) // Prints field names and values
	}

	println("*************** Race ********************")
	resRace, err := promise.PromiseRace(promise1, promise2, promise3)
	if err != nil {
		println(err.Error())
	} else {
		fmt.Printf("Struct: %+v\n", resRace) // Prints field names and values
	}

	time.Sleep(5 * time.Second) // waiting to see all the internal logs of promises
	// println("exited")
}
