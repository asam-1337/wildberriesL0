package cache

import (
	"github.com/asam-1337/wildberriesL0/internal/domain/entity"
	"github.com/asam-1337/wildberriesL0/internal/localErrors"
	"sync"
	"time"
)

type Item struct {
	Value      entity.Order
	Expiration int64
}

type Cache struct {
	items             map[string]Item
	mu                sync.RWMutex
	cleanupInterval   time.Duration
	defaultExpiration time.Duration
}

func NewCache(cleanupInterval, defaultExpiration time.Duration) *Cache {
	cache := &Cache{
		items:             make(map[string]Item),
		mu:                sync.RWMutex{},
		cleanupInterval:   cleanupInterval,
		defaultExpiration: defaultExpiration,
	}

	if cleanupInterval > 0 {
		go cache.garbageCollector()
	}

	return cache
}

func (c *Cache) garbageCollector() {
	for {
		<-time.After(c.cleanupInterval)

		c.mu.RLock()
		for key, item := range c.items {
			if item.Expiration < time.Now().UnixNano() {
				c.Delete(key)
			}
		}
		c.mu.RUnlock()
	}
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

func (c *Cache) Load(key string) (entity.Order, error) {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()

	if !ok {
		return entity.Order{}, localErrors.ErrCashNotFound
	}

	return item.Value, nil
}

func (c *Cache) Store(value entity.Order) error {
	exp := time.Now().Add(c.defaultExpiration).UnixNano()
	key := value.OrderUID
	c.mu.Lock()

	c.items[key] = Item{
		Value:      value,
		Expiration: exp,
	}

	c.mu.Unlock()
	return nil
}

func (c *Cache) Exist(key string) bool {
	c.mu.RLock()
	_, ok := c.items[key]
	c.mu.RUnlock()
	return ok
}
