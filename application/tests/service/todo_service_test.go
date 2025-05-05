package service

import (
	"testing"

	"github.com/yuktake/todo-webapp/domain/todo"
	"github.com/yuktake/todo-webapp/mock/repository"
	"github.com/yuktake/todo-webapp/service"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestTodoService_CreateTodo_Success(t *testing.T) {
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
