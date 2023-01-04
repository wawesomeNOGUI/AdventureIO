package main

import (
	"fmt"
	// "math"
)

// all entities implement EntityInterface 
// this file has monsters, animals, things that move

type LockedDoor struct {
	EntityBase
	locked bool
	unlockKey EntityInterface
}

var numOfLockedDoors int
func newLockedDoor(room *Room, x, y float64, unlockKey EntityInterface) (string, *LockedDoor)  {
	d := LockedDoor{}
	d.X = x
	d.Y = y
	d.width = 46
	d.height = 8
	d.vX = 0
	d.vY = 0
	d.s = 1
	d.K = "lD"
	d.room = room
	d.canChangeRooms = false

	d.locked = true
	d.unlockKey = unlockKey

	numOfLockedDoors++
	d.key = fmt.Sprintf("lD%d", numOfLockedDoors)

	return d.key, &d
}

func (d *LockedDoor) Update(oX, oY float64) {
	if d.owner != nil {
		d.X += oX
		d.Y += oY

		if d.held != nil {
			d.held.Update(oX, oY)
		}

		return
	}

	if !d.locked {
		// prevX := d.X
		// prevY := d.Y
		d.X += d.vX *d.s
		d.Y += d.vY *d.s
	} else {
		//
	}

	WallCheck(d)
}



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
	b.canChangeRooms = true

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
			// if b.X < 15 || b.X > 145 {
			// 	goto WallCheck
			// } else if b.Y < 15 || b.Y > 90 {
			// 	goto WallCheck
			// }

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
		itemKey, vX, vY := b.room.Entities.nonConcurrentSafeClosestEntity(b.key, []string{}, 20, 100, b.X, b.Y)		

		if itemKey != "" {
			b.vX = vX
			b.vY = vY
		}

		// fmt.Println(b.vX)

		// Try to pick up an item
		// b.room.Entities.nonConcurrentSafeTryPickUpItem(b, b.X+2, b.Y+2)
		gotEntity, _ := b.room.Entities.nonConcurrentSafeTryPickUpEntity(b, b.X, b.Y)
		if gotEntity {
			p, ok := b.held.(*Player)
			if ok {
				p.BeingHeld = b.key
			}
		}
	}

	// fmt.Println("flap")

	// WallCheck:
	WallCheck(b)
}

//==================Dragons====================
type Dragon struct {
	EntityBase
	playersHeld map[string]*Player
	health int
	waitCounter int // time delay before allowed to fly towards players and attack again
	invincibleCounter int //dragon gains invincibility for a short time after being hit 
}

var numOfDragons int
func newDragon(room *Room, x, y float64) (string, *Dragon) {
	b := Dragon{}
	b.X = x 
	b.Y = y
	b.width = 8
	b.height = 20
	b.vX = 0.5
	b.vY = 0.5
	b.s = 1.25
	b.K = "drg"
	b.room = room
	b.canChangeRooms = false
	b.playersHeld = make(map[string]*Player)
	b.health = 5

	numOfDragons++
	b.key = fmt.Sprintf("drg%d", numOfDragons)

	return b.key, &b
}


const drgWaitCounterThreshold = 150
const lungeLength = 30
const dragonInvincibleDelay = 25
func (d *Dragon) Update(oX, oY float64) {
	if d.owner != nil {
		d.X += oX
		d.Y += oY

		for _, v := range d.playersHeld {
			v.Update(oX, oY)
		}

		return
	}

	if d.invincibleCounter > 0 {
		d.invincibleCounter--
	}

	if d.health <= 0 {
		d.room.specialVars["dragonBeat"] = true

		//drop held players
		for k, v := range d.playersHeld {
			d.room.Entities.entities[k] = v
			v.SetOwner(nil)
			v.SetRoom(d.room)
			v.BeingHeld = ""

			delete(d.playersHeld, k)
		}

		//then go to respawn room
		delete(d.room.Entities.entities, d.key)
		RespawnRoomPtr.Entities.StoreEntity(d.key, d)
		return
	}

	prevX := d.X
	prevY := d.Y
	d.X += d.vX * d.s
	d.Y += d.vY * d.s

	for _, v := range d.playersHeld {
		v.Update(d.X-prevX, d.Y-prevY)
	}

	if d.waitCounter++; d.waitCounter > drgWaitCounterThreshold {
		d.waitCounter = 0
		d.s = 3  //so dragon lunges at player 
		// we can run the non concurrent safe one here cause UpdateEntities() locks the mutex to the map of entities
		// itemKey, vX, vY := b.room.Entities.nonConcurrentSafeClosestItem(b.key, 20, 100, b.X, b.Y)
		entityKey, vX, vY := d.room.Entities.nonConcurrentSafeClosestEntity(d.key, []string{"p"}, 20, 100, d.X, d.Y)		

		if entityKey != "" {
			d.vX = vX
			d.vY = vY
		}
	} else if d.waitCounter == lungeLength {
		d.s = 1.25
	}

	yes, key := d.room.Entities.isEntityHere(d, []string{"p"}, d.X, d.Y)
	if yes {
		_, ok := d.room.Entities.entities[key].(*Player)
		if ok {
			gotEntity, _ := d.room.Entities.nonConcurrentSafeTryPickUpEntity(d, d.X, d.Y)
			if gotEntity {
				p, ok := d.held.(*Player)
				if ok {
					d.playersHeld[p.key] = p
					p.BeingHeld = d.key

					d.held = nil
				}
			}
		}
	}

	// WallCheck:
	WallCheck(d)
}