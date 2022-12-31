package main

import (
	"sync"
	"fmt"
	"encoding/json"
	"math"
)

type EntityInterface interface {
	Update(float64, float64)

	Held() EntityInterface    // This method and below implemented by EntityBase
	SetHeld(EntityInterface)
	Owner() EntityInterface
	SetOwner(EntityInterface)
	Key() string
	GetKind() string
	GetX() float64
	SetX(float64)
	GetY() float64
	SetY(float64)
	GetWidth() float64
	SetWidth(float64)
	GetHeight() float64
	SetHeight(float64)
	GetvX() float64
	SetvX(float64)
	GetvY() float64
	SetvY(float64)
	GetS() float64
	SetS(float64)
	GetRoom() *Room 
	SetRoom(*Room)
	CanChangeRooms() bool
	SetCanChangeRooms(bool)
}

// recursively traverse all entities being held by the caller entity 
func traverseEntities(e EntityInterface, output map[string]EntityInterface) {
	if e.Held() != nil {
		traverseEntities(e.Held(), output)
	}

	output[e.Key()] = e
}

func changeRoom(e EntityInterface, r *Room) {
	delete(e.GetRoom().Entities.entities, e.Key())
	e.SetRoom(r)
	r.Entities.StoreEntity(e.Key(), e)
}

func WallCheck(e EntityInterface) {
	if e.GetX() <= 0 {
		if e.GetRoom().leftRoom != nil && e.CanChangeRooms() {
			changeRoom(e, e.GetRoom().leftRoom)
			prevX := e.GetX()
			e.SetX(160 - e.GetWidth() - 1)

			if e.Held() != nil {
				e.Held().Update(e.GetX()-prevX, 0)
			}

			p, ok := e.(*Player)
			if ok {
				reliableChans.SendToPlayer(p.key, p.room.roomKey + "," + p.room.wallColor)
				reliableChans.SendToPlayer(p.key, fmt.Sprintf("P%f,%f", p.X, p.Y))
				p.roomChangeChan <- p.room
			}
		} else {
			e.SetX(0)
			e.SetvX(-e.GetvX())
		}
	} else if e.GetX() + e.GetWidth() >= 160 {
		if e.GetRoom().rightRoom != nil && e.CanChangeRooms() {
			changeRoom(e, e.GetRoom().rightRoom)
			prevX := e.GetX()
			e.SetX(1)

			if e.Held() != nil {
				e.Held().Update(e.GetX()-prevX, 0)
			}

			p, ok := e.(*Player)
			if ok {
				reliableChans.SendToPlayer(p.key, p.room.roomKey + "," + p.room.wallColor)
				reliableChans.SendToPlayer(p.key, fmt.Sprintf("P%f,%f", p.X, p.Y))
				p.roomChangeChan <- p.room
			}
		} else {
			e.SetX(160 - e.GetWidth())
			e.SetvX(-e.GetvX())
		}
	}	
	
	if e.GetY() <= 0 {
		if e.GetRoom().aboveRoom != nil && e.CanChangeRooms() {
			changeRoom(e, e.GetRoom().aboveRoom)
			prevY := e.GetY()
			e.SetY(105 - e.GetHeight() - 1)

			if e.Held() != nil {
				e.Held().Update(0, e.GetY()-prevY)
			}

			p, ok := e.(*Player)
			if ok {
				reliableChans.SendToPlayer(p.key, p.room.roomKey + "," + p.room.wallColor)
				reliableChans.SendToPlayer(p.key, fmt.Sprintf("P%f,%f", p.X, p.Y))
				p.roomChangeChan <- p.room
			}
		} else {
			e.SetY(0)
			e.SetvY(-e.GetvY())
		}
	} else if e.GetY() + e.GetHeight() >= 105 {
		if e.GetRoom().belowRoom != nil && e.CanChangeRooms() {
			changeRoom(e, e.GetRoom().belowRoom)
			prevY := e.GetY()
			e.SetY(1)

			if e.Held() != nil {
				e.Held().Update(0, e.GetY()-prevY)
			}

			p, ok := e.(*Player)
			if ok {
				reliableChans.SendToPlayer(p.key, p.room.roomKey + "," + p.room.wallColor)
				reliableChans.SendToPlayer(p.key, fmt.Sprintf("P%f,%f", p.X, p.Y))
				p.roomChangeChan <- p.room
			}
		} else {
			e.SetY(105 - e.GetHeight())
			e.SetvY(-e.GetvY())
		}
	}

	wallHit := false

	// Pixel perfect wall hit check
	for x := e.GetX() + 1; x < e.GetX() + e.GetWidth() - 1; x++ {  // check for top hit
			// make sure x in range
			if x < 0 || x >= 160 {
				continue
			} 
			// test for hit 
			if !e.GetRoom().layout[int(x)][int(e.GetY())] {
				continue
			}

			e.SetY(e.GetY()+1)
			e.SetvY(-e.GetvY())
			wallHit = true
			break
	}

	for x := e.GetX() + 1; x < e.GetX() + e.GetWidth() - 1; x++ {  // check for bottom hit
		// make sure x in range
		if x < 0 || x >= 160  || int(e.GetY() + e.GetHeight()) >= 105 {
			continue
		} 
		// test for hit 
		if !e.GetRoom().layout[int(x)][int(e.GetY() + e.GetHeight()) - 1] {
			continue
		}

		e.SetY(e.GetY()-1)
		e.SetvY(-e.GetvY())
		wallHit = true
		break
	}

	for y := e.GetY() + 1; y < e.GetY() + e.GetHeight() - 1; y++ {  // check for left hit
		// make sure y in range
		if y < 0 || y >= 105 {
			continue
		} 
		// test for hit 
		if !e.GetRoom().layout[int(e.GetX())][int(y)] {
			continue
		}

		e.SetX(e.GetX()+1)
		e.SetvX(-e.GetvX())
		wallHit = true
		break
	}

	for y := e.GetY() + 1; y < e.GetY() + e.GetHeight() - 1; y++ {  // check for right hit
		// make sure y in range
		if y < 0 || y >= 105  || int(e.GetX() + e.GetWidth()) >= 160 {
			continue
		} 
		// test for hit 
		if !e.GetRoom().layout[int(e.GetX() + e.GetWidth()) - 1][int(y)] {
			continue
		}

		e.SetX(e.GetX()-1)
		e.SetvX(-e.GetvX())
		wallHit = true
		break
	}	

	// make sure player knows they hit the wall and change they're position if they haven't already
	if wallHit {
	// 	p, ok := e.(*Player)
	// 	if ok {
	// 		reliableChans.SendToPlayer(p.key, fmt.Sprintf("P%f,%f", p.X, p.Y))
	// 	}
	}
}

