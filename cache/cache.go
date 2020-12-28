package cache

import (
	"sync"
	"time"
)

// Cache ...
type Cache struct {
	sync.RWMutex
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
	Items             map[string]Item
}

// Item ...
type Item struct {
	Value      interface{}
	Created    time.Time
	Expiration int64
}

// New ...
func New(defaultExpiration, cleanupInterval time.Duration) *Cache {

	// инициализируем карту(map) в паре ключ(string)/значение(Item)
	items := make(map[string]Item)

	cache := Cache{
		Items:             items,
		DefaultExpiration: defaultExpiration,
		CleanupInterval:   cleanupInterval,
	}

	// Если интервал очистки больше 0, запускаем GC (удаление устаревших элементов)
	if cleanupInterval > 0 {
		cache.StartGC()
	}

	return &cache
}

// StartGC ...
func (c *Cache) StartGC() {
	go c.GC()
}

// GC ...
func (c *Cache) GC() {

	for {
		// ожидаем время установленное в cleanupInterval
		<-time.After(c.CleanupInterval)

		if c.Items == nil {
			return
		}

		// Ищем элементы с истекшим временем жизни и удаляем из хранилища
		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)

		}

	}

}

// expiredKeys возвращает список "просроченных" ключей
func (c *Cache) expiredKeys() (keys []string) {

	c.RLock()

	defer c.RUnlock()

	for k, i := range c.Items {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return
}

// clearItems удаляет ключи из переданного списка, в нашем случае "просроченные"
func (c *Cache) clearItems(keys []string) {

	c.Lock()

	defer c.Unlock()

	for _, k := range keys {
		delete(c.Items, k)
	}
}
