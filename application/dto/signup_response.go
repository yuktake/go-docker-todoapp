package dto

import (
	"github.com/yuktake/todo-webapp/domain/user"
)

// ユーザー登録リクエスト
type SignupResponse struct {
	Message string    `json:"message"`
	User    user.User `json:"user"`
}
