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
func newItem(kind string, r *Room, x, y, w, h float64) (string, *Item) {
	i := Item{}
	i.X = x 
	i.Y = y
	i.width = w
	i.height = h
	i.K = kind
	i.room = r
	i.canChangeRooms = true

	numOfItems++
	i.key = fmt.Sprintf(kind + "%d", numOfItems)

	return i.key, &i
}

func newSword(r *Room, x, y float64) (string, *Item) {
	return newItem("sword", r, x, y, 10, 5)
}

func newDoorGrate(r *Room, x, y float64) (string, *Item) {
	return newItem("dG", r, x, y, 46, 8)
}

func (b *Item) Update(oX, oY float64) {
	if b.owner != nil {
		b.X += oX
		b.Y += oY

		if b.held != nil {
			b.held.Update(oX, oY)
		}

		return
	}

	prevX := b.X
	prevY := b.Y
	b.X += b.vX * b.s
	b.Y += b.vY * b.s

	if b.held != nil {
		b.held.Update(b.X-prevX, b.Y-prevY)
	}

	// WallCheck:
	WallCheck(b)
}

func (c *EntityContainer) isItemHere(self EntityInterface, x, y float64) (bool, string) {
	// c.mu.Lock()
    // defer c.mu.Unlock()

	for k, v := range c.entities {
		if v == self {
			continue
		}
		_, ok := v.(*Item)  // check if EntityInterface holds type Item
		if !ok {
			continue
		}

		// check for rectangle overlap between two entities
		if x < v.GetX() + v.GetWidth() && v.GetX() < x + self.GetWidth() {
			if y < v.GetY() + v.GetHeight() && v.GetY() < y + self.GetHeight() {
				return true, k
			}
		}
	}

	return false, ""
} 

// TryPickUpItem allows an entity to request trying to pick up an item.
// The function uses EntityContainer mutex so only one goroutine can try to pick up an item at a time (concurrent safe)
// If the ref entity can pick up the item, set the ref entity's held member to the pointer to the Item
// and then return true, the Item key
func (c *EntityContainer) TryPickUpItem(ref EntityInterface, x, y float64) (bool, string) {
	c.mu.Lock()
    defer c.mu.Unlock()

	gotItem, itemKey := c.nonConcurrentSafeTryPickUpItem(ref, x, y)

	return gotItem, itemKey
}

func (c *EntityContainer) nonConcurrentSafeTryPickUpItem(ref EntityInterface, x, y float64) (bool, string) {
	itemHere, itemKey := c.isItemHere(ref, x, y)

	if itemHere {
		ref.SetHeld(c.entities[itemKey])
		c.entities[itemKey].SetOwner(ref)
		delete(c.entities, itemKey)
		return true, itemKey
	}

	return false, ""
}

// run the below one for concurrent safe calling
func (c *EntityContainer) nonConcurrentSafeClosestItem(self string, closeParam, d, x, y float64) (string, float64, float64) {
	var closest string
	var dX float64
	var dY float64

	for k, v := range c.entities {
		if self == k {
			continue
		}
		_, ok := v.(*Item)  // check if EntityInterface holds type Item
		if !ok {
			continue
		}

		tmpDX := v.GetX() - x 
		tmpDY := v.GetY() - y
		tmpD := math.Sqrt((tmpDX)*(tmpDX) + (tmpDY)*(tmpDY))


		if tmpD < d {
			closest = k
			d = tmpD
			dX = tmpDX
			dY = tmpDY
		}

		if d < closeParam {  // so don't have to search through all entity just return early after finding pretty close one
			break
		}
	}
	if d != 0 {
		return closest, dX/d, dY/d
	} else {
		return "", 0, 0
	}
}

