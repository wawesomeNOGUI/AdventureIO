package main

import (
	"sync"
	"fmt"
)

type EntityInterface interface {
	Update()
}

// recursively traverse all entities being held by the caller entity 
func (e EntityInterface) TraverseEntities(output map[string]EntityInterface) {
	if e.held != nil {
		e.held.TraverseEntities(output)
		output[e.key] = e
	}
}

// For storing pointers to entities with mutex so goroutines can access with no contention
type EntityContainer struct {
	mu sync.Mutex
	entities map[string]EntityInterface
}

// All entities will have at least this info
type EntityBase struct {
	key string	// this entity's unique key for inserting into map
	X float64
	Y float64
	s float64	// speed, how much can move each update (not exported)
	vX float64  // current direction vectors (normalized = hypotenuse of 1) 
	vY float64
	Kind string // what kind of entity
	held EntityInterface // should be a pointer reference to an entity, this entity will only be accessed through parent entity
}

func (c *EntityContainer) LoadEntity(k string) EntityInterface {
	c.mu.Lock()
    defer c.mu.Unlock()

	return c.entities[k]
}

func (c *EntityContainer) StoreEntity(k string, v EntityInterface) {
	c.mu.Lock()
    defer c.mu.Unlock()

	c.entities[k] = v
}

func (c *EntityContainer) DeleteEntity(k string) EntityInterface {
	c.mu.Lock()
    defer c.mu.Unlock()

	tmpEntity := c.entities[k]
	delete(c.entities, k)

	return tmpEntity
}

// Return map of all entities currently contained in the EntityContainer
func (c *EntityContainer) Entities() map[string]EntityInterface {
	c.mu.Lock()
    defer c.mu.Unlock()

	tmpMap := make(map[string]EntityInterface)
	for k, v := range c.entities {
		tmpMap[k] = v
	}

	return tmpMap
}

// Return serializtion of all entities ready for sending to clients
func (c *EntityContainer) SerializeEntities() string {
	c.mu.Lock()
    defer c.mu.Unlock()

	tmpMap := make(map[string]EntityInterface)
	for _, v := range c.entities {
		v.TraverseEntities(tmpMap)
	}

	jsonTemp, err := json.Marshal(tmpMap)
	if err != nil {
		fmt.Println(err)
	}

	return string(jsonTemp)
}

// Calls each entity's Update function
func (c *EntityContainer) UpdateEntities() {
	c.mu.Lock()
    defer c.mu.Unlock()

	for _, v := range c.entities {
		v.Update()
	}
}

// returns entity key, and normalized vector pointing from x,y to the entity
// closeParam tells distance when the search should break cause found a close enough entity
// d is the largest search radius
func (c *EntityContainer) ClosestEntity(closeParam, d, x, y float64) (string, float64, float64) {
	c.mu.Lock()
    defer c.mu.Unlock()

	var closest string
	var dX float64
	var dY float64

	for k, v := range c.entities {
		tmpDX := v.X - x 
		tmpDY := v.Y - y
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

	return closest, dX/d, dY/d
}