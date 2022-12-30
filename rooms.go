package main 

import (
	"fmt"
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
	// leftX float64 // simple rectangle room bounds
	// rightX float64
	// upperY float64
	// lowerY float64

	layout *[160][105]bool

	leftRoom  *Room
	rightRoom *Room
	aboveRoom *Room
	belowRoom *Room
}

func newRoom(key string, uF func(*Room), rL *[160][105]bool, l, r, u, d *Room) (string, *Room) {
	room := Room{}
	room.Entities = EntityContainer{entities: make(map[string]EntityInterface)} 
	room.updateFunc = uF

	// room.leftX = lX
	// room.rightX = rX
	// room.upperY = uY
	// room.lowerY = lY
	room.layout = rL

	room.leftRoom = l
	room.rightRoom = r
	room.aboveRoom = u
	room.belowRoom = d

	room.roomKey = key

	return key, &room
}

func defaultRoomUpdate(r *Room) {
	r.Entities.UpdateEntities()
}

func castleRoomUpdate(r *Room) {
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
					p.roomChangeChan <- p.room
				}
			} else {
				v.SetY(51)
			}
		}
		v.Update(0, 0)
	}
}