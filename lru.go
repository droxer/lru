package lru

import (
	"container/list"
	"sync"
)

type Cache struct {
	size      int
	evictList *list.List
	entries   map[interface{}]*list.Element
	lock      sync.Mutex
}

type entry struct {
	key   interface{}
	value interface{}
}

func New(size int) *Cache {
	return &Cache{
		size:      size,
		evictList: list.New(),
		entries:   make(map[interface{}]*list.Element, size),
	}
}

func (c *Cache) Reset() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.evictList = list.New()
	c.entries = make(map[interface{}]*list.Element, c.size)
}

func (c *Cache) Add(key, value interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if item, ok := c.entries[key]; ok {
		c.evictList.MoveToFront(item)
		item.Value.(*entry).value = value
		return
	}

	if c.evictList.Len() == c.size {
		c.removeOldest()
	}

	item := &entry{key, value}
	c.entries[key] = c.evictList.PushFront(item)
}

func (c *Cache) Get(key interface{}) (value interface{}, ok bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if item, ok := c.entries[key]; ok {
		c.evictList.MoveToFront(item)
		return item.Value.(*entry).value, true
	}

	return nil, false
}

func (c *Cache) removeOldest() {
	item := c.evictList.Back()
	if item == nil {
		return
	}

	c.evictList.Remove(item)
	kv := item.Value.(*entry)
	delete(c.entries, kv.key)
}
