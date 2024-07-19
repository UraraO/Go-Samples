package main

import (
	"fmt"
	"heartcheck/src/heartcheck"
	"time"
)

func main() {
	handler := func() {
		fmt.Println("handler handle")
	}
	dur := 5 * time.Second
	hc := heartcheck.InitHeartChecker(dur, handler)
	go hc.Check()
	for {
		select {
		case <-hc.Offline: // 超时
			fmt.Println("Client: server offline")
			time.Sleep(1 * time.Second)
			return
		case time := <-hc.Server:
			fmt.Println("times:", time)
		}
	}
}
