package main

// All generic structs and methods for players

type Player struct {
	X float64
	Y float64
	Held string  // string representing the key for an item in ownedItems
}