package domain

import (
	"testing"

	"github.com/yuktake/todo-webapp/domain/user"
)

type User = user.User

func TestNameLengthValidation(t *testing.T) {
	// ユーザー名が10文字以上の場合
	user := User{
		Name:     "12345678901",
		Password: "password",
		Email:    "test@test.com",
	}

	// バリデーションを実行
	err := user.Validate()
	// エラーメッセージを確認
	if err.Error() != "Name must be less than 10 characters" {
		t.Errorf("Expected error message 'Name must be less than 10 characters', got '%s'", err.Error())
	}
}
