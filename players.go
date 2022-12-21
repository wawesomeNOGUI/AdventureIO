package main

// All generic structs and methods for players (Player implements EntityInterface)

type Player struct {
	EntityBase
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
	// dummy update function to implement EntityInterface
}