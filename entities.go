package main

import (
	"fmt"
	// "math"
)

// all entities implement EntityInterface 
// this file has monsters, animals, things that move

//==================Bats=======================
type Bat struct {
	EntityBase
	heldCounter int	// how long the bat has held this item, it drops it after the counter reaches a certain point
	waitCounter int // time delay before allowed to fly towards items again
}

var numOfBats int
func newBat(room *Room, x, y float64) (string, *Bat) {
	b := Bat{}
	b.X = x 
	b.Y = y
	b.vX = 0.5
	b.vY = 0.5
	b.s = 1
	b.K = "bat"	
	b.room = room

	numOfBats++
	b.key = fmt.Sprintf("bat%d", numOfBats)

	return b.key, &b
}

const heldCounterThreshold = 100
const waitCounterThreshold = 200
func (b *Bat) Update() {
	b.X += b.vX * b.s
	b.Y += b.vY * b.s

	if b.held != nil {
		b.held.SetX(b.held.GetX() + b.vX * b.s)
		b.held.SetY(b.held.GetY() + b.vY * b.s)

		b.heldCounter++

		if b.heldCounter > heldCounterThreshold {
			if b.X < 15 || b.X > 145 {
				goto WallCheck
			} else if b.Y < 15 || b.Y > 90 {
				goto WallCheck
			}

			// fly away from dropped item
			if b.held.GetX() > b.X && b.vX > 0 {
				b.vX = -b.vX
			}
			if b.held.GetY() > b.Y && b.vY > 0 {
				b.vY = -b.vY
			}
			
			b.held.SetX(b.held.GetX() + b.vX * 5)
			b.held.SetY(b.held.GetY() + b.vY * 5)

			// non concurrent safe store here is ok cause UpdateEntities() locks mutex
			b.room.Entities.entities[b.held.Key()] = b.held
			b.held = nil
			b.heldCounter = 0
			b.waitCounter = waitCounterThreshold
		}
	} else if b.waitCounter--; b.waitCounter < 0 {	// chase items

		// we can run the non concurrent safe one here cause UpdateEntities() locks the mutex to the map of entities
		itemKey, vX, vY := b.room.Entities.nonConcurrentSafeClosestItem(b.key, 20, 100, b.X, b.Y)

		if itemKey != "" {
			b.vX = vX
			b.vY = vY
		}

		// fmt.Println(b.vX)

		// Try to pick up an item
		b.room.Entities.nonConcurrentSafeTryPickUpItem(b, b.X+2, b.Y+2)
	}

	// fmt.Println("flap")

	WallCheck:

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
//==================Dragons====================
/*
func (e *Entity) initializeDragon() {
	e.s = 0.15
	e.behaviorFunc = dragonBehaviorFunc
}

func dragonBehaviorFunc(e *Entity) {

}
*/