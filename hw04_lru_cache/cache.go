package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type queueItem struct {
	value interface{}
	key   Key
}

func (c *lruCache) set(key Key, value interface{}) bool {
	exists := false

	if item, ok := c.items[key]; ok {
		if qItem, ok := item.Value.(queueItem); ok {
			qItem.value = value
			item.Value = qItem
		}
		c.queue.MoveToFront(item)
		exists = true
	} else {
		c.items[key] = c.queue.PushFront(queueItem{
			value: value,
			key:   key,
		})
	}

	if c.queue.Len() > c.capacity {
		item := c.queue.Back()
		c.queue.Remove(item)
		if qItem, ok := item.Value.(queueItem); ok {
			delete(c.items, qItem.key)
		}
	}

	return exists
}

func (c *lruCache) get(key Key) (interface{}, bool) {
	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		if qi, ok := item.Value.(queueItem); ok {
			return qi.value, true
		}
	}

	return nil, false
}

func (c *lruCache) clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.Lock()
	defer c.Unlock()
	return c.set(key, value)
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.Lock()
	defer c.Unlock()
	return c.get(key)
}

func (c *lruCache) Clear() {
	c.Lock()
	defer c.Unlock()
	c.clear()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
