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
	s float64   // speed, how much can move each update (not exported)
	Kind string // what kind of entity
}

//==================Bats=======================
type Bat struct {
	EntityBase
	Held string // bats can pick up items
}

var numOfBats int
func newBat(x, y float64) (string, *Bat) {
	b := Bat{}
	b.X = x 
	b.Y = y
	b.s = 0.25
	b.Kind = "bat"	

	numOfBats++
	return fmt.Sprintf("bat%d", numOfBats), &b
}

func (b *Bat) behaviorFunc() {
	b.X += b.s
	b.Y += b.s

	if b.X < 2 {
		b.X = 2
		b.s = -b.s
	} else if b.X > 154 {
		b.X = 154
	}	
	
	if b.Y < 2 {
		b.Y = 2
	} else if b.Y > 99 {
		b.Y = 99
	}
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