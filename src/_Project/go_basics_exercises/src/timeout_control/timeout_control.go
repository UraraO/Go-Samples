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
