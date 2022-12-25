package main

import "sync"

func InitializeRooms(m *sync.Map) {
	r1key, r1ptr := newRoom("r1", defaultRoomUpdate, &r1Layout, nil, nil, nil, nil)
	r1ptr.wallColor = "#8c58b8"
	r2key, r2ptr := newRoom("r2", defaultRoomUpdate, &r1Layout, nil, nil, nil, r1ptr)
	r2ptr.wallColor = "#8c58b8"
	r1ptr.aboveRoom = r2ptr
	m.Store(r1key, r1ptr)
	m.Store(r2key, r2ptr)
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