package user

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        int64  `bun:"id,pk,autoincrement" json:"id"`
	Name      string `bun:"name,notnull" json:"name"`
	Password  string `bun:"password,notnull" json:"password"`
	Email     string `bun:"email,notnull,unique" json:"email"`
	CreatedAt time.Time
	UpdatedAt time.Time `bun:",nullzero"`
}
