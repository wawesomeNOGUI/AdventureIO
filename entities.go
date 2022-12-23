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
	b.width = 8
	b.height = 11
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
func (b *Bat) Update(oX, oY float64) {
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
			
			// b.held.Update(b.vX * 5, b.vY * 5)
			//b.held.SetX(b.held.GetX() + b.vX * 5)
			//b.held.SetY(b.held.GetY() + b.vY * 5)

			// non concurrent safe store here is ok cause UpdateEntities() locks mutex
			b.room.Entities.entities[b.held.Key()] = b.held
			b.held.SetOwner(nil)
			p, ok := b.held.(*Player)
			if ok {
				p.BeingHeld = ""
			}
			b.held = nil
			b.heldCounter = 0
			b.waitCounter = waitCounterThreshold
		}
	} else if b.waitCounter--; b.waitCounter < 0 {	// chase items

		// we can run the non concurrent safe one here cause UpdateEntities() locks the mutex to the map of entities
		// itemKey, vX, vY := b.room.Entities.nonConcurrentSafeClosestItem(b.key, 20, 100, b.X, b.Y)
		itemKey, vX, vY := b.room.Entities.nonConcurrentSafeClosestEntity(b.key, 20, 100, b.X, b.Y)		

		if itemKey != "" {
			b.vX = vX
			b.vY = vY
		}

		// fmt.Println(b.vX)

		// Try to pick up an item
		// b.room.Entities.nonConcurrentSafeTryPickUpItem(b, b.X+2, b.Y+2)
		gotEntity, _ := b.room.Entities.nonConcurrentSafeTryPickUpEntity(b, b.X+2, b.Y+2)
		if gotEntity {
			p, ok := b.held.(*Player)
			if ok {
				p.BeingHeld = b.key
			}
		}
	}

	// fmt.Println("flap")

	WallCheck:
	WallCheck(b)
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