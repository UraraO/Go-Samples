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
	Elems     map[string]elem
	dataMutex sync.Mutex
	Capacity  uint
}

func InitCache(cap uint) Cache {
	return Cache{
		Elems:     make(map[string]elem),
		dataMutex: sync.Mutex{},
		Capacity:  cap,
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.dataMutex.Lock()
	val, ok := c.Elems[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(val.ExpireTime) {
		delete(c.Elems, key)
		return nil, false
	}
	c.dataMutex.Unlock()
	return val, true
}

func (c *Cache) Put(key string, val interface{}, dur time.Duration) bool {
	c.dataMutex.Lock()
	if len(c.Elems) == int(c.Capacity) {
		for k, v := range c.Elems {
			if time.Now().After(v.ExpireTime) {
				delete(c.Elems, k)
			}
		}
	}
	elem, ok := c.Elems[key]
	if len(c.Elems) == int(c.Capacity) && !ok { // 缓存仍满，且元素不在缓存中，插入失败
		return false
	}
	elem.ExpireTime = time.Now().Add(dur)
	elem.Val = val
	c.Elems[key] = elem
	c.dataMutex.Unlock()
	return true
}

func (c *Cache) Delete(key string) {
	c.dataMutex.Lock()
	delete(c.Elems, key)
	c.dataMutex.Unlock()
}
