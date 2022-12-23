package main 

type Room struct {
	roomKey string  // name of room layout picture (client has all room layouts already, send room number only once TCP)
					// roomKey should be of format r1, r2, r3, etc.
	Entities EntityContainer
	updateFunc func(*Room)
	leftX float64 // simple rectangle room bounds
	rightX float64
	upperY float64
	lowerY float64

	leftRoom  *Room
	rightRoom *Room
	aboveRoom *Room
	belowRoom *Room
}

func newRoom(key string, uF func(*Room), lX, rX, uY, lY float64, l, r, u, d *Room) (string, *Room) {
	room := Room{}
	room.Entities = EntityContainer{entities: make(map[string]EntityInterface)} 
	room.updateFunc = uF
	
	room.leftX = lX
	room.rightX = rX
	room.upperY = uY
	room.lowerY = lY

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