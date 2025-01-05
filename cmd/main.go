package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/AdarshMishrji/go-async/pkg/promise"
)

func promise1() (any, error) {
	time.Sleep(1 * time.Second)
	println("[INTERNAL] resolved promise1")
	return 1, nil
}
func promise2() (any, error) {
	time.Sleep(2 * time.Second)
	println("[INTERNAL] resolved promise2")
	return nil, errors.New("promise 2 error")
}
func promise3() (any, error) {
	time.Sleep(3 * time.Second)
	println("[INTERNAL] resolved promise3")
	return 3, nil
}

func main() {
	println("*************** ALL ********************")
	resAll, err := promise.PromiseAll(promise1, promise2, promise3)
	if err != nil {
		println(err.Error())
	} else {
		println(len(resAll))
		for _, val := range resAll {
			fmt.Printf("Struct: %+v\n", val) // Prints field names and values
		}
	}

	println("*************** ALL Settled ********************")
	resAllSettled, err := promise.PromiseAllSettled(promise1, promise2, promise3)
	if err != nil {
		println(err.Error())
	} else {
		println(len(resAllSettled))
		for _, val := range resAllSettled {
			fmt.Printf("Struct: %+v\n", val) // Prints field names and values
		}
	}

	println("*************** Race ********************")
	resRace, err := promise.PromiseRace(promise1, promise2, promise3)
	if err != nil {
		println(err.Error())
	} else {
		fmt.Printf("Struct: %+v\n", resRace) // Prints field names and values
	}

	time.Sleep(5 * time.Second) // waiting to see all the internal logs of promises
	println("exited")
}
