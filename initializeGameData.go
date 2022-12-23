package main

import "sync"

func InitializeRooms(m *sync.Map) {
	m.Store(newRoom("r1", defaultRoomUpdate, &r1Layout, nil, nil, nil, nil))
}

func InitializeEntities(m *sync.Map) {
	// List all the entities you want here
	r, _ := m.Load("r1")
	tmpR := r.(*Room)
	tmpR.Entities.StoreEntity(newSword(tmpR, 20, 20))
	tmpR.Entities.StoreEntity(newSword(tmpR, 20, 40))
	//tmpR.Entities.StoreEntity(newBat(tmpR, 50, 75))
	tmpR.Entities.StoreEntity(newBat(tmpR, 10, 15))
	m.Store("r1", tmpR)

	/*
	m.Store(newBat(50, 75))
	m.Store(newBat(50, 6))
	m.Store(newBat(50, 6))
	tmpB, _ := m.Load("bat2")
	tmpB.(*Bat).Kind = "key"
	m.Store("bat2", tmpB)
	*/
}