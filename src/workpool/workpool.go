/*
简单的工作池实现，可提交任务，后台进程取任务并执行
*/
package workpool

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type WorkPool struct {
	cap   int64
	size  atomic.Int64
	mut   *sync.Mutex
	cond  *sync.Cond
	cwork chan func()
	exit  chan int
}

func InitWorkPool(cap int64) *WorkPool {
	if cap <= 0 {
		return nil
	}
	mut := &sync.Mutex{}
	p := &WorkPool{
		cap:   cap,
		size:  atomic.Int64{},
		mut:   mut,
		cond:  sync.NewCond(mut),
		cwork: make(chan func()),
		exit:  make(chan int),
	}
	p.size.Store(0)
	return p
}

func work() {
	fmt.Println("work start")
	time.Sleep(time.Millisecond * 500)
	fmt.Println("work done")
}

func (p *WorkPool) RunLoop() {
	for {
		select {
		case work := <-p.cwork:
			fmt.Println("run a work")
			go func(work func()) {
				p.size.Add(1)
				work()
				p.size.Add(-1)
				p.cond.Signal()
			}(work)
		case <-p.exit:
			fmt.Println("receive EXIT cmd")
			for p.size.Load() != 0 {
				time.Sleep(time.Second)
			}
			<-p.exit
			fmt.Println("---EXIT---")
			return
		}
	}
}

func (p *WorkPool) Exit() {
	p.exit <- 1
	p.exit <- 1
}

func (p *WorkPool) Submit(f func()) {
	p.mut.Lock()
	for p.size.Load() == p.cap {
		p.cond.Wait()
	}
	p.cwork <- f
	fmt.Println("submit")
	p.mut.Unlock()
}

func WorkPoolTest() {
	wp := InitWorkPool(10)
	go wp.RunLoop()
	for i := 0; i < 4; i++ {
		go func(p *WorkPool) {
			for i := 0; i < 10; i++ {
				p.Submit(work)
			}
		}(wp)
	}
	time.Sleep(1 * time.Second)
	wp.Exit()
}

// package main

// import "workpool/src/workpool"

// func main() {
// 	workpool.WorkPoolTest()
// }
