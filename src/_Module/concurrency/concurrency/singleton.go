package concurrency

import (
	"fmt"
	"sync"
	"time"
)

type Instance struct {
	message string
}

var pins *Instance = nil
var mut *sync.Mutex = &sync.Mutex{}

func GetInstance() *Instance {
	if pins != nil {
		return pins
	}
	mut.Lock()
	defer mut.Unlock()
	if pins != nil {
		return pins
	}
	pins = &Instance{
		message: "hello",
	}
	return pins
}

func getIns() {
	ins := GetInstance()
	fmt.Printf("%v\n", ins.message)
}

func SingletonTest() {
	for i := 0; i < 100; i++ {
		go getIns()
	}
	time.Sleep(1 * time.Second)
}
