package todo

import (
	"time"

	"github.com/yuktake/todo-webapp/domain/user"

	"github.com/go-playground/validator/v10"
	"github.com/uptrace/bun"
)

type Todo struct {
	bun.BaseModel `bun:"table:todos,alias:t"`

	ID        int64      `bun:"id,pk,autoincrement" json:"id"`
	UserID    int64      `bun:"user_id,notnull" json:"user_id" validate:"required"`
	User      *user.User `bun:"rel:belongs-to,on_delete:RESTRICT,on_update:RESTRICT" json:"-"`
	Content   string     `bun:"content,notnull" json:"content" validate:"required"`
	Done      bool       `bun:"done,notnull" json:"done" validate:"required"`
	Until     time.Time  `bun:"until,nullzero" json:"until"`
	CreatedAt time.Time  `bun:"created_at,notnull" json:"created_at"`
	UpdatedAt time.Time  `bun:"updated_at,notnull" json:"updated_at"`
}

func (t *Todo) Indexes() []func(*bun.DB) *bun.CreateIndexQuery {
	return []func(*bun.DB) *bun.CreateIndexQuery{}
}

// バリデーションメソッド（エンティティ自身でバリデーション）
func (t *Todo) Validate() error {
	v := validator.New()

	// validate tagged fields
	if err := v.Struct(t); err != nil {
		return err
	}

	// custom validation

	return nil
}
