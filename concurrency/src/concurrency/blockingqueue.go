package concurrency

import (
	"fmt"
	"sync"
	"time"
)

type BlockQueue struct {
	fullCond  *sync.Cond
	emptyCond *sync.Cond
	elemsMut  *sync.Mutex
	Cap       int64
	Size      int64
	beg       int64
	end       int64
	elems     []interface{}
}

func InitBlockQueue(cap int64) *BlockQueue {
	if cap <= 0 {
		return nil
	}
	emut := &sync.Mutex{}
	return &BlockQueue{
		fullCond:  sync.NewCond(emut),
		emptyCond: sync.NewCond(emut),
		elemsMut:  emut,
		Cap:       cap,
		Size:      0,
		beg:       0,
		end:       0,
		elems:     make([]interface{}, cap+cap/2),
	}
}

// end进，beg出
func (q *BlockQueue) Put(item interface{}) {
	q.elemsMut.Lock()
	defer q.elemsMut.Unlock()
	for q.Size >= q.Cap {
		q.fullCond.Wait()
	}
	q.elems[q.end] = item
	q.end++
	q.Size++
	if q.end == q.Cap {
		q.end = 0
	}
	q.emptyCond.Signal()
}

func (q *BlockQueue) Get() (item interface{}, err error) {
	q.elemsMut.Lock()
	defer q.elemsMut.Unlock()
	for q.Size <= 0 {
		q.emptyCond.Wait()
	}
	item = q.elems[q.beg]
	q.beg++
	q.Size--
	if q.beg == q.Cap {
		q.beg = 0
	}
	q.fullCond.Signal()
	return item, nil
}

var cap = 30
var testTimes = 3

func PutTest(q *BlockQueue) {
	for i := 0; i < cap; i++ {
		fmt.Printf("put %v\n", i)
		q.Put(i)
	}
}

func GetTest(q *BlockQueue) {
	for i := 0; i < cap; i++ {
		v, _ := q.Get()
		fmt.Printf("get %v\n", v)
	}
}

func BlockQueueTest() {
	q := InitBlockQueue(10)
	for i := 0; i < testTimes; i++ {
		go PutTest(q)
		go GetTest(q)
	}
	time.Sleep(time.Second)
}
