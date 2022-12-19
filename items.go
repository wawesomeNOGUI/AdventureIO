package main

import (
	"math"
	"fmt"
)

// Generic structs and methods for items (Items implement EntityInterface)
// Items are entities that can be picked up by other entities

type Item struct {
	EntityBase
}

var numOfItems int
func newItem(kind string, x, y float64) (string, *Item) {
	i := Item{}
	i.X = x 
	i.Y = y
	i.Kind = kind

	numOfItems++
	i.key = fmt.Sprintf(kind + "%d", numOfItems)

	return i.key, &i
}

func (b *Item) Update() {
	b.X += b.vX * b.s
	b.Y += b.vY * b.s

	// WallCheck:

	if b.X < 2 {
		b.X = 2
		b.vX = -b.vX
	} else if b.X > 154 {
		b.X = 154
		b.vX = -b.vX
	}	
	
	if b.Y < 2 {
		b.Y = 2
		b.vY = -b.vY
	} else if b.Y > 99 {
		b.Y = 99
		b.vY = -b.vY
	}
}

func (c *EntityContainer) isItemHere(x, y float64) (bool, string) {
	// c.mu.Lock()
    // defer c.mu.Unlock()

	for k, v := range c.entities {
		_, ok := v.(*Item)  // check if EntityInterface holds type Item
		if !ok {
			continue
		}

		d := math.Sqrt(math.Pow(x - v.GetX(), 2) + math.Pow(y - v.GetY(), 2))

		if d < 10 {
			return true, k
		}
	}

	return false, ""
} 

// TryPickUpItem allows an entity to request trying to pick up an item.
// The function uses EntityContainer mutex so only one goroutine can try to pick up an item at a time (concurrent safe)
// If the entity can pick up the item, set the entity's held member to the pointer to the Item
// put the item in the ownedItems map (parameter o) and then return true, the Item key
func (c *EntityContainer) TryPickUpItem(ref EntityInterface, x, y float64) (bool, string) {
	c.mu.Lock()
    defer c.mu.Unlock()

	itemHere, itemKey := c.isItemHere(x, y)

	if itemHere {
		ref.SetHeld(c.DeleteEntity(itemKey))
	}

	return false, ""
}