package main

// All generic structs and methods for players (Player implements EntityInterface)

type Player struct {
	EntityBase
	BeingHeld string
}

func newPlayer(tag string, x, y float64) *Player {
	p := Player{}
	p.X = x 
	p.Y = y
	p.K = "p"
	p.key = tag

	return &p
}

func (p *Player) Update() {
	if p.owner != nil {
		p.X = p.owner.GetX()
		p.Y = p.owner.GetY()
	}

	if p.held != nil {
		p.held.Update()
	}
}