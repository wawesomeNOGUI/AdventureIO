package main

import (
	"sync"
	"fmt"
	"encoding/json"
	"math"
)

type EntityInterface interface {
	Update()

	Held() EntityInterface    // This method and below implemented by EntityBase
	SetHeld(EntityInterface)
	Key() string
	GetKind() string
	GetX() float64
	SetX(float64)
	GetY() float64
	SetY(float64)
	GetRoom() *Room 
}

// recursively traverse all entities being held by the caller entity 
func TraverseEntities(e EntityInterface, output map[string]EntityInterface) {
	if e.Held() != nil {
		TraverseEntities(e.Held(), output)
	}

	output[e.Key()] = e
}

// For storing pointers to entities with mutex so goroutines can access with no contention
type EntityContainer struct {
	mu sync.Mutex
	entities map[string]EntityInterface
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

// Return map of all player entities currently contained in the EntityContainer
func (c *EntityContainer) Players() map[string]*Player {
	c.mu.Lock()
    defer c.mu.Unlock()

	tmpMap := make(map[string]*Player)
	for k, v := range c.entities {
		switch z := v.(type) {
		case *Player:
			tmpMap[k] = z
		}
	}

	return tmpMap
}

// Return serializtion of all entities ready for sending to clients
func (c *EntityContainer) SerializeEntities() string {
	c.mu.Lock()
    defer c.mu.Unlock()

	tmpMap := make(map[string]EntityInterface)
	for _, v := range c.entities {
		TraverseEntities(v, tmpMap)
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

	return closest, dX/d, dY/d
}