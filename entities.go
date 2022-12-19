package main


// all entities implement EntityInterface 
// this file has monsters, animals, things that move

//==================Bats=======================
type Bat struct {
	EntityBase
	Held string // bats can pick up items
	heldCounter int	// how long the bat has held this item, it drops it after the counter reaches a certain point
	waitCounter int // time delay before allowed to fly towards items again
}

var numOfBats int
func newBat(x, y float64) (string, *Bat) {
	b := Bat{}
	b.X = x 
	b.Y = y
	b.s = 1
	b.Kind = "bat"	

	numOfBats++
	b.key = fmt.Sprintf("bat%d", numOfBats)

	return b.key, &b
}

var heldCounterThreshold int = 100
var waitCounterThreshold int = 200
func (b *Bat) Update() {
	b.X += b.vX * b.s
	b.Y += b.vY * b.s

	if b.Held != "" {
		tmpItem := ownedItems.LoadItem(b.Held)
		tmpItem.X += b.vX * b.s
		tmpItem.Y += b.vY * b.s
		ownedItems.StoreItem(b.Held, tmpItem)

		b.heldCounter++

		if b.heldCounter > heldCounterThreshold {
			if b.X < 15 || b.X > 145 {
				goto WallCheck
			} else if b.Y < 15 || b.Y > 90 {
				goto WallCheck
			}
			ownedItems.DeleteItem(b.Held)
			k := b.Held
			tmpItem.Owner = ""
			b.Held = ""
			b.heldCounter = 0
			b.waitCounter = waitCounterThreshold

			// fly away from dropped item
			if tmpItem.X > b.X && b.vX > 0 {
				b.vX = -b.vX
			}
			if tmpItem.Y > b.Y && b.vY > 0 {
				b.vY = -b.vY
			}
			
			tmpItem.X -= b.vX * 5
			tmpItem.Y -= b.vY * 5

			strayItems.StoreItem(k, tmpItem)
		}
	} else if b.waitCounter--; b.waitCounter < 0 {	// chase items
		itemKey, vX, vY := strayItems.ClosestItem(20, 100, b.X, b.Y)

		if itemKey != "" {
			b.vX = vX
			b.vY = vY
		}
	}

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

func (b *Bat) tryPickUpItem(batKey string, s, o *EntityInterface) (bool) {
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