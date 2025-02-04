package dto

import (
	"github.com/go-playground/validator/v10"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

func (l *LoginRequest) Validate() error {
	v := validator.New()

	if err := v.Struct(l); err != nil {
		return err
	}

	return nil
}
