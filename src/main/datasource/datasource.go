package datasource

import (
	"context"

	"github.com/uptrace/bun"
)

type Datasource struct {
	Db      *bun.DB
	Context context.Context
}
