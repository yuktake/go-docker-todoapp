package dto

import (
	"github.com/go-playground/validator/v10"
)

type CreateTodoRequest struct {
	UserID  int64  `json:"user_id" validate:"required"`
	Content string `json:"content" validate:"required"`
	Done    bool   `json:"done"`
}

func (c *CreateTodoRequest) Validate() error {
	v := validator.New()

	if err := v.Struct(c); err != nil {
		return err
	}

	return nil
}
