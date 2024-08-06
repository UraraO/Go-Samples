/*
* @Author: chaidaxuan chaidaxuan@wps.cn
* @Date: 2024-07-26 15:45:10
* @LastEditors: chaidaxuan chaidaxuan@wps.cn
* @LastEditTime: 2024-07-26 17:46:12
* @FilePath: /urarao/GoProjects/Golang-Samples/src/stream_control/timeout_control.go
* @Description:

流控制示例
分别使用context withtimeout 和 time.Timer实现

* Copyright (c) 2024 by ${git_name_email}, All Rights Reserved.
*/
package timeoutcontrol

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func SlowQuery() error {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	n := r.Intn(1000)
	time.Sleep(time.Duration(n) * time.Millisecond)
	return nil
}

func TimerTest() {
	tm := time.NewTimer(500 * time.Millisecond)
	begin := time.Now()
	fc := make(chan int)
	go wrappedFuncCall(SlowQuery, fc)
	select {
	case <-tm.C:
		fmt.Println("timeout occur: timer")
	case <-fc:
		fmt.Println("SlowQuery return, use time:", (time.Since(begin)))
	}
}

func wrappedFuncCall(f func() error, c chan int) {
	f()
	c <- 1
}

func ContextTest() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()
	begin := time.Now()
	fc := make(chan int)
	go wrappedFuncCall(SlowQuery, fc)
	select {
	case <-ctx.Done():
		fmt.Println("timeout occur: timer")
	case <-fc:
		fmt.Println("SlowQuery return, use time:", (time.Since(begin)))
	}
}

// package main

// import (
// 	timeoutcontrol "stream_control/src/timeout_control"
// 	"time"
// )

// func main() {
// 	timeoutcontrol.TimerTest()
// 	time.Sleep(1 * time.Second)
// 	timeoutcontrol.ContextTest()
// }
