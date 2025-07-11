package cache

import (
	"github.com/maypok86/otter/v2"
	"wb-l0/internal/types/dto"
)

type OrderCache struct {
	store *otter.Cache[string, dto.Order]
}

func New() *OrderCache {
	opts := &otter.Options[string, dto.Order]{
		MaximumSize: 1000,

		//	InitialCapacity:  0,
		//	Logger:           nil,
	}

	cache := otter.Must[string, dto.Order](opts)

	return &OrderCache{store: cache}
}

func (c *OrderCache) Set(key string, order dto.Order) {
	c.store.Set(key, order)
}

func (c *OrderCache) Get(key string) (dto.Order, bool) {
	order, ok := c.store.GetIfPresent(key)

	return order, ok
}
