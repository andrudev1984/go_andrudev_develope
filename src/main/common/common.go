package common

import (
	"time"

	"github.com/google/uuid"
)

type Identifiable struct {
	ID uuid.UUID `bun:"type:uuid,default:uuid_generate_v4(),scanonly,primary"`
}

type Nameable struct {
	Name        string `bun:"type:varchar(255),notnull"`
	Description string `bun:"type:text,default:''"`
}

type NotModifiable struct {
	Identifiable
	Created time.Time `bun:"type:timestamp not null,scanonly,default:now()"`
}

type Modifiable struct {
	NotModifiable
	Changed time.Time `bun:"type:timestamp not null,default:now()"`
}
