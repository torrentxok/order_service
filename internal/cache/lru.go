package cache

import (
	"container/list"
	"sync"

	"github.com/torrentxok/order_service/internal/models"
)

type entry struct {
	key   string
	value *models.Order
}

type LRUCache struct {
	capacity int
	items    map[string]*list.Element
	list     *list.List
	mu       sync.Mutex
}

func NewLRUCache(capacity int) *LRUCache {
	if capacity <= 0 {
		panic("cache capacity must be positive")
	}

	return &LRUCache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

func (c *LRUCache) Get(key string) (*models.Order, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.list.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}

	return nil, false
}

func (c *LRUCache) Set(key string, value *models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		elem.Value.(*entry).value = value
		c.list.MoveToFront(elem)
		return
	}

	elem := c.list.PushFront(&entry{key: key, value: value})
	c.items[key] = elem

	if c.list.Len() > c.capacity {
		c.evict()
	}
}

func (c *LRUCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.list.Remove(elem)
		ent := elem.Value.(*entry)
		delete(c.items, ent.key)
	}
}

func (c *LRUCache) Capacity() int {
	return c.capacity
}

func (c *LRUCache) evict() {
	elem := c.list.Back()
	if elem == nil {
		return
	}

	c.list.Remove(elem)
	ent := elem.Value.(*entry)
	delete(c.items, ent.key)
}
