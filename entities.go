package main

import (
	"sync"
	"fmt"
)
// Structs and Movement Functions For Game Entities
/*
type Entity struct {
	X float64
	Y float64
	s float64   // speed, how much can move each update (not exported)
	Kind string // what kind of entity
	Held string // key of item the entity holds
	behaviorFunc func(*Entity)
}

type EntityContainer struct {
	mu sync.Mutex
	entities map[string]Entity
}

func (c *EntityContainer) LoadEntity(k string) Entity {
	c.mu.Lock()
    defer c.mu.Unlock()

	return c.entities[k]
}

func (c *EntityContainer) StoreEntity(k string, v Entity) {
	c.mu.Lock()
    defer c.mu.Unlock()

	c.entities[k] = v
}

func (c *EntityContainer) DeleteEntity(k string) Entity {
	c.mu.Lock()
    defer c.mu.Unlock()

	tmpEntity := c.entities[k]
	delete(c.entities, k)

	return tmpEntity
}

// Return map of all Entities currently contained in the EntityContainer
func (c *EntityContainer) GetEntities() map[string]Entity {
	c.mu.Lock()
    defer c.mu.Unlock()

	tmpMap := make(map[string]Entity)
	for k, v := range c.entities {
		tmpMap[k] = v
	}

	return tmpMap
}
*/
func InitializeEntities(m *sync.Map) {
	// List all the entities you want here
	m.Store(newBat(50, 75))
	m.Store(newBat(50, 6))
}

type EntityInterface interface {
	behaviorFunc()
}

type EntityBase struct {
	X float64
	Y float64
	sX float64   // speed, how much can move each update (not exported)
	sY float64
	Kind string // what kind of entity
}

//==================Bats=======================
type Bat struct {
	EntityBase
	Held string // bats can pick up items
	heldCounter int	// how long the bat has held this item, it drops it after the counter reaches a certain point
}

var numOfBats int
func newBat(x, y float64) (string, *Bat) {
	b := Bat{}
	b.X = x 
	b.Y = y
	b.sX = 1
	b.sY = -1
	b.Kind = "bat"	

	numOfBats++
	return fmt.Sprintf("bat%d", numOfBats), &b
}

func (b *Bat) behaviorFunc() {
	b.X += b.sX
	b.Y += b.sY

	if b.Held != "" {
		tmpItem := ownedItems.LoadItem(b.Held)
		tmpItem.X += b.sX
		tmpItem.Y += b.sY
		ownedItems.StoreItem(b.Held, tmpItem)

		b.heldCounter++

		if b.heldCounter > 100 {
			ownedItems.DeleteItem(b.Held)
			k := b.Held
			tmpItem.Owner = ""
			b.Held = ""
			b.heldCounter = 0

			// fly away from dropped item
			if tmpItem.X > b.X && b.sX > 0 {
				b.sX = -b.sX
			}
			if tmpItem.Y > b.Y && b.sY > 0 {
				b.sY = -b.sY
			}
			
			tmpItem.X -= b.sX * 2
			tmpItem.Y -= b.sY * 2

			strayItems.StoreItem(k, tmpItem)
		}
	}

	if b.X < 2 {
		b.X = 2
		b.sX = -b.sX
	} else if b.X > 154 {
		b.X = 154
		b.sX = -b.sX
	}	
	
	if b.Y < 2 {
		b.Y = 2
		b.sY = -b.sY
	} else if b.Y > 99 {
		b.Y = 99
		b.sY = -b.sY
	}
}

func (b *Bat) tryPickUpItem(batKey string, s, o *ItemContainer) (bool) {
	gotItem, itemKey := s.TryPickUpItem(o, batKey, b.X+2, b.Y+2)
	if gotItem {
		b.Held = itemKey
	}
	return gotItem
}
//==================Dragons====================
/*
func (e *Entity) initializeDragon() {
	e.s = 0.15
	e.behaviorFunc = dragonBehaviorFunc
}

func dragonBehaviorFunc(e *Entity) {

}
*/