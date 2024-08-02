/*
 * @Author: chaidaxuan chaidaxuan@wps.cn
 * @Date: 2024-07-29 14:15:25
 * @LastEditors: chaidaxuan chaidaxuan@wps.cn
 * @LastEditTime: 2024-07-29 16:33:55
 * @FilePath: /urarao/GoProjects/chaidaxuan/filesync/src/utils/blockqueue.go
 * @Description:
 *
 * Copyright (c) 2024 by ${git_name_email}, All Rights Reserved.
 */
package utils

import (
	"filesync/src/def"
	"fmt"
	"sync"
	"time"
)

type Mission struct {
	Type     int
	FileName string
}

func NewMission(mType int, fileName string) *Mission {
	return &Mission{
		Type:     mType,
		FileName: fileName,
	}
}

func (m *Mission) Content() string {
	switch m.Type {
	case def.MISSION_TYPE_INIT:
		return "MISSION_TYPE_INIT"
	case def.MISSION_TYPE_READ:
		return "MISSION_TYPE_READ"
	case def.MISSION_TYPE_WRITE:
		return "MISSION_TYPE_WRITE"
	case def.MISSION_TYPE_UPLOAD:
		return "MISSION_TYPE_UPLOAD"
	case def.MISSION_TYPE_DOWNLOAD:
		return "MISSION_TYPE_DOWNLOAD"
	case def.MISSION_TYPE_LIST:
		return "MISSION_TYPE_LIST"
	case def.MISSION_TYPE_DELETE:
		return "MISSION_TYPE_DELETE"
	default:
		return "UNKNOWN MISSION TYPE"
	}
}

type BlockQueue struct {
	fullCond  *sync.Cond  // 队列满时阻塞
	emptyCond *sync.Cond  // 队列空时阻塞
	elemsMut  *sync.Mutex // 保护共享数据，该锁不仅保护elems数组，也用于两个条件变量的触发
	Cap       int64       // 容量
	Size      int64       // 元素数量
	beg       int64       // 队列起始位置
	end       int64       // 队列末尾位置
	elems     []*Mission
}

func NewBlockQueue(cap int64) *BlockQueue {
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
		elems:     make([]*Mission, cap),
	}
}

// end进，beg出
func (q *BlockQueue) Put(item *Mission) {
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

func (q *BlockQueue) Get() (item *Mission) {
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
	return item
}

var cap = 30      // 测试容量
var testTimes = 3 // 测试循环次数

func PutTest(q *BlockQueue) {
	for i := 0; i < cap; i++ {
		fmt.Printf("put %v\n", i)
		q.Put(NewMission(def.MISSION_TYPE_INIT, "test"))
	}
}

func GetTest(q *BlockQueue) {
	for i := 0; i < cap; i++ {
		v := q.Get()
		fmt.Printf("get %v\n", v)
	}
}

func BlockQueueTest() {
	q := NewBlockQueue(10)
	for i := 0; i < testTimes; i++ {
		go PutTest(q)
		go GetTest(q)
	}
	time.Sleep(time.Second)
}
