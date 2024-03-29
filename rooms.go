package main 

import (
	"fmt"
	//"time"
)

type Vector2 struct {
	x float64
	y float64
}

type Room struct {
	roomKey string  // name of room layout picture (client has all room layouts already, send room number only once TCP)
					// roomKey should be of format r1, r2, r3, etc.
	wallColor string // to tell player what color they can't walk through

	Entities EntityContainer
	updateFunc func(*Room)
	updateLoop func(*Room)
	updateEntitiesChan chan bool
	sendGameStateUnreliableChan chan bool
	// leftX float64 // simple rectangle room bounds
	// rightX float64
	// upperY float64
	// lowerY float64

	layout *[160][105]bool

	leftRoom  *Room
	rightRoom *Room
	aboveRoom *Room
	belowRoom *Room

	specialVars map[string]interface{}
}

func newRoom(key string, uF func(*Room), rL *[160][105]bool, l, r, u, d *Room) (string, *Room) {
	room := Room{}
	room.roomKey = key
	room.Entities = EntityContainer{entities: make(map[string]EntityInterface)} 
	room.updateFunc = uF
	room.updateLoop = roomUpdateLoop
	room.updateEntitiesChan = make(chan bool)
	room.sendGameStateUnreliableChan = make(chan bool)

	// room.leftX = lX
	// room.rightX = rX
	// room.upperY = uY
	// room.lowerY = lY
	room.layout = rL

	room.leftRoom = l
	room.rightRoom = r
	room.aboveRoom = u
	room.belowRoom = d

	room.specialVars = make(map[string]interface{})

	//go room.updateLoop(&room)

	return key, &room
}

func roomUpdateLoop(r *Room) {
	for {
		// check for if it's time to update entities or send game state
		select {
		case <-r.updateEntitiesChan:
			r.updateFunc(r)
		case <-r.sendGameStateUnreliableChan:
			s := r.Entities.SerializeEntities()

			for k, _ := range r.Entities.Players() {
				unreliableChans.SendToPlayer(k, s)
			}
		}
	}
}

func defaultRoomUpdate(r *Room) {
	r.Entities.UpdateEntities()
}

func respawnRoomUpdate(r *Room) {
	//
}

const fallAsleepThreshold = 500
func batRoomUpdate(r *Room) {
	if r.specialVars["batsAwake"] == false {
		// prevent player goroutines from sending updates updating
		for _, p := range r.Entities.Players() {
			p.mu.Lock()
			defer p.mu.Unlock()
		}

		r.Entities.mu.Lock()
    	defer r.Entities.mu.Unlock()

		r.specialVars["fallAsleepTimer"] = 0  // reset timer

		for _, v := range r.Entities.entities {
			if v.GetKind() == "bat" {
				// check if player touches this bat, waking it up
				yes, _ := r.Entities.isEntityHere(v, []string{"p"}, v.GetX(), v.GetY())
				if yes {
					r.specialVars["batsAwake"] = true
				}
			} else {
				v.Update(0, 0)
			}
		}
	} else {
		if len(r.Entities.Players()) <= 0 {
			t := r.specialVars["fallAsleepTimer"].(int)
			t++
			r.specialVars["fallAsleepTimer"] = t

			if t > fallAsleepThreshold {
				r.specialVars["batsAwake"] = false
			}
		}

		r.Entities.UpdateEntities()
	}
}


func dragonRoomUpdate(r *Room) {
	playersPresent := len(r.Entities.Players())
	doorGrateMap := r.Entities.GetEntitiesByKind("dG")

	r.Entities.mu.Lock()

	if r.specialVars["dragonBeat"] == false {
		// close gate when players enter
		if playersPresent > 0 {
			
	
			for _, doorGrate := range doorGrateMap {
				if doorGrate.GetX() == 57 && doorGrate.GetY() == 3 {
					doorGrate.SetvX(0)
					doorGrate.SetvY(0)
					continue
				}

				if doorGrate.GetX() < 57 {
					doorGrate.SetvX(0.25)
				} else if doorGrate.GetX() > 57 {
					doorGrate.SetvX(-0.25)
				}
		
				if doorGrate.GetX() == 57 && doorGrate.GetY() != 3 {
					doorGrate.SetvX(0)
					doorGrate.SetvY(-0.25)
				}
			}	
		}
	} else {
		// open gate when dragon defeated
		if playersPresent > 0 {
	
			for _, doorGrate := range doorGrateMap {
				if doorGrate.GetX() == 10 && doorGrate.GetY() == 10 {
					doorGrate.SetvX(0)
					doorGrate.SetvY(0)
					continue
				}

				if doorGrate.GetX() < 10 {
					doorGrate.SetvX(0.25)
				} else if doorGrate.GetX() > 10 {
					doorGrate.SetvX(-0.25)
				} else {
					doorGrate.SetvX(0)
				}
				
				if doorGrate.GetY() > 10 {
					doorGrate.SetvY(-0.25)
				} else if doorGrate.GetY() < 10 {
					doorGrate.SetvY(0.25)
				} else {
					doorGrate.SetvY(0)
				}
			}	
		}
	}

	r.Entities.mu.Unlock()
	
	r.Entities.UpdateEntities()
}

func castleRoomUpdate(r *Room) {
	// prevent player goroutines from sending updates updating
	for _, p := range r.Entities.Players() {
		p.mu.Lock()
		defer p.mu.Unlock()
	}

	r.Entities.mu.Lock()
    defer r.Entities.mu.Unlock()

	for _, v := range r.Entities.entities {
		if v.GetY() < 50 && v.GetX() > 49 && v.GetX() < 110 - v.GetWidth() {
			prevY := v.GetY()
			prevX := v.GetX()
			v.SetY(53) 

			if v.GetX() < 75 {
				v.SetX(75)
			} else if v.GetX() > 85 - v.GetWidth() {
				v.SetX(85 - v.GetWidth())
			}
			if v.Held() != nil {
				v.Held().Update(v.GetX() - prevX, v.GetY() - prevY)
			}

			p, ok := v.(*Player)
			if ok {
				reliableChans.SendToPlayer(p.key, fmt.Sprintf("P%f,%f", p.X, p.Y))
			}
		} else if v.GetY() <= 52 && v.GetX() > 73 && v.GetX() <= 86 - v.GetWidth() {
			if v.GetRoom().aboveRoom != nil && v.CanChangeRooms() {
				changeRoom(v, v.GetRoom().aboveRoom)
				prevY := v.GetY()
				v.SetY(105 - v.GetHeight() - 1)
	
				if v.Held() != nil {
					v.Held().Update(0, v.GetY()-prevY)
				}
	
				p, ok := v.(*Player)
				if ok {
					reliableChans.SendToPlayer(p.key, p.room.roomKey + "," + p.room.wallColor)
					reliableChans.SendToPlayer(p.key, fmt.Sprintf("P%f,%f", p.X, p.Y))
					//p.roomChangeChan <- p.room
				}
			} else {
				v.SetY(51)
			}
		}
		v.Update(0, 0)
	}
}