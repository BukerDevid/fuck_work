package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mx sync.RWMutex

	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (c *cacheItem) GetValue() interface{} { return c.value }

func (c *cacheItem) GetKey() Key { return c.key }

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	lc.mx.Lock()
	defer lc.mx.Unlock()

	if item, ok := lc.items[key]; ok {

		item.Value = cacheItem{
			key:   key,
			value: value,
		}

		lc.queue.MoveToFront(item)
		return true
	}

	if lc.queue.Len() == lc.capacity {
		rmItem := lc.queue.Back()

		//INFO: небезопасная дрянь
		key := rmItem.Value.(cacheItem)

		delete(lc.items, key.GetKey())
		lc.queue.Remove(rmItem)
	}

	lc.items[key] = lc.queue.PushFront(cacheItem{
		key:   key,
		value: value,
	})

	return false
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	lc.mx.Lock()
	defer lc.mx.Unlock()

	if item, ok := lc.items[key]; ok {
		cacheValue := item.Value.(cacheItem)
		return cacheValue.GetValue(), ok
	}

	return nil, false
}

func (lc *lruCache) Clear() {
	lc.mx.Lock()
	defer lc.mx.Unlock()

	lc.items = make(map[Key]*ListItem, lc.capacity)
	lc.queue = NewList()
}
