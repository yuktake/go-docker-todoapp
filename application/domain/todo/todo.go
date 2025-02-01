package todo

import (
	"time"

	"github.com/uptrace/bun"
)

type Todo struct {
	bun.BaseModel `bun:"table:todos,alias:t"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	Content   string    `bun:"content,notnull" json:"content"`
	Done      bool      `bun:"done" json:"done"`
	Until     time.Time `bun:"until,nullzero"`
	CreatedAt time.Time
	UpdatedAt time.Time `bun:",nullzero"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}
