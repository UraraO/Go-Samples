package concurrency

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"
)

type Spinlock struct {
	Spinning atomic.Bool
	Waiters  atomic.Int64
}

func InitSpinlock() *Spinlock {
	sl := &Spinlock{}
	sl.Spinning.Store(false)
	sl.Waiters.Store(0)
	return sl
}

func (sl *Spinlock) Lock() {
	sl.Waiters.Add(1)
	for {
		if sl.Spinning.CompareAndSwap(false, true) {
			sl.Waiters.Add(-1)
			break
		}
		runtime.Gosched()
	}
}

func (sl *Spinlock) Unlock() {
	sl.Spinning.CompareAndSwap(true, false)
}

func testsl(sl *Spinlock) {
	sl.Lock()
	fmt.Printf("sl.waiters = %v\n", sl.Waiters.Load())
	sl.Unlock()
}

func SlTest() {
	sl := InitSpinlock()
	for i := 0; i < 100; i++ {
		go testsl(sl)
	}
	time.Sleep(time.Second)

}
