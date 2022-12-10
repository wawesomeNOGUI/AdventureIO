package main

import (
	"sync"
	"fmt"
)

// Generic structs and methods for items

type Item struct {
	X float64
	Y float64
	Owner string  // should be a playerTag
	Kind string   // tells what kind of item it is
}

// For storing lists of items with mutex so goroutines can access with no contention
type ItemContainer struct {
	mu sync.Mutex
	items map[string]Item
}

func (c *ItemContainer) AddItem(k string, v Item) {
	c.mu.Lock()
    defer c.mu.Unlock()

	c.items[k] = v
}

func (c *ItemContainer) DeleteItem(k string) {
	c.mu.Lock()
    defer c.mu.Unlock()

	delete(c.items, k)
}

func (c *ItemContainer) GetItems() map[string]Item {
	c.mu.Lock()
    defer c.mu.Unlock()

	tmpMap := make(map[string]Item)
	for k, v := range strayItems.items {
		tmpMap[k] = v
	}

	return tmpMap
}