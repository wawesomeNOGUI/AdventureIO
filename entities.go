package main

// Structs and Movement Functions For Game Entities

type Entity struct {
	X float64
	Y float64
	Kind string // what kind of entity
	Held string // key of item the entity holds
	behaviorFunc func()
}

//==================Bats=======================
func batBehaviorFunc() {

}

//==================Dragons====================
func dragonBehaviorFunc() {

}
