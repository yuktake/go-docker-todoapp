package user

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	Name      string    `bun:"name,notnull" json:"name" validate:"required"`
	Password  string    `bun:"password,notnull" json:"-"`
	Email     string    `bun:"email,notnull,unique" json:"email" validate:"required,email"`
	CreatedAt time.Time `bun:"created_at,notnull"`
	UpdatedAt time.Time `bun:"updated_at,notnull"`
}

// バリデーションメソッド（エンティティ自身でバリデーション）
func (u *User) Validate() error {
	v := validator.New()

	// validate tagged fields
	if err := v.Struct(u); err != nil {
		return err
	}

	// custom validation

	// name length
	if len(u.Name) > 10 {
		return errors.New("Name must be less than 10 characters")
	}

	return nil
}
