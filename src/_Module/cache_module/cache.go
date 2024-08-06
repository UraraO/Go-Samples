/*===========
 Author: UraraO Haru_UraraO@outlook.com
 Date: 2024-08-06 19:47:01
 LastEditors: UraraO Haru_UraraO@outlook.com
 LastEditTime: 2024-08-06 22:35:16
 FilePath: /Golang-Samples/src/_Module/cache_module/cache.go
 Description:

 简单实现的内存k-v缓存模块

 Copyright (c) 2024 by ${git_name_email}, All Rights Reserved.
===========*/

package cache

import (
	"sync"
	"time"
)

type elem struct {
	Val        interface{}
	ExpireTime time.Time
}

type Cache struct {
	elems     map[string]elem
	dataMutex sync.Mutex
	Capacity  uint
}

func InitCache(cap uint) Cache {
	return Cache{
		elems:     make(map[string]elem),
		dataMutex: sync.Mutex{},
		Capacity:  cap,
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.dataMutex.Lock()
	defer c.dataMutex.Unlock()
	val, ok := c.elems[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(val.ExpireTime) {
		delete(c.elems, key)
		return nil, false
	}
	return val, true
}

func (c *Cache) Put(key string, val interface{}, dur time.Duration) bool {
	c.dataMutex.Lock()
	defer c.dataMutex.Unlock()
	if len(c.elems) == int(c.Capacity) {
		for k, v := range c.elems {
			if time.Now().After(v.ExpireTime) {
				delete(c.elems, k)
			}
		}
	}
	_, ok := c.elems[key]
	if len(c.elems) == int(c.Capacity) && !ok { // 缓存仍满，且元素不在缓存中，插入失败
		return false
	}
	elem := elem{
		Val:        val,
		ExpireTime: time.Now().Add(dur),
	}
	c.elems[key] = elem
	return true
}

func (c *Cache) Delete(key string) {
	c.dataMutex.Lock()
	defer c.dataMutex.Unlock()
	delete(c.elems, key)
}
