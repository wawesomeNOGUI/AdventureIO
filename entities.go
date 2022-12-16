package main

import "sync"

// Structs and Movement Functions For Game Entities

type Entity struct {
	X float64
	Y float64
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

//==================Bats=======================
func batBehaviorFunc(e *Entity) {
	e.X++
	e.Y++

	if e.X < 2 {
		e.X = 2
	} else if e.X > 154 {
		e.X = 154
	}	
	
	if e.Y < 2 {
		e.Y = 2
	} else if e.Y > 99 {
		e.Y = 99
	}
}

//==================Dragons====================
func dragonBehaviorFunc(e *Entity) {

}
