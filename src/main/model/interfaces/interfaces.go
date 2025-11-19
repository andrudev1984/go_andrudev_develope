package interfaces

import "github.com/google/uuid"

type Identifiable interface {
	GetId() uuid.UUID
}

type Nameable interface {
	Identifiable
	GetName() string
	GetDescription() string
}
