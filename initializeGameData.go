package main

import (
    "sync"
	"github.com/wawesomeNOGUI/adventureIO/roomMapData"
)

func InitializeRooms(m *sync.Map) {
	r1key, r1ptr := newRoom("r1", defaultRoomUpdate, &roomMapData.R1Layout, nil, nil, nil, nil)
	r1ptr.wallColor = "#8c58b8"
	r2key, r2ptr := newRoom("r2", defaultRoomUpdate, &roomMapData.R2Layout, nil, nil, nil, r1ptr)
	r2ptr.wallColor = "#442800"
	r1ptr.aboveRoom = r2ptr
	r3key, r3ptr := newRoom("r3", defaultRoomUpdate, &roomMapData.UpDownLayout, nil, nil, r1ptr, nil)
	r3ptr.wallColor = "#fcfc68"
	r1ptr.belowRoom = r3ptr
	m.Store(r1key, r1ptr)
	m.Store(r2key, r2ptr)
	m.Store(r3key, r3ptr)
}

func InitializeEntities(m *sync.Map) {
	// List all the entities you want here
	r, _ := m.Load("r2")
	tmpR := r.(*Room)
	tmpR.Entities.StoreEntity(newSword(tmpR, 15, 55))
	tmpR.Entities.StoreEntity(newSword(tmpR, 30, 55))
	tmpR.Entities.StoreEntity(newSword(tmpR, 45, 55))
	tmpR.Entities.StoreEntity(newSword(tmpR, 60, 55))
	tmpR.Entities.StoreEntity(newSword(tmpR, 75, 55))
	tmpR.Entities.StoreEntity(newSword(tmpR, 90, 55))
	tmpR.Entities.StoreEntity(newSword(tmpR, 105, 55))
	tmpR.Entities.StoreEntity(newSword(tmpR, 120, 55))
	tmpR.Entities.StoreEntity(newSword(tmpR, 135, 55))
	//tmpR.Entities.StoreEntity(newBat(tmpR, 50, 75))
	// tmpR.Entities.StoreEntity(newBat(tmpR, 10, 15))
	m.Store("r2", tmpR)

	/*
	m.Store(newBat(50, 75))
	m.Store(newBat(50, 6))
	m.Store(newBat(50, 6))
	tmpB, _ := m.Load("bat2")
	tmpB.(*Bat).Kind = "key"
	m.Store("bat2", tmpB)
	*/
}