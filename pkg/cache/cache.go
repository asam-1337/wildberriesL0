package cache

import (
	"fmt"
	"github.com/asam-1337/wildberriesL0/internal/domain/entity"
	"sync"
)

type Cache struct {
	mu     sync.RWMutex
	orders map[string]entity.Order
}

func NewCache(mu *sync.RWMutex) *Cache {
	return &Cache{
		orders: make(map[string]entity.Order),
	}
}

func (c *Cache) Load(key string) (entity.Order, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, ok := c.orders[key]
	if !ok {
		return val, fmt.Errorf("not found")
	}

	return val, nil
}

func (c *Cache) Store(key string, value entity.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.orders[key] = value
}

