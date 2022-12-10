package main

import (
	"sync"
	"math"
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

func (c *ItemContainer) StoreItem(k string, v Item) {
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

func (c *ItemContainer) isItemHere(x, y float64) (bool, string) {
	// c.mu.Lock()
    // defer c.mu.Unlock()

	for k, v := range strayItems.items {
		d := math.Sqrt(math.Pow(x - v.X, 2) + math.Pow(y - v.Y, 2))

		if d < 10 {
			return true, k
		}
	}

	return false, ""
} 

func (c *ItemContainer) TryPickUpItem(o *ItemContainer, pTag string, x, y float64) (bool, string) {
	c.mu.Lock()
    defer c.mu.Unlock()

	itemHere, itemKey := c.isItemHere(x, y)

	if itemHere {
		tmpItem := c.items[itemKey]
		tmpItem.Owner = pTag
		o.StoreItem(itemKey, tmpItem)

		delete(c.items, itemKey)
		return true, itemKey
	}

	return false, ""
}