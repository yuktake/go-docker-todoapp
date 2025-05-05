package service

import (
	"strconv"
	"testing"

	"github.com/yuktake/todo-webapp/domain/todo"
	"github.com/yuktake/todo-webapp/domain/user"
	"github.com/yuktake/todo-webapp/mock/repository"
	"github.com/yuktake/todo-webapp/service"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestTodoService_CreateTodo_Success(t *testing.T) {
	// Arange
	// 他のテストで生成されたデータを削除する
	CleanupAllTables(testDB)

	// テスト用のデータを作成
	// Todoに紐付けるユーザを作成
	userEntity := user.User{
		ID:       1,
		Name:     "Test User",
		Password: "password",
		Email:    "test@test.com",
	}
	_, err := userRepository.CreateUser(&userEntity)
	if err != nil {
		t.Errorf("user does not be created properly: %v", err)
	}

	todoRepository := todo.NewTodoRepository(todo.TodoRepositoryParams{
		DB: testDB,
	})
	todoService := service.NewTodoService(
		service.TodoServiceParams{
			Repo: todoRepository,
		},
	)
	todo := todo.Todo{
		UserID:  1,
		Content: "Test Todo",
		Done:    false,
	}

	// Act
	createdTodo, err := todoService.CreateTodo(todo)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, todo.Content, createdTodo.Content)

	// 作成したTodoを取得して確認
	todo, err = todoService.GetTodoByID(strconv.Itoa(int(createdTodo.ID)))
	if err != nil {
		t.Errorf("todo does not be created properly: %v", err)
	}
	// 後処理
}

func TestTodoService_CreateTodo_ByMock_Success(t *testing.T) {
	// 前処理
	// 他のテストで生成されたデータを削除する
	CleanupAllTables(testDB)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockTodoRepository(ctrl)
	todoService := service.NewTodoService(service.TodoServiceParams{
		Repo: mockRepo,
	})
	todo := todo.Todo{
		UserID:  1,
		Content: "Test Todo",
		Done:    false,
	}
	mockRepo.EXPECT().CreateTodo(&todo).Return(todo, nil)
	createdTodo, err := todoService.CreateTodo(todo)

	// モックの期待値と実際の結果を比較
	assert.NoError(t, err)
	assert.Equal(t, todo, createdTodo)

	// 後処理
}
