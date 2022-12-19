package main

// All generic structs and methods for players (Player implements EntityInterface)

type Player struct {
	EntityBase
}

func (p *Player) Update() {
	// dummy update function to implement EntityInterface
}