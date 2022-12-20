package main

import "sync"

func InitializeRooms(m *sync.Map) {
	m.Store(newRoom("r1", defaultRoomUpdate, nil, nil, nil, nil))
}

func InitializeEntities(m *sync.Map) {
	// List all the entities you want here
	m.Store(newBat(50, 75))
	m.Store(newBat(50, 6))
	m.Store(newBat(50, 6))
	tmpB, _ := m.Load("bat2")
	tmpB.(*Bat).Kind = "key"
	m.Store("bat2", tmpB)
}