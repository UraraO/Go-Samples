package concurrency

import "sync"

type KVCache struct {
	kvs sync.Map
}

func (c *KVCache) Put(k string, v any) {
	c.kvs.Store(k, v)
}

func (c *KVCache) Get(k string) (v any, ok bool) {
	return c.kvs.Load(k)
}

func (c *KVCache) Delete(k string) (v any, ok bool) {
	return c.kvs.LoadAndDelete(k)
}
