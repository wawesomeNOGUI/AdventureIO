package main

import (
	"sync"
	"math"
)

// Generic structs and methods for items

type Item struct {
	X float64
	Y float64
	Owner string  // should be a playerTag or entity key
	Kind string   // tells what kind of item it is
}

// For storing lists of items with mutex so goroutines can access with no contention
type ItemContainer struct {
	mu sync.Mutex
	items map[string]Item
}

func (c *ItemContainer) LoadItem(k string) Item {
	c.mu.Lock()
    defer c.mu.Unlock()

	return c.items[k]
}

func (c *ItemContainer) StoreItem(k string, v Item) {
	c.mu.Lock()
    defer c.mu.Unlock()

	c.items[k] = v
}

func (c *ItemContainer) DeleteItem(k string) Item {
	c.mu.Lock()
    defer c.mu.Unlock()

	tmpItem := c.items[k]
	delete(c.items, k)

	return tmpItem
}

// Return map of all Items currently contained in the ItemContainer
func (c *ItemContainer) GetItems() map[string]Item {
	c.mu.Lock()
    defer c.mu.Unlock()

	tmpMap := make(map[string]Item)
	for k, v := range c.items {
		tmpMap[k] = v
	}

	return tmpMap
}

// returns item key, and normalized vector pointing from x,y to the item
// closeParam tells distance when the search should break cause found a close enough item
// d is the largest search radius
func (c *ItemContainer) ClosestItem(closeParam, d, x, y float64) (string, float64, float64) {
	c.mu.Lock()
    defer c.mu.Unlock()

	var closest string
	var dX float64
	var dY float64

	for k, v := range c.items {
		tmpDX := v.X - x 
		tmpDY := v.Y - y
		tmpD := math.Sqrt((tmpDX)*(tmpDX) + (tmpDY)*(tmpDY))

		if tmpD < d {
			closest = k
			d = tmpD
			dX = tmpDX
			dY = tmpDY
		}

		if d < closeParam {  // so don't have to search through all items just return early after finding pretty close one
			break
		}
	}

	return closest, dX/d, dY/d
}

func (c *ItemContainer) isItemHere(x, y float64) (bool, string) {
	// c.mu.Lock()
    // defer c.mu.Unlock()

	for k, v := range c.items {
		d := math.Sqrt(math.Pow(x - v.X, 2) + math.Pow(y - v.Y, 2))

		if d < 10 {
			return true, k
		}
	}

	return false, ""
} 

// TryPickUpItem allows a player to request trying to pick up an item.
// The function uses ItemContainer mutex so only one player goroutine can try to pick up an item at a time (concurrent safe)
// If the player can pick up the item, set the Item Owned member to the playerTag or entity, 
// put the item in the ownedItems map (parameter o) and then return true, the Item key
func (c *ItemContainer) TryPickUpItem(o *ItemContainer, tag string, x, y float64) (bool, string) {
	c.mu.Lock()
    defer c.mu.Unlock()

	o.mu.Lock()
	defer o.mu.Unlock()

	itemHere, itemKey := c.isItemHere(x, y)

	if itemHere {
		tmpItem := c.items[itemKey]
		tmpItem.Owner = tag
		o.items[itemKey] = tmpItem

		delete(c.items, itemKey)
		return true, itemKey
	}

	return false, ""
}