// https://tip.golang.org/doc/go1.8#mapiter
// https://go.dev/blog/race-detector
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

	tmpMap := make(map[string]EntityInterface)
	for _, v := range c.entities {
		traverseEntities(v, tmpMap)
	}

	tmpPlayerMap := make(map[string]*Player)
	for k, v := range tmpMap {
		switch z := v.(type) {
		case *Player:
			tmpPlayerMap[k] = z
		}
	}

	return tmpPlayerMap
}

// Return serializtion of all entities ready for sending to clients
func (c *EntityContainer) SerializeEntities() string {
	c.mu.Lock()
    defer c.mu.Unlock()

	tmpMap := make(map[string]EntityInterface)
	for _, v := range c.entities {
		traverseEntities(v, tmpMap)
	}

	jsonTemp, err := json.Marshal(tmpMap)
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(string(jsonTemp))


	return string(jsonTemp)
}

// Calls each entity's Update function
func (c *EntityContainer) UpdateEntities() {
	c.mu.Lock()
    defer c.mu.Unlock()

	for _, v := range c.entities {
		v.Update(0, 0)
	}
}

func (c *EntityContainer) isEntityHere(self EntityInterface, x, y float64) (bool, string) {
	// c.mu.Lock()
    // defer c.mu.Unlock()

	for k, v := range c.entities {
		if v == self {
			continue
		}
		// to not let entities pick up same type as themselves (kinda fun interaction though maybe make an option)
		if v.GetKind() == self.GetKind() {
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


func (c *EntityContainer) nonConcurrentSafeTryPickUpEntity(ref EntityInterface, x, y float64) (bool, string) {
	entityHere, entityKey := c.isEntityHere(ref, x, y)

	if entityHere {
		ref.SetHeld(c.entities[entityKey])
		c.entities[entityKey].SetOwner(ref)
		delete(c.entities, entityKey)
		return true, entityKey
	}

	return false, ""
}

// run the below one for concurrent safe calling
func (c *EntityContainer) nonConcurrentSafeClosestEntity(self string, filter []string, closeParam, d, x, y float64) (string, float64, float64) {
	var closest string
	var dX float64
	var dY float64

	for k, v := range c.entities {
		if self == k {
			continue
		}
		// to not let entities pick up same type as themselves (kinda fun interaction though maybe make an option)
		if v.GetKind() == c.entities[self].GetKind() {
			continue
		}

		if len(filter) == 0 {
			goto NEXT
		}

		for _, f := range filter {
			if v.GetKind() == f {
				goto NEXT
			} 		
		}
		continue
		
		NEXT:

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

// returns entity key, and normalized vector pointing from x,y to the entity
// closeParam tells distance when the search should break cause found a close enough entity
// d is the largest search radius
func (c *EntityContainer) ClosestEntity(self string, filter []string, closeParam, d, x, y float64) (string, float64, float64) {
	c.mu.Lock()
    defer c.mu.Unlock()

	itemKey, vX, vY := c.nonConcurrentSafeClosestEntity(self, filter, closeParam, d, x, y)

	return itemKey, vX, vY
}