package main

// All entities will have at least this info (implements most of EntityInterface except for Update function)
type EntityBase struct {
	key string	// this entity's unique key for inserting into map
	X float64
	Y float64
	s float64	// speed, how much can move each update (not exported)
	vX float64  // current direction vectors (normalized = hypotenuse of 1) 
	vY float64
	Kind string // what kind of entity
	held EntityInterface // should be a pointer reference to an entity, this entity will only be accessed through parent entity
}

func (e *EntityBase) Held() EntityInterface {
	return e.held
}

func (e *EntityBase) SetHeld(p EntityInterface) {
	e.held = p
}

func (e *EntityBase) Key() string {
	return e.key
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

