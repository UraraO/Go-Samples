package concurrency

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var ID int64 = 1
var IDMut sync.Mutex

func GetID() (id int64) {
	IDMut.Lock()
	defer IDMut.Unlock()
	id = ID
	ID++
	return id
}

func GetIDAtomic() (id int64) {
	id = atomic.AddInt64(&ID, 1)
	return id
}

func GetIDTest() {
	for i := 0; i < 1000; i++ {
		go fmt.Printf("%d\n", GetID())
	}
	time.Sleep(time.Second)
}
