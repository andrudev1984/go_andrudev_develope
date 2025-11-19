package common

import (
	"github.com/google/uuid"
)

type IRepository[T any] interface {
	FindById(uuid uuid.UUID) (*T, error)
}
