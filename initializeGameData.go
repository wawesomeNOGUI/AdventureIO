package main

import (
    "sync"
	"github.com/wawesomeNOGUI/adventureIO/roomMapData"
)

var RespawnRoomPtr *Room

func InitializeRooms(m *sync.Map) {
	// this is a special room just to store entities waiting to respawn
	r0key, r0ptr := newRoom("r0", respawnRoomUpdate, nil, nil, nil, nil, nil)
	RespawnRoomPtr = r0ptr

	r1key, r1ptr := newRoom("r1", defaultRoomUpdate, &roomMapData.R1Layout, nil, nil, nil, nil)
	r1ptr.wallColor = "#8c58b8"

	r2key, r2ptr := newRoom("r2", defaultRoomUpdate, &roomMapData.R2Layout, nil, nil, nil, r1ptr)
	r2ptr.wallColor = "#442800"
	r1ptr.aboveRoom = r2ptr

	r3key, r3ptr := newRoom("r3", castleRoomUpdate, &roomMapData.CastleLayout, nil, nil, r1ptr, nil)
	r3ptr.wallColor = "#fcfc68"
	r1ptr.belowRoom = r3ptr

	r4key, r4ptr := newRoom("r4", defaultRoomUpdate, &roomMapData.R4Layout, nil, nil, r3ptr, nil)
	r4ptr.wallColor = "#74b474"
	r3ptr.belowRoom = r4ptr

	r5key, r5ptr := newRoom("r5", dragonRoomUpdate, &roomMapData.R5Layout, nil, nil, r4ptr, nil)
	r5ptr.wallColor = "#404040"
	r4ptr.belowRoom = r5ptr
	r5ptr.specialVars["dragonBeat"] = false
	
	r6key, r6ptr := newRoom("r6", defaultRoomUpdate, &roomMapData.R6Layout, nil, r4ptr, nil, nil)
	r6ptr.wallColor = "#d084c0"
	r4ptr.leftRoom = r6ptr

	r8key, r8ptr := newRoom("r8", defaultRoomUpdate, &roomMapData.R8Layout, nil, r6ptr, nil, nil)
	r8ptr.wallColor = "#6c6c6c"
	r6ptr.leftRoom = r8ptr
	
	r9key, r9ptr := newRoom("r9", batRoomUpdate, &roomMapData.UpDownLayout, nil, nil, nil, r8ptr)
	r9ptr.wallColor = "#6c6c6c"
	r9ptr.specialVars["batsAwake"] = false
	r9ptr.specialVars["fallAsleepTimer"] = 0
	r8ptr.aboveRoom = r9ptr

	r10key, r10ptr := newRoom("r10", castleRoomUpdate, &roomMapData.CastleLayout, nil, nil, nil, r9ptr)
	r10ptr.wallColor = "#000000"
	r9ptr.aboveRoom = r10ptr

	r11key, r11ptr := newRoom("r11", defaultRoomUpdate, &roomMapData.UpDownLayout, nil, nil, nil, r10ptr)
	r11ptr.wallColor = "#000000"
	r10ptr.aboveRoom = r11ptr
	

	m.Store(r0key, r0ptr)
	m.Store(r1key, r1ptr)
	m.Store(r2key, r2ptr)
	m.Store(r3key, r3ptr)
	m.Store(r4key, r4ptr)
	m.Store(r5key, r5ptr)
	m.Store(r6key, r6ptr)
	m.Store(r8key, r8ptr)
	m.Store(r9key, r9ptr)
	m.Store(r10key, r10ptr)
	m.Store(r11key, r11ptr)

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
	m.Store("r2", tmpR)

	r, _ = m.Load("r1")
	tmpR = r.(*Room)
	tmpR.Entities.StoreEntity(newBat(tmpR, 50, 75))
	m.Store("r1", tmpR)

	r, _ = m.Load("r5")
	tmpR = r.(*Room)
	tmpR.Entities.StoreEntity(newDoorGrate(tmpR, 10, 10))
	tmpR.Entities.StoreEntity(newDragon(tmpR, 50, 75))
	m.Store("r5", tmpR)

	r, _ = m.Load("r9")
	tmpR = r.(*Room)
	tmpR.Entities.StoreEntity(newBat(tmpR, 23, 20))
	tmpR.Entities.StoreEntity(newBat(tmpR, 48, 20))
	tmpR.Entities.StoreEntity(newBat(tmpR, 73, 20))
	tmpR.Entities.StoreEntity(newBat(tmpR, 98, 20))
	tmpR.Entities.StoreEntity(newBat(tmpR, 123, 20))
	// tmpR.Entities.StoreEntity(newBat(tmpR, 140, 20))

	tmpR.Entities.StoreEntity(newBat(tmpR, 23, 40))
	tmpR.Entities.StoreEntity(newBat(tmpR, 48, 40))
	tmpR.Entities.StoreEntity(newBat(tmpR, 73, 40))
	tmpR.Entities.StoreEntity(newBat(tmpR, 98, 40))
	tmpR.Entities.StoreEntity(newBat(tmpR, 123, 40))
	// tmpR.Entities.StoreEntity(newBat(tmpR, 140, 40))

	tmpR.Entities.StoreEntity(newBat(tmpR, 23, 60))
	tmpR.Entities.StoreEntity(newBat(tmpR, 48, 60))
	tmpR.Entities.StoreEntity(newBat(tmpR, 73, 60))
	tmpR.Entities.StoreEntity(newBat(tmpR, 98, 60))
	tmpR.Entities.StoreEntity(newBat(tmpR, 123, 60))
	// tmpR.Entities.StoreEntity(newBat(tmpR, 140, 60))

	key, keyPtr := newKey(tmpR, 47, 53)
	tmpR.Entities.StoreEntity(key, keyPtr)

	tmpR.Entities.StoreEntity(newLockedDoor(tmpR, 57, 3, keyPtr))

	m.Store("r9", tmpR)

	r, _ = m.Load("r11")
	tmpR = r.(*Room)
	tmpR.Entities.StoreEntity(newTrophy(tmpR, 50, 75))
	m.Store("r11", tmpR)
}