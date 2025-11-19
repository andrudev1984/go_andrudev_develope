package common

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

type Identifiable struct {
	ID uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),pk"`
}

func (i Identifiable) GetId() uuid.UUID {
	return i.ID
}

type Nameable struct {
	Name        string `bun:"type:varchar(255),notnull"`
	Description string `bun:"type:text,default:''"`
}

func (i Nameable) GetName() string {
	return i.Name
}

func (i Nameable) GetDescription() string {
	return i.Description
}

type NotModifiable struct {
	Identifiable
	Created time.Time `bun:"type:timestamp not null"`
}

type Modifiable struct {
	NotModifiable
	Changed time.Time `bun:"type:timestamp not null"`
}

var _ bun.BeforeAppendModelHook = (*NotModifiable)(nil)
var _ bun.BeforeAppendModelHook = (*Modifiable)(nil)

func (i *NotModifiable) BeforeAppendModel(_ context.Context, query schema.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		i.Created = time.Now().UTC()
	}
	return nil
}

func (m *Modifiable) BeforeAppendModel(_ context.Context, query schema.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.Created = time.Now().UTC()
		m.Changed = time.Now().UTC()
	case *bun.UpdateQuery:
		m.Changed = time.Now().UTC()
	}
	return nil
}
