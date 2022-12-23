package main

// All entities will have at least this info (implements most of EntityInterface except for Update function)
type EntityBase struct {
	key string	// this entity's unique key for inserting into map
	X float64
	Y float64
	width float64
	height float64
	s float64	// speed, how much can move each update (not exported)
	vX float64  // current direction vectors (normalized = hypotenuse of 1) 
	vY float64
	K string // what kind of entity
	held EntityInterface // should be a pointer reference to an entity, this entity will only be accessed through parent entity
	owner EntityInterface // for this entity to use the x, y data from parent entity's Update
	room *Room
}

func (e *EntityBase) Held() EntityInterface {
	return e.held
}

func (e *EntityBase) SetHeld(p EntityInterface) {
	e.held = p
}

func (e *EntityBase) Owner() EntityInterface {
	return e.owner
}

func (e *EntityBase) SetOwner(p EntityInterface) {
	e.owner = p
}

func (e *EntityBase) Key() string {
	return e.key
}

func (e *EntityBase) GetKind() string {
	return e.K
}

func (e *EntityBase) GetX() float64 {
	return e.X
}

func (e *EntityBase) SetX(x float64) {
	e.X = x
}

func (e *EntityBase) GetY() float64 {
	return e.Y
}

func (e *EntityBase) SetY(y float64) {
	e.Y = y
}

func (e *EntityBase) GetWidth() float64 {
	return e.width
}

func (e *EntityBase) SetWidth(w float64) {
	e.width = w
}

func (e *EntityBase) GetHeight() float64 {
	return e.width
}

func (e *EntityBase) SetHeight(h float64) {
	e.height = h
}

func (e *EntityBase) GetvX() float64 {
	return e.vX
}

func (e *EntityBase) SetvX(vX float64) {
	e.vX = vX
}

func (e *EntityBase) GetvY() float64 {
	return e.vY
}

func (e *EntityBase) SetvY(vY float64) {
	e.vY = vY
}

func (e *EntityBase) GetS() float64 {
	return e.s
}

func (e *EntityBase) SetS(s float64) {
	e.s = s
}

func (e *EntityBase) GetRoom() *Room {
	return e.room
}

func (e *EntityBase) SetRoom(r *Room) {
	e.room = r
}

