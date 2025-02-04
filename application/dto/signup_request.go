package dto

import (
	"github.com/go-playground/validator/v10"
)

// ユーザー登録リクエスト
type SignupRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=64"`
	Email    string `json:"email" validate:"required,email"`
}

// パスワードバリデーションメソッド（任意で特殊ルール追加可）
func (s *SignupRequest) Validate() error {
	v := validator.New()

	// validate tagged fields
	if err := v.Struct(v); err != nil {
		return err
	}

	// custom validation

	return nil
}
