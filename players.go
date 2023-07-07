package main

import (
	"sync"
)

// All generic structs and methods for players (Player implements EntityInterface)

type Player struct {
	EntityBase
	BeingHeld string
	mu sync.Mutex
}

func newPlayer(tag string, x, y float64) *Player {
	p := Player{}
	p.X = x 
	p.Y = y
	p.width = 4
	p.height = 4
	p.K = "p"
	p.key = tag
	p.canChangeRooms = true
	// p.roomChangeChan = make(chan *Room)
	// p.canUpdate = make(chan bool)
	// p.updateDone = make(chan bool)

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

	if p.held != nil {
		//mostly moved in player connecition in main, but this updates if touching other entities
		p.held.Update(oX, oY)
	}

	wallCheck(p)
}