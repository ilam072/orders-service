package cache

import (
	"errors"
	"github.com/maypok86/otter/v2"
	"wb-l0/internal/types/dto"
	"wb-l0/pkg/e"
)

var (
	ErrOrderNotFound = errors.New("order not found")
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

func (c *OrderCache) Get(key string) (dto.Order, error) {
	const op = "cache.Get()"

	order, ok := c.store.GetIfPresent(key)
	if !ok {
		return dto.Order{}, e.Wrap(op, ErrOrderNotFound)
	}

	return order, nil
}
