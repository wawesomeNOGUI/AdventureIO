package main

// All generic structs and interfaces for items, players, etc.

type Item struct {
	X float64
	Y float64
	Owner string  // should be a playerTag
	Kind string   // tells what kind of item it is
}

type Player struct {
	X float64
	Y float64
	Held Item
}