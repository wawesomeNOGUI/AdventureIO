package main

// All generic structs and methods for players (Player implements EntityInterface)

type Player struct {
	EntityBase
	BeingHeld string
	roomChangeChan chan *Room
}

func newPlayer(tag string, x, y float64) *Player {
	p := Player{}
	p.X = x 
	p.Y = y
	p.width = 3
	p.height = 3
	p.K = "p"
	p.key = tag
	p.canChangeRooms = true
	p.roomChangeChan = make(chan *Room)

	return &p
}

func (p *Player) Update(oX, oY float64) {
	if p.owner != nil {
		p.X += oX
		p.Y += oY

		if p.held != nil {
			p.held.Update(oX, oY)
		}

		return
	}

	WallCheck(p)
